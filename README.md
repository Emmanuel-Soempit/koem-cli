# koem-cli

> **koem** — meaning *sacred shrine* in Nigerian cultural myth.

koem-cli acts like a shrine keeper — a local directory where developers come to make sensible port categorisations, reserve ports across environments, and avoid conflicts on their machines.

---

## Installation

Requires [Go](https://go.dev/dl/) 1.21+.

```bash
go install github.com/Emmanuel-Soempit/koem-cli@latest
```

Make sure `$HOME/go/bin` is in your `PATH`:

```bash
echo 'export PATH="$HOME/go/bin:$PATH"' >> ~/.zshrc && source ~/.zshrc
```

---

## Features

### 1. Port Range Labels — [`cmd/add.go`](cmd/add.go) · [`impls/label.go`](impls/label.go)

Define named labels with a minimum and maximum port range. Configurations are persisted to a local `.koem.yaml` file.

```bash
# Add a label
koem-cli label add backend 8000 9000
koem-cli label add database 5432 5499

# List all labels with port usage
koem-cli get labels

# Filter by label
koem-cli get labels -l backend
```

The list shows each label, its port range, ports currently in use, and the number of active reserves. Port usage is detected in real time using `net.Listen` across the range concurrently.

---

### 2. Port Suggestions — [`cmd/free.go`](cmd/free.go) · [`impls/label.go`](impls/label.go)

Suggest 3 free ports within a label's range — one each for Production, Preview, and Development.

```bash
koem-cli get free -l backend
```

```
Suggested free ports for "backend":
┌─────────────┬──────┐
│ ENVIRONMENT │ PORT │
├─────────────┼──────┤
│ Production  │ 8009 │
│ Preview     │ 8010 │
│ Development │ 8011 │
└─────────────┴──────┘
```

---

### 3. Port Reservation — [`cmd/reserves.go`](cmd/reserves.go) · [`impls/daemon.go`](impls/daemon.go) · [`impls/daemon_server.go`](impls/daemon_server.go)

Reserve 3 ports for an application under a label. Unlike suggestions, reserved ports are **actually held open** by a lightweight background daemon using `net.Listen`, so no other process can claim them.

The daemon is auto-started on the first reservation and communicates over a Unix socket (`/tmp/koem.sock`).

```bash
# Reserve ports for an app
koem-cli get reserves add myapp -l backend

# List all reserves
koem-cli get reserves

# List reserves for a specific label
koem-cli get reserves -l backend

# Free up reserves for a label
koem-cli get reserves clear -l backend

# Free up all reserves
koem-cli get reserves clear
```

```
Reserved ports for "myapp" under "backend":
┌─────────────┬──────┐
│ ENVIRONMENT │ PORT │
├─────────────┼──────┤
│ Production  │ 8009 │
│ Preview     │ 8010 │
│ Development │ 8011 │
└─────────────┴──────┘
```

---

## All Commands

```
koem-cli label add <name> <min> <max>        Add a new label

koem-cli get labels                          List all labels
koem-cli get labels -l <label>               Filter by label name
koem-cli get free -l <label>                 Suggest 3 free ports

koem-cli get reserves                        List all reserves
koem-cli get reserves -l <label>             List reserves for a label
koem-cli get reserves add <app> -l <label>   Reserve 3 ports for an app
koem-cli get reserves clear                  Clear all reserves
koem-cli get reserves clear -l <label>       Clear reserves for a label
```

---

## Project Structure

```
.
├── main.go
├── cmd/
│   ├── root.go           # Root command, viper config init
│   ├── label.go          # label command
│   ├── add.go            # label add
│   ├── get.go            # get command + get labels
│   ├── free.go           # get free
│   ├── reserves.go       # get reserves + add + clear
│   └── daemon.go         # hidden daemon run command
└── impls/
    ├── label.go          # Label CRUD, port scanning, free port finding, reserve persistence
    ├── daemon.go         # Daemon client (EnsureDaemon, SendReserve, SendRelease)
    └── daemon_server.go  # Daemon server (holds ports open via net.Listen)
```

---

## Version

```bash
koem-cli --version
```

---

## Repo

[github.com/Emmanuel-Soempit/koem-cli](https://github.com/Emmanuel-Soempit/koem-cli)
