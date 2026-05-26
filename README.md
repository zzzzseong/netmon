<div align="center">

# 📡 netmon

**A modern, beautiful CLI tool for network monitoring and process management**

[![Go Version](https://img.shields.io/badge/Go-1.25%2B-00ADD8?style=flat&logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Version](https://img.shields.io/badge/version-1.6.1-brightgreen.svg)](https://github.com/zzzzseong/netmon/releases)
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
- **📊 Network Statistics** - View network statistics summary at a glance
  - Connection counts (TCP/UDP)
  - Listening ports count
  - Network interfaces count
  - Default gateway information
  - Top processes by connection count
- **🌐 DNS Lookup** - Perform DNS lookups with clean output
  - Forward lookup (A, AAAA, CNAME, MX, NS, and TXT records)
  - Reverse lookup (PTR records)
  - Response time measurement
  - No external dependencies (uses Go's net package)


### 🔍 Process Management
- **Smart Process Search** - Find processes by PID, port, or name with automatic detection
  - Provide a number to search by both PID and port simultaneously
  - Provide a string to search by process name or command line substring
  - Shows full command line (like `ps -ef`) for easy identification
  - Displays all active ports used by the process
- **Process Shutdown** - Safely terminate processes with interactive confirmation
- **Real-time Metrics** - Monitor CPU and memory usage per process

### 🎨 User Experience
- **Beautiful UI** - Modern terminal interface with color-coded output
- **Center-aligned Headers** - Clean, organized table layouts
- **Optimized Columns** - Compact display for better terminal compatibility
- **Automatic Sorting** - Port listings sorted by port number for easy scanning
- **Shell Completion** - Tab completion support for Bash, Zsh, and Fish
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
curl -fsSL https://raw.githubusercontent.com/zzzzseong/netmon/main/scripts/install.sh | bash -s v1.5.0
```

The script automatically:
- Detects OS and architecture (Linux AMD64/ARM64)
- Installs `traceroute` dependency (required for `traceroute` command)
- Downloads the latest version
- Installs to `/usr/local/bin`
- Installs shell completions (Bash, Zsh, Fish)
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

**For Windows:**
- `netmon-windows-amd64.zip` (Intel/AMD 64-bit)

**Installation (Linux / macOS):**

```bash
# Download and extract
tar -xzf netmon-<platform>-<arch>.tar.gz

# Install (optional)
sudo mv netmon /usr/local/bin/

# Verify executable permissions
chmod +x /usr/local/bin/netmon
```

**Installation (Windows):**

```powershell
# Extract the zip file
Expand-Archive netmon-windows-amd64.zip -DestinationPath .

# Run directly
.\netmon.exe --help

# Or move to a directory in your PATH (optional, run as Administrator)
Move-Item netmon.exe C:\Windows\System32\netmon.exe
```

### 🔨 Build from Source

**Prerequisites:**
- Go 1.25 or higher
- `traceroute` command (for `traceroute` functionality)
  - Linux: `sudo apt-get install traceroute` (Debian/Ubuntu) or `sudo yum install traceroute` (RHEL/CentOS)
  - macOS: Usually pre-installed, or `brew install traceroute`
  - Windows: Uses built-in `tracert` command

**Linux / macOS:**

```bash
# Clone the repository
git clone https://github.com/zzzzseong/netmon.git
cd netmon

# Build
go build .

# Install (optional)
sudo mv netmon /usr/local/bin/
```

**Windows:**

```powershell
# Clone the repository
git clone https://github.com/zzzzseong/netmon.git
cd netmon

# Build (produces netmon.exe)
go build .

# Run directly
.\netmon.exe --help

# Or move to a directory in your PATH (optional, run as Administrator)
Move-Item netmon.exe C:\Windows\System32\netmon.exe
```

---

## 🚀 Usage

### 📋 List Active Ports

Display all active listening ports with detailed process information:

```bash
# Show TCP LISTEN connections only (default)
netmon ls
```

**Output:**
```
╭────────────┬─────────────────────┬────────────┬───────────┬───────────────────────────┬─────────────────┬───────────┬──────────╮
│  PROTOCOL  │    LOCAL ADDRESS    │   STATUS   │   PID     │       PROCESS NAME        │    USERNAME     │  CPU %    │  MEM %   │
├────────────┼─────────────────────┼────────────┼───────────┼───────────────────────────┼─────────────────┼───────────┼──────────┤
│ TCP        │ 127.0.0.1:8080      │ LISTEN     │ 12345     │ node                      │ jisung          │ 15.2%     │ 2.1%     │
│ TCP        │ *:3000              │ LISTEN     │ 23456     │ nginx                     │ jisung          │ 2.5%      │ 0.8%     │
╰────────────┴─────────────────────┴────────────┴───────────┴───────────────────────────┴─────────────────┴───────────┴──────────╯
```

```bash
# Include UDP connections with -a flag
netmon ls -a
```

**Output:**
```
╭────────────┬─────────────────────┬────────────┬───────────┬───────────────────────────┬─────────────────┬───────────┬──────────╮
│  PROTOCOL  │    LOCAL ADDRESS    │   STATUS   │   PID     │       PROCESS NAME        │    USERNAME     │  CPU %    │  MEM %   │
├────────────┼─────────────────────┼────────────┼───────────┼───────────────────────────┼─────────────────┼───────────┼──────────┤
│ TCP        │ 127.0.0.1:8080      │ LISTEN     │ 12345     │ node                      │ jisung          │ 15.2%     │ 2.1%     │
│ TCP        │ *:3000              │ LISTEN     │ 23456     │ nginx                     │ jisung          │ 2.5%      │ 0.8%     │
│ UDP        │ 127.0.0.1:53         │ NONE       │ 567       │ systemd-resolved          │ root            │ 0.1%      │ 0.3%     │
╰────────────┴─────────────────────┴────────────┴───────────┴───────────────────────────┴─────────────────┴───────────┴──────────╯
```

---

### 🔗 Show Active Connections

Display active ESTABLISHED TCP connections with owning process information:

```bash
netmon conn
```

**Output:**
```
╭────────────┬───────────────────────────┬───────────────────────────┬───────────┬─────────────────────────╮
│  PROTOCOL  │       LOCAL ADDRESS       │      REMOTE ADDRESS       │   PID     │         PROCESS         │
├────────────┼───────────────────────────┼───────────────────────────┼───────────┼─────────────────────────┤
│ TCP        │ 127.0.0.1:8080            │ 192.168.1.100:54321       │ 12345     │ node                    │
│ TCP        │ 10.0.0.15:5173            │ 142.250.206.206:443       │ 23456     │ chrome                  │
╰────────────┴───────────────────────────┴───────────────────────────┴───────────┴─────────────────────────╯
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

### 📊 View Network Statistics

Display network statistics summary:

```bash
netmon stats
```

**Output:**
```
╭────────────────────────────────────────────────────────────╮
│                                                            │
│  Network Summary                                           │
│                                                            │
│  Active TCP Connections:       150                         │
│  Active UDP Connections:       0                           │
│  Listening Ports:              42                          │
│  Network Interfaces:           19                          │
│  Default Gateway:              192.168.1.1                 │
│                                                            │
│  Top Processes by Connections:                             │
│    • chrome (25 connections)                               │
│    • node (12 connections)                                 │
│    • docker (8 connections)                                │
│                                                            │
╰────────────────────────────────────────────────────────────╯
```

---

### 🌐 Perform DNS Lookup

Lookup DNS records for a domain or IP address:
- Domain lookup returns `A`, `AAAA`, `CNAME`, `MX`, `NS`, and `TXT` records when available
- IP lookup returns `PTR` records for reverse DNS

```bash
# Forward lookup (domain to IP)
netmon dns google.com
```

**Output:**
```
🔍 DNS Lookup: google.com

╭───────────────┬──────────────────────────────────────────────────────────────╮
│     TYPE      │                            VALUE                             │
├───────────────┼──────────────────────────────────────────────────────────────┤
│ A             │ 142.250.206.206                                              │
│ AAAA          │ 2404:6800:400a:813::200e                                     │
│ CNAME         │ google.com.                                                  │
│ MX            │ smtp.google.com. (priority: 10)                              │
│ NS            │ ns1.google.com.                                              │
│ TXT           │ v=spf1 include:_spf.google.com ~all                          │
╰───────────────┴──────────────────────────────────────────────────────────────╯

Response Time: 6ms
```

```bash
# Reverse lookup (IP to domain)
netmon dns 8.8.8.8
```

**Output:**
```
🔍 DNS Lookup: 8.8.8.8

╭───────────────┬──────────────────────────────────────────────────────────────╮
│     TYPE      │                            VALUE                             │
├───────────────┼──────────────────────────────────────────────────────────────┤
│ PTR           │ dns.google.                                                  │
╰───────────────┴──────────────────────────────────────────────────────────────╯

Response Time: 3ms
```


---

### 🔍 Find Process by PID, Port, or Name

Find processes using a smart search that automatically detects PID, port, or process name:

```bash
netmon find <pid|port|name>
```

**How it works:**
- Provide a number, and `find` automatically searches by both PID and port
- Provide a string, and `find` searches by process name or command line substring
- If both match, both results are shown (duplicates removed)
- If only one matches, only that result is displayed
- Shows full command line (like `ps -ef`) for easy process identification
- Displays active ports (LISTEN) and active connections (ESTABLISHED) for the found process

**Examples:**

```bash
# Search by port 8080
netmon find 8080

# Search by PID 8922
netmon find 8922

# Search by process name
netmon find nginx

# If PID 8080 exists AND port 8080 is in use, both results shown
netmon find 8080
```

**Output:**
```
🔍 Found by Port: 8080

╭──────────────────────────────────────────────────╮
│ PID:        12345                               │
│ Name:       node                                │
│ Status:     [sleep]                             │
│                                                 │
│ Command:                                        │
│   /usr/bin/node /path/to/server.js              │
│                                                 │
│ Active Ports:                                   │
│   • 127.0.0.1:8080 (TCP)                        │
│                                                 │
│ Active Connections:                             │
│   • 127.0.0.1:8080 → 192.168.1.100:54321       │
╰──────────────────────────────────────────────────╯
```

*Smart search that works with both PIDs and ports - shows active ports and connections!*

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
⚠️  Shutdown Confirmation for PID 12345

╭─────────────────────────────────────╮
│ PID:        12345                   │
│ Name:       nginx                   │
│ Status:     [S]                     │
│                                     │
│ Active Ports:                       │
│   • 0.0.0.0:8080 (TCP)             │
╰─────────────────────────────────────╯

Are you sure you want to shutdown this process?
▸ ✓ Yes, shutdown
  ✗ No, cancel
```

---

### 🔄 Update netmon

Upgrade netmon to the latest version without re-running the install script:

```bash
netmon update
```

**Output:**
```
# Already on the latest version
Checking for updates...
Already up to date (v1.6.0)

# Update available
Checking for updates...
Updating v1.5.0 → v1.6.0
Downloading netmon-linux-amd64.tar.gz...
Verifying checksum...
Installing to /usr/local/bin/netmon...
Updated to v1.6.0
```

> **Note:** If the install directory (e.g. `/usr/local/bin`) requires root access, netmon automatically retries with `sudo`.

---

## 📚 Commands

| Command | Description | Usage | Flags |
|---------|-------------|-------|-------|
| `ls` | List all active ports with process metrics | `netmon ls` | `-a` (include UDP) |
| `ip` | Show network interfaces | `netmon ip` | `-a` (show IPv6) |
| `route` | Display routing table with smart filtering | `netmon route` | - |
| `stats` | Display network statistics summary | `netmon stats` | - |
| `dns` | Perform DNS lookup | `netmon dns <domain|ip>` | - |
| `conn` | Show active ESTABLISHED connections | `netmon conn` | - |
| `find` | Find process by PID, port, or name (auto-detects) | `netmon find <pid|port|name>` | - |
| `shutdown` | Shutdown a process | `netmon shutdown <pid>` | - |
| `traceroute` | Trace route to network host | `netmon traceroute <host>` | - |
| `version` | Show version information | `netmon version` | - |
| `update` | Update netmon to the latest version | `netmon update` | - |
| `help` | Show help information | `netmon help` | - |

---

## 🎯 Use Cases

- 🔧 **Port Conflict Resolution** - Quickly find which process is using a port (smart PID/port search with active connections)
- 📊 **Network Monitoring** - Monitor network statistics and connections with real-time metrics
- 🛡️ **Process Management** - Safely terminate processes with confirmation
- 💻 **Development** - Check if your development server port is available
- 🌐 **Network Debugging** - View network interfaces and routing information
- 📈 **Performance Monitoring** - Track CPU and memory usage per process

---

## 🆕 What's New in v1.6.1

Released: 2026-05-26

- 👁️ **Watch Mode** - `--watch` / `-w` flag added to `ls`, `conn`, `stats`, `ip`, `route` commands. Refresh interval configurable with `-n` (default: 2s)

### Previous Releases

**v1.6.0:**

- 🔄 **Self-Upgrade Command** - `netmon update` upgrades netmon in place with checksum verification, no need to re-run the install script
- 🐛 **Install Script Fix** - Checksum verification no longer fails on Linux (tarball filename now matches SHA256SUMS entry)

**v1.5.0:**

- 🪟 **Windows Installation Guide** - Full Windows download and installation documentation added
- 🔍 **Find Command Fix** - Full command line now correctly displayed (ps -ef style) when finding processes
- 🌐 **DNS Enhancement** - Underscore-containing domains (SRV, DMARC records like `_dmarc.google.com`) now supported
- 🛡️ **Route Filtering** - IPv6 routes properly excluded from routing table output on all platforms
- 🔧 **Cross-Platform Reliability** - Resolved Windows CI failures (`go vet` unsafe.Pointer, binary extension)
- ⚡ **Signal Handling** - Traceroute cancellation signal handling split by platform (Unix: SIGTERM+Interrupt, Windows: Interrupt only)

**v1.4.0:**

- 📊 **Network Statistics Command** - Quick overview of network activity with connection counts, listening ports, and top processes
- 🌐 **DNS Lookup Command** - Fast DNS lookups (forward and reverse) with response time measurement
- 🔍 **Enhanced Find Command** - Now displays active connections for found processes
- 🎨 **Enhanced Formatting** - Beautiful table layouts for all new commands
- ⚡ **Performance Improvements** - Optimized connection filtering and DNS lookup performance

**v1.2.2:**
- 🔍 Enhanced Find Command - Smart search that automatically detects PID or port, shows full command line
- 📋 Automatic Port Sorting - Port listings are now automatically sorted by port number
- 🛑 Improved Shutdown UI - Better confirmation prompts and styled success/cancel messages
- 🏗️ Code Quality Improvements - Following Go best practices with comprehensive documentation
- 🧪 Enhanced Testing - Added test coverage for multiple packages including benchmarks

**v1.2.1:**
- 🎨 Custom Help Output Restored
- 🧹 Completion Command Cleanup
- 📁 Code Organization improvements

### Shell Completion

Shell completion is automatically installed when you use the install script or Homebrew.

**Linux Installation (via install script):**
```bash
# Install netmon (completions are automatically installed to system directories)
curl -fsSL https://raw.githubusercontent.com/zzzzseong/netmon/main/scripts/install.sh | bash

# Completions are installed to system directories:
# - Bash: /etc/bash_completion.d or /usr/local/etc/bash_completion.d
# - Zsh: /usr/share/zsh/site-functions or /usr/local/share/zsh/site-functions
# 
# Most Linux distributions work without additional configuration.
# Just restart your shell. If completion doesn't work:
# - Bash: Ensure bash-completion package is installed
# - Zsh: Ensure compinit is enabled in your ~/.zshrc
```

**macOS Installation (via Homebrew):**
```bash
# Homebrew automatically installs completions
brew install zzzzseong/netmon/netmon

# Add to ~/.zshrc (shown in installation caveats):
if type brew &>/dev/null; then
  FPATH=$(brew --prefix)/share/zsh/site-functions:$FPATH
  autoload -Uz compinit
  compinit
fi

# Then restart your shell or run:
source ~/.zshrc
```

> 📜 For detailed changelog and previous versions, see [GitHub Releases](https://github.com/zzzzseong/netmon/releases)

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
