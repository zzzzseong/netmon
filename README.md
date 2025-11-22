<div align="center">

# 📡 netmon

**A modern, beautiful CLI tool for network monitoring and process management**

[![Go Version](https://img.shields.io/badge/Go-1.25%2B-00ADD8?style=flat&logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Version](https://img.shields.io/badge/version-1.1.4-brightgreen.svg)](https://github.com/zzzzseong/netmon/releases)
[![Platform](https://img.shields.io/badge/platform-Linux%20%7C%20macOS%20%7C%20Windows-lightgrey.svg)](https://github.com/zzzzseong/netmon)

*A beautifully designed network monitoring tool built with Go that provides an intuitive interface for viewing active ports, network interfaces, routing tables, and managing processes on Linux, macOS, and Windows.*

[Features](#-features) • [Installation](#-installation) • [Usage](#-usage) • [Commands](#-commands)

</div>

---

## ✨ Features

### 🔌 Network Monitoring
- **📋 Port Listing** - Display all active TCP/UDP ports with detailed process metrics
  - Process owner (username)
  - CPU usage percentage
  - Memory usage percentage
- **🌐 Network Interfaces** - View network interface information with IP addresses
  - IPv4 addresses by default
  - IPv6 support with `-a` flag
  - Filter interfaces with active IPs
- **🛣️ Routing Table** - Display system routing information with smart filtering
  - Native OS APIs (Netlink on Linux, BSD routing socket on macOS, Win32 API on Windows)
  - Smart filtering (excludes /32 hosts, link-local, multicast, broadcast)
  - Cross-platform compatible
- **🗺️ Traceroute** - Trace network path to destination with animated loading
  - Cross-platform support (traceroute/tracert)
  - Beautiful table format with color-coded RTT values
  - Real-time animated spinner during execution

### 🔍 Process Management
- **Port Search** - Find processes using specific ports (like `lsof -i :port`)
- **Process Shutdown** - Safely terminate processes with interactive confirmation
- **Real-time Metrics** - Monitor CPU and memory usage per process

### 🎨 User Experience
- **Beautiful UI** - Modern terminal interface with color-coded output
- **Center-aligned Headers** - Clean, organized table layouts
- **Optimized Columns** - Compact display for better terminal compatibility
- **Fast & Lightweight** - Built with Go for optimal performance

---

## 📦 Installation

### 🍺 Homebrew (macOS Recommended)

```bash
brew install zzzzseong/netmon/netmon
```

### 🐧 Linux - Quick Install Script

The easiest way to install on Linux:

```bash
# Install latest version
curl -fsSL https://raw.githubusercontent.com/zzzzseong/netmon/main/scripts/install.sh | bash

# Install specific version
curl -fsSL https://raw.githubusercontent.com/zzzzseong/netmon/main/scripts/install.sh | bash -s v1.1.4
```

The script automatically:
- Detects OS and architecture (Linux AMD64/ARM64)
- Installs `traceroute` dependency (required for `traceroute` command)
- Downloads the latest version
- Installs to `/usr/local/bin`
- Verifies installation

**After installation, you can immediately use netmon:**

```bash
netmon ls        # List active ports
netmon ip        # Show network interfaces
netmon --help    # Show all commands
```

> **Note:** If `netmon` is not found after installation, try opening a new terminal or run `hash -r` to refresh the command cache.

### 📥 Download Pre-built Binaries

Pre-built binaries are available in the [Releases](https://github.com/zzzzseong/netmon/releases) section.

**For Linux:**
- `netmon-linux-amd64.tar.gz` (Intel/AMD 64-bit)
- `netmon-linux-arm64.tar.gz` (ARM 64-bit)

**For macOS:**
- `netmon-darwin-amd64.tar.gz` (Intel)
- `netmon-darwin-arm64.tar.gz` (Apple Silicon)

**Installation:**

```bash
# Download and extract
tar -xzf netmon-<platform>-<arch>.tar.gz

# Install (optional)
sudo mv netmon /usr/local/bin/

# Verify executable permissions
chmod +x /usr/local/bin/netmon
```

### 🔨 Build from Source

**Prerequisites:**
- Go 1.25 or higher
- `traceroute` command (for `traceroute` functionality)
  - Linux: `sudo apt-get install traceroute` (Debian/Ubuntu) or `sudo yum install traceroute` (RHEL/CentOS)
  - macOS: Usually pre-installed, or `brew install traceroute`
  - Windows: Uses built-in `tracert` command

```bash
# Clone the repository
git clone https://github.com/zzzzseong/netmon.git
cd netmon

# Build
go build -o netmon .

# Install (optional)
sudo mv netmon /usr/local/bin/
```

---

## 🚀 Usage

### 📋 List Active Ports

Display all active listening ports with detailed process information:

```bash
netmon ls
```

**Output:**
```
╭────────────┬─────────────────────┬────────────┬───────────┬───────────────────────────┬─────────────────┬───────────┬──────────╮
│  PROTOCOL  │    LOCAL ADDRESS    │   STATUS   │   PID     │       PROCESS NAME        │    USERNAME     │  CPU %    │  MEM %   │
├────────────┼─────────────────────┼────────────┼───────────┼───────────────────────────┼─────────────────┼───────────┼──────────┤
│ TCP        │ 127.0.0.1:8080      │ LISTEN     │ 12345     │ node                      │ jisung          │ 15.2%     │ 2.1%     │
│ TCP        │ *:3000              │ LISTEN     │ 23456     │ nginx                     │ jisung          │ 2.5%      │ 0.8%     │
│ UDP        │ 127.0.0.1:53        │ LISTEN     │ 567       │ systemd-resolved          │ root            │ 0.1%      │ 0.3%     │
╰────────────┴─────────────────────┴────────────┴───────────┴───────────────────────────┴─────────────────┴───────────┴──────────╯
```

---

### 🌐 View Network Interfaces

Display network interfaces with IP addresses:

```bash
# Show IPv4 addresses only (default)
netmon ip
```

**Output:**
```
╭──────────────┬──────────────────────────────────────────┬──────────────────────┬─────────────┬─────────────╮
│  INTERFACE   │                IP ADDRESS                │     MAC ADDRESS      │  STATUS     │  MTU        │
├──────────────┼──────────────────────────────────────────┼──────────────────────┼─────────────┼─────────────┤
│ en0          │ 192.168.1.100/24                         │ a4:83:e7:5c:5d:3e    │ UP          │ 1500        │
│ lo0          │ 127.0.0.1/8                              │ N/A                  │ UP          │ 16384       │
╰──────────────┴──────────────────────────────────────────┴──────────────────────┴─────────────┴─────────────╯
```

```bash
# Show all addresses including IPv6
netmon ip -a
```

---

### 🛣️ View Routing Table

Display system routing information with smart filtering:

```bash
netmon route
```

**Output:**
```
╭───────────────────────────┬───────────────────────────┬───────────────────────────┬───────────────────────────┬──────────────────────────╮
│  DESTINATION              │  GATEWAY                  │  INTERFACE                │  METRIC                   │  SOURCE                  │
├───────────────────────────┼───────────────────────────┼───────────────────────────┼───────────────────────────┼──────────────────────────┤
│ default                   │ 172.16.3.254              │ en0                       │ -                         │ 172.16.0.99              │
│ 172.16.0.0/22             │ -                         │ en0                       │ -                         │ 172.16.0.99              │
│ 10.10.0.0/24              │ -                         │ docker0                   │ -                         │ 10.10.0.1                │
╰───────────────────────────┴───────────────────────────┴───────────────────────────┴───────────────────────────┴──────────────────────────╯
```

---

### 🗺️ Trace Route to Host

Trace the network path to a destination with animated loading indicator:

```bash
netmon traceroute <host>
```

**Example:**
```bash
netmon traceroute google.com
```

**Output:**
```
◐ Tracing route to google.com...
╭──────────┬──────────────────────────────────────────┬──────────────┬──────────────┬──────────────╮
│  HOP     │                   HOST                   │    RTT 1     │    RTT 2     │    RTT 3     │
├──────────┼──────────────────────────────────────────┼──────────────┼──────────────┼──────────────┤
│ 1        │ 192.168.1.1                              │ 2.5 ms       │ 2.3 ms       │ 2.1 ms       │
│ 2        │ 10.0.0.1                                 │ 15.2 ms      │ 14.8 ms      │ 15.0 ms      │
│ 3        │ 172.217.160.46                           │ 25.3 ms      │ 24.9 ms      │ 25.1 ms      │
╰──────────┴──────────────────────────────────────────┴──────────────┴──────────────┴──────────────╯
```

---

### 🔍 Find Process by Port

Find which process is using a specific port:

```bash
netmon find <port>
```

**Example:**
```bash
netmon find 8080
```

**Output:**
```
⚙️  Process Information
╭─────────────────────────────────────╮
│ PID:        12345                   │
│ Name:       node                    │
│ Status:     [S]                     │
│                                     │
│ Active Ports:                       │
│   • 127.0.0.1:8080 (TCP)           │
╰─────────────────────────────────────╯
```

*This is equivalent to `lsof -i :8080` but with a more beautiful output format.*

---

### 🛑 Shutdown Process

Safely shutdown a process with interactive confirmation:

```bash
netmon shutdown <pid>
```

**Example:**
```bash
netmon shutdown 12345
```

The command will:
1. 📊 Display detailed process information (PID, name, status, active ports)
2. ❓ Show an interactive prompt for confirmation
3. ✅ Safely shutdown the process if confirmed

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

---

## 📚 Commands

| Command | Description | Usage | Flags |
|---------|-------------|-------|-------|
| `ls` | List all active ports with process metrics | `netmon ls` | - |
| `ip` | Show network interfaces | `netmon ip` | `-a` (show IPv6) |
| `route` | Display routing table with smart filtering | `netmon route` | - |
| `find` | Find process using a port | `netmon find <port>` | - |
| `shutdown` | Shutdown a process | `netmon shutdown <pid>` | - |
| `traceroute` | Trace route to network host | `netmon traceroute <host>` | - |
| `version` | Show version information | `netmon version` | - |
| `help` | Show help information | `netmon help` | - |

---

## 🎯 Use Cases

- 🔧 **Port Conflict Resolution** - Quickly find which process is using a port
- 📊 **Network Monitoring** - Monitor active network connections with real-time metrics
- 🛡️ **Process Management** - Safely terminate processes with confirmation
- 💻 **Development** - Check if your development server port is available
- 🌐 **Network Debugging** - View network interfaces and routing information
- 📈 **Performance Monitoring** - Track CPU and memory usage per process

---

## 🆕 What's New in v1.1.4

- 🐧 **Linux installation script** - Easy one-command installation for Linux users
- 📁 **Improved project structure** - Installation scripts organized in `scripts/` directory
- 📝 **README optimization** - Cleaner documentation with removed duplicates


> 📜 For detailed changelog, see [GitHub Releases](https://github.com/zzzzseong/netmon/releases)

---

## 🛠️ Tech Stack

- **Language**: Go 1.25+
- **Core Dependencies**:
  - [`github.com/shirou/gopsutil/v3`](https://github.com/shirou/gopsutil) - System and process utilities
  - [`github.com/charmbracelet/lipgloss`](https://github.com/charmbracelet/lipgloss) - Terminal styling
  - [`github.com/manifoldco/promptui`](https://github.com/manifoldco/promptui) - Interactive prompts
- **Routing System**:
  - [`github.com/vishvananda/netlink`](https://github.com/vishvananda/netlink) - Linux netlink interface
  - [`golang.org/x/net/route`](https://golang.org/x/net/route) - BSD routing socket interface
  - [`golang.org/x/sys/windows`](https://golang.org/x/sys/windows) - Windows system calls

---

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

---

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

## 🔗 Links

- **GitHub**: https://github.com/zzzzseong/netmon
- **Issues**: https://github.com/zzzzseong/netmon/issues
- **Releases**: https://github.com/zzzzseong/netmon/releases

---

<div align="center">

Made with by [zzzzseong](https://github.com/zzzzseong)

⭐ Star this repository if you find it helpful!

</div>
