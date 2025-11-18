# netmon

`netmon` is a Go-based CLI tool for macOS/Linux that allows you to check active TCP/UDP ports with a modern interface and manage processes.

## Features

- **Port List**: Display active TCP/UDP ports in a table format with PID and process name
- **Process Termination**: Interactive cursor-based prompt to confirm and execute process termination
- **Color Support**: Color-coded port status (LISTEN: green, ESTABLISHED: blue)
- **Cross-platform**: Supports both macOS and Linux

## Installation

### Requirements

- Go 1.16 or higher

### Build

```bash
go build -o netmon main.go
```

Or run directly:

```bash
go run main.go ls
```

## Usage

### List Ports

```bash
netmon ls
```

Displays active TCP/UDP ports in the following format:

```
PROTOCOL    LOCAL ADDRESS    STATUS      PID    PROCESS NAME
--------    ------------    ------      ---    ------------
TCP         0.0.0.0:8080    LISTEN      1234   nginx
UDP         127.0.0.1:53     LISTEN      567    systemd-resolved
```

### Kill Process

```bash
netmon kill <pid>
```

Displays process information and prompts for confirmation before termination:

1. Display process information (PID, name, status, listening ports)
2. Select "Kill" or "Cancel" using arrow keys
3. Press Enter to confirm
4. Process is terminated if "Kill" is selected

**Example:**

```bash
netmon kill 1234
```

```
=== Process Information ===
PID:        1234
Name:       nginx
Status:     [S]

Listening Ports:
  - 0.0.0.0:8080 (TCP)

Do you want to kill this process?
Select: [Kill/Cancel]
```

## Dependencies

- `github.com/shirou/gopsutil/v3` - Network and process information retrieval
- `github.com/manifoldco/promptui` - Interactive prompts
- `github.com/fatih/color` - Terminal color output

## License

MIT
