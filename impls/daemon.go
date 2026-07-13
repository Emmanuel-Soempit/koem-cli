package impls

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"os/exec"
	"strconv"
	"time"
)

func EnsureDaemon() error {
	if isDaemonRunning() {
		return nil
	}
	exe, err := os.Executable()
	if err != nil {
		return fmt.Errorf("could not find koem executable: %w", err)
	}
	cmd := exec.Command(exe, "daemon", "run")
	cmd.Stdout = nil
	cmd.Stderr = nil
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("could not start daemon: %w", err)
	}
	for i := 0; i < 20; i++ {
		time.Sleep(100 * time.Millisecond)
		if isDaemonRunning() {
			return nil
		}
	}
	return fmt.Errorf("daemon did not start in time")
}

func isDaemonRunning() bool {
	conn, err := net.DialTimeout("unix", SocketPath, 500*time.Millisecond)
	if err != nil {
		return false
	}
	defer conn.Close()
	enc := json.NewEncoder(conn)
	dec := json.NewDecoder(conn)
	enc.Encode(DaemonCommand{Action: "ping"})
	var resp DaemonResponse
	dec.Decode(&resp)
	return resp.OK
}

func SendReserve(ports []int) error {
	conn, err := net.DialTimeout("unix", SocketPath, time.Second)
	if err != nil {
		return fmt.Errorf("could not connect to daemon: %w", err)
	}
	defer conn.Close()
	enc := json.NewEncoder(conn)
	dec := json.NewDecoder(conn)
	enc.Encode(DaemonCommand{Action: "reserve", Ports: ports})
	var resp DaemonResponse
	if err := dec.Decode(&resp); err != nil {
		return err
	}
	if !resp.OK {
		return fmt.Errorf("daemon error: %s", resp.Error)
	}
	return nil
}

func SendRelease(ports []int) error {
	conn, err := net.DialTimeout("unix", SocketPath, time.Second)
	if err != nil {
		return nil
	}
	defer conn.Close()
	enc := json.NewEncoder(conn)
	dec := json.NewDecoder(conn)
	enc.Encode(DaemonCommand{Action: "release", Ports: ports})
	var resp DaemonResponse
	if err := dec.Decode(&resp); err != nil {
		return err
	}
	if !resp.OK {
		return fmt.Errorf("daemon error: %s", resp.Error)
	}
	return nil
}

func PortsFromReserve(r Reserve) ([]int, error) {
	prod, err := strconv.Atoi(r.Production)
	if err != nil {
		return nil, err
	}
	prev, err := strconv.Atoi(r.Preview)
	if err != nil {
		return nil, err
	}
	dev, err := strconv.Atoi(r.Development)
	if err != nil {
		return nil, err
	}
	return []int{prod, prev, dev}, nil
}
