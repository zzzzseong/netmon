# netmon

📡 **netmon** - A powerful CLI tool for monitoring network connections and managing processes.

A modern, beautifully designed network monitoring tool built with Go that provides an intuitive interface for viewing active ports and managing processes on macOS and Linux.

## ✨ Features

- **📋 Port Listing**: Display all active TCP/UDP ports in a beautifully formatted table
- **🔍 Port Search**: Find processes using a specific port (similar to `lsof -i :port`)
- **🛑 Process Management**: Safely shutdown processes with interactive confirmation
- **🎨 Beautiful UI**: Modern terminal interface with ASCII art and color-coded output
- **⚡ Fast & Lightweight**: Built with Go for optimal performance
- **🔄 Cross-platform**: Works on macOS and Linux

## 📦 Installation

### Homebrew (Recommended)

```bash
brew install zzzzseong/netmon/netmon
```

### Build from Source

**Prerequisites:**
- Go 1.25 or higher

```bash
# Clone the repository
git clone https://github.com/zzzzseong/netmon.git
cd netmon

# Build
go build -o netmon .

# Install (optional)
sudo mv netmon /usr/local/bin/
```

### Run Directly

```bash
go run main.go ls
```

### Download Pre-built Binaries

Pre-built binaries for macOS and Linux are available in the [Releases](https://github.com/zzzzseong/netmon/releases) section.

## 🚀 Usage

### Get Help

```bash
netmon help
# or
netmon --help
# or
netmon -h
```

### List Active Ports

Display all active listening ports with process information:

```bash
netmon ls
```

**Output:**
```
╭───────────────────────┬──────────────────────┬──────────────────────┬──────────────────────┬─────────────────────────╮
│  PROTOCOL             │  LOCAL ADDRESS       │  STATUS              │  PID                 │  PROCESS NAME           │
├───────────────────────┼──────────────────────┼──────────────────────┼──────────────────────┼─────────────────────────┤
│ TCP                   │ *:8080               │ LISTEN               │ 12345                │ nginx                   │
│ UDP                   │ 127.0.0.1:53         │ LISTEN               │ 567                  │ systemd-resolved        │
╰───────────────────────┴──────────────────────┴──────────────────────┴──────────────────────┴─────────────────────────╯
```

### Find Process by Port

Find which process is using a specific port:

```bash
netmon find <port>
```

**Example:**
```bash
netmon find 8080
```

This is equivalent to `lsof -i :8080` but with a more beautiful output format.

### Shutdown Process

Safely shutdown a process with interactive confirmation:

```bash
netmon shutdown <pid>
```

**Example:**
```bash
netmon shutdown 12345
```

The command will:
1. Display detailed process information (PID, name, status, active ports)
2. Show an interactive prompt for confirmation
3. Safely shutdown the process if confirmed

**Interactive Prompt:**
```
⚙️  Process Information
╭─────────────────────────────────────╮
│ PID:        12345                   │
│ Name:       nginx                   │
│ Status:     [S]                     │
│                                     │
│ Active Ports:                       │
│   • 0.0.0.0:8080 (TCP)             │
╰─────────────────────────────────────╯

⚠️  Do you want to shutdown this process?
[Shutdown/Cancel]
```

## 📚 Commands

| Command | Description | Usage |
|---------|-------------|-------|
| `ls` | List all active ports | `netmon ls` |
| `find` | Find process using a port | `netmon find <port>` |
| `shutdown` | Shutdown a process | `netmon shutdown <pid>` |
| `help` | Show help information | `netmon help` |

## 🎯 Use Cases

- **Port Conflict Resolution**: Quickly find which process is using a port
- **Network Monitoring**: Monitor active network connections
- **Process Management**: Safely terminate processes with confirmation
- **Development**: Check if your development server port is available

## 🛠️ Dependencies

- `github.com/shirou/gopsutil/v3` - System and process utilities
- `github.com/charmbracelet/lipgloss` - Terminal styling
- `github.com/manifoldco/promptui` - Interactive prompts

## 📝 Examples

### Check if port 3000 is in use

```bash
netmon find 3000
```

### List all active ports

```bash
netmon ls
```

### Shutdown a process by PID

```bash
netmon shutdown 12345
```

### Get help

```bash
netmon help
```

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🔗 Links

- **GitHub**: https://github.com/zzzzseong/netmon
- **Issues**: https://github.com/zzzzseong/netmon/issues
- **Releases**: https://github.com/zzzzseong/netmon/releases
