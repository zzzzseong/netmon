<div align="center">

# 📡 netmon

**A modern, beautiful CLI tool for network monitoring and process management**

[![Go Version](https://img.shields.io/badge/Go-1.25%2B-00ADD8?style=flat&logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Version](https://img.shields.io/badge/version-1.1.2-brightgreen.svg)](https://github.com/zzzzseong/netmon/releases)
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
  - Linux-style output format support
  - Cross-platform compatible

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

### 🍺 Homebrew (Recommended)

```bash
brew install zzzzseong/netmon/netmon
```

### 🔨 Build from Source

**Prerequisites:** Go 1.25 or higher

```bash
# Clone the repository
git clone https://github.com/zzzzseong/netmon.git
cd netmon

# Build
go build -o netmon .

# Install (optional)
sudo mv netmon /usr/local/bin/
```

### 📥 Download Pre-built Binaries

Pre-built binaries for macOS and Linux are available in the [Releases](https://github.com/zzzzseong/netmon/releases) section.

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

**New in v1.1.1:**
- ✨ **USERNAME** column showing process owner
- ✨ **CPU %** column displaying CPU usage
- ✨ **MEM %** column displaying memory usage
- ✨ Center-aligned headers
- ✨ Optimized column widths

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
# or
netmon ip --all
```

**New in v1.1.1:**
- ✨ Shows only **IPv4 addresses** by default
- ✨ Hides interfaces **without IP addresses**
- ✨ Use `-a` flag to show **IPv6 addresses** and all interfaces
- ✨ Center-aligned headers with optimized widths

---

### 🛣️ View Routing Table

Display system routing information with smart filtering:

```bash
# Table format (default)
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

```bash
# Linux-style format (ip route style)
netmon route --format linux
```

**Output:**
```
default via 172.16.3.254 dev en0 src 172.16.0.99
172.16.0.0/22 dev en0 src 172.16.0.99
10.10.0.0/24 dev docker0 src 10.10.0.1
```

**New in v1.1.2:**
- ✨ **Native OS APIs** - Uses Netlink (Linux), BSD routing socket (macOS), Win32 API (Windows)
- ✨ **Smart filtering** - Automatically excludes /32 hosts, link-local, multicast, broadcast routes
- ✨ **Linux-style output** - Support for `--format linux` flag (ip route style)
- ✨ **Reduced output** - From 130+ routes to essential 2-3 routes for better readability
- ✨ **Cross-platform** - Works seamlessly on Linux, macOS (Intel & ARM), and Windows

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
| `ip` | Show network interfaces | `netmon ip` | `-a, --all` (show IPv6) |
| `route` | Display routing table with smart filtering | `netmon route` | `--format table\|linux` |
| `find` | Find process using a port | `netmon find <port>` | - |
| `shutdown` | Shutdown a process | `netmon shutdown <pid>` | - |
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

## 🆕 What's New in v1.1.2

### 🚀 Major Route Command Refactoring
- ✨ **Native OS APIs** - Direct system calls instead of command parsing
  - Linux: `vishvananda/netlink` library (Netlink RTM_GETROUTE)
  - macOS/BSD: `golang.org/x/net/route` (BSD routing socket)
  - Windows: `GetIpForwardTable2()` Win32 API
- ✨ **Smart Route Filtering** - Automatically excludes unnecessary routes
  - /32 single host routes
  - Link-local addresses (169.254.0.0/16)
  - Multicast addresses (224.0.0.0/4)
  - Broadcast addresses (255.255.255.255/32)
  - Loopback addresses (127.0.0.0/8)
- ✨ **Linux-style Output** - New `--format linux` flag for ip route style
- ✨ **Dramatic Output Reduction** - From 130+ routes to 2-3 essential routes
- ✨ **True Cross-platform** - Tested on Linux, macOS (Intel & ARM), Windows

### 🛠️ Technical Improvements
- 🏗️ **Provider Architecture** - Clean OS abstraction layer with RouteProvider interface
- ⚡ **Performance** - Direct API calls eliminate external process overhead
- 🎯 **Accuracy** - Native APIs provide more reliable and detailed routing information
- 🔧 **Maintainability** - OS-specific implementations cleanly separated with build tags

### Previous Updates (v1.1.1)
- ✨ Enhanced port listing with username, CPU, and memory metrics
- ✨ Improved IP command with IPv4 default and smart filtering
- ✨ Center-aligned headers and optimized column widths

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

## 📖 Examples

### Check if port 3000 is in use

```bash
netmon find 3000
```

### Monitor all active ports with process metrics

```bash
netmon ls
```

### View only your active network interfaces

```bash
netmon ip
```

### View all network interfaces including IPv6

```bash
netmon ip -a
```

### View routing table in Linux style

```bash
netmon route --format linux
```

### Shutdown a process by PID

```bash
netmon shutdown 12345
```

### Get help

```bash
netmon help
# or
netmon --help
# or
netmon -h
```

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
