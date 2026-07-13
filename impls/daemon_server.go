package impls

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

const SocketPath = "/tmp/koem.sock"
const PidFile = "/tmp/koem.pid"

type DaemonCommand struct {
	Action string `json:"action"`
	Ports  []int  `json:"ports"`
}

type DaemonResponse struct {
	OK    bool   `json:"ok"`
	Error string `json:"error,omitempty"`
}

type daemonState struct {
	mu        sync.Mutex
	listeners map[int]net.Listener
}

func RunDaemon() error {
	if err := os.WriteFile(PidFile, []byte(fmt.Sprintf("%d", os.Getpid())), 0644); err != nil {
		return fmt.Errorf("could not write pid file: %w", err)
	}
	defer os.Remove(PidFile)
	defer os.Remove(SocketPath)

	state := &daemonState{
		listeners: make(map[int]net.Listener),
	}

	reserves, err := loadAllReserves()
	if err == nil {
		for _, r := range reserves {
			for _, portStr := range []string{r.Production, r.Preview, r.Development} {
				port := 0
				fmt.Sscanf(portStr, "%d", &port)
				if port > 0 {
					state.holdPort(port)
				}
			}
		}
	}

	os.Remove(SocketPath)
	ln, err := net.Listen("unix", SocketPath)
	if err != nil {
		return fmt.Errorf("could not open unix socket: %w", err)
	}
	defer ln.Close()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		<-sigs
		state.mu.Lock()
		for _, l := range state.listeners {
			l.Close()
		}
		state.mu.Unlock()
		ln.Close()
		os.Exit(0)
	}()

	for {
		conn, err := ln.Accept()
		if err != nil {
			return nil
		}
		go state.handleConn(conn)
	}
}

func (s *daemonState) holdPort(port int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.listeners[port]; exists {
		return
	}
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return
	}
	s.listeners[port] = ln
	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				return
			}
			conn.Close()
		}
	}()
}

func (s *daemonState) handleConn(conn net.Conn) {
	defer conn.Close()
	dec := json.NewDecoder(conn)
	enc := json.NewEncoder(conn)

	var cmd DaemonCommand
	if err := dec.Decode(&cmd); err != nil {
		enc.Encode(DaemonResponse{OK: false, Error: err.Error()})
		return
	}

	switch cmd.Action {
	case "reserve":
		for _, port := range cmd.Ports {
			s.holdPort(port)
		}
		enc.Encode(DaemonResponse{OK: true})
	case "release":
		s.mu.Lock()
		for _, port := range cmd.Ports {
			if ln, exists := s.listeners[port]; exists {
				ln.Close()
				delete(s.listeners, port)
			}
		}
		s.mu.Unlock()
		enc.Encode(DaemonResponse{OK: true})
	case "ping":
		enc.Encode(DaemonResponse{OK: true})
	default:
		enc.Encode(DaemonResponse{OK: false, Error: "unknown action"})
	}
}

func loadAllReserves() ([]Reserve, error) {
	labels, err := LoadAll()
	if err != nil {
		return nil, err
	}
	all := []Reserve{}
	for _, l := range labels {
		r, err := LoadReserves(l.Name)
		if err != nil {
			continue
		}
		all = append(all, r...)
	}
	return all, nil
}
