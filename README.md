# Network Toolkit 🔧

Swiss army knife for network management activities developed in Go.

## 📋 Description

Network Toolkit is a command-line application that provides advanced tools for system administrators and security professionals to manage, monitor, and audit network connections. The application offers an interactive and easy-to-use interface, with functionalities equivalent to nmap, netstat, and other essential network tools.

## ✨ Implemented Features

### 1. List Listening Ports
Alternative to `netstat -tuln` command (Linux) or `Get-NetTCPConnection -State Listen` (PowerShell).

Displays all TCP ports in listening state with:
- ✅ Local address
- ✅ Port
- ✅ Connection state
- ✅ Process PID
- ✅ Process name

**Helper Functions:**
- `GetListeningPortsCount()` - Returns the number of listening ports
- `IsPortListening(port)` - Checks if a specific port is listening
- `GetProcessByPort(port)` - Returns the process using a port

### 2. Network Scanner (nmap -sS -sV -p-)
Complete network scanner for multiple hosts in CIDR notation.

Features:
- ✅ CIDR network parsing (e.g., 192.168.1.0/24)
- ✅ Automatic detection of active hosts
- ✅ Parallel TCP port scanning
- ✅ Identification of 20+ common services
- ✅ Banner grabbing for advanced detection
- ✅ Thread configuration (1-100)
- ✅ Multiple port range options
- ✅ Detailed report with statistics

**Port Options:**
- Common ports (~20 main ports)
- Specific range (e.g., 1-1024)
- Custom ports (e.g., 80,443,8080)

### 3. Stealth Single-Host Scanner (nmap -sS -sV -p- -T4 --reason)
Aggressive scanner focused on a single target with maximum performance.

Features:
- ✅ TCP SYN Scan (stealth mode)
- ✅ Service version detection (-sV)
- ✅ Aggressive T4 timing (up to 200 threads)
- ✅ Reason analysis (--reason): syn-ack, conn-refused, timeout
- ✅ Port states: open, closed, filtered
- ✅ Banner grabbing with version extraction
- ✅ Real-time progress
- ✅ Time estimation before scan

**Scan Modes:**
- **Quick**: Ports 1-1024 (~20 seconds)
- **Full**: All 65535 ports (~5-10 minutes)
- **Custom**: User-defined range

## 🚀 Installation

### Prerequisites
- Go 1.21 or higher
- Administrator privileges (recommended to view all processes)

### Compile

```bash
# Navigate to the project directory
cd network-toolkit

# Download dependencies
go mod download

# Compile the executable
go build -o network-toolkit.exe
```

## 💻 Usage

### Run the Application

```bash
# Windows (recommended: run as Administrator)
.\network-toolkit.exe
```

### Interactive Menu
The application presents an interactive menu:

```
============================================================
  Network Toolkit 🔧 - v1.2.0
  Swiss army knife for network management activities
============================================================

------------------------------------------------------------
MAIN MENU
------------------------------------------------------------
[1] List Listening Ports (netstat -tuln)
[2] Network Scanner (nmap -sS -sV -p-)
[3] Stealth Single-Host Scanner (nmap -sS -sV -p- -T4)
[0] Exit
------------------------------------------------------------
```

### Example Output - Listening Ports

```
=== LISTENING PORTS ===
ADDRESS              PORT       STATE           PID        PROCESS
--------------------------------------------------------------------------------------------
0.0.0.0              80         LISTEN          1234       nginx.exe
0.0.0.0              443        LISTEN          1234       nginx.exe
127.0.0.1            3306       LISTEN          5678       mysqld.exe
0.0.0.0              8080       LISTEN          9012       java.exe

Total: 4 listening port(s)
```

### Example Output - Network Scanner

```
🔍 Starting network scan: 192.168.1.0/24
📊 Hosts to scan: 254
🔌 Ports per host: 20
⚙️  Threads: 10

✅ 192.168.1.1 - 4 open port(s)
✅ 192.168.1.20 - 6 open port(s)

================================================================================
📊 NETWORK SCAN REPORT
================================================================================

🖥️  HOST: 192.168.1.1 (router.local)
   Scan time: 2.3s
   🔓 Open ports: 4

   PORT       SERVICE              BANNER
   ----------------------------------------------------------------------
   80         HTTP                 nginx/1.18.0
   443        HTTPS                
   22         SSH                  OpenSSH_8.2p1
   8080       HTTP-Proxy           
```

### Example Output - Stealth Scanner

```
🎯 TARGET: 192.168.1.20 (server.local)
🔍 Scanning 65535 ports (range: 1-65535)
⚙️  Threads: 100 | Timeout: 1s | Timing: Aggressive (T4)

✅ Port 22/tcp      open    SSH
✅ Port 80/tcp      open    HTTP
✅ Port 443/tcp     open    HTTPS
⏳ Progress: 25% (16384/65535 ports scanned)

================================================================================
🎯 STEALTH SCAN REPORT (NMAP-LIKE)
================================================================================

📍 TARGET: 192.168.1.20 (server.local)
⏱️  Duration: 5m 23s

📊 STATISTICS
   🟢 Open:   8
   🔴 Closed:  65520
   🟡 Filtered: 7

🔓 DETECTED OPEN PORTS
PORT       STATE      SERVICE         REASON               VERSION/BANNER
----------------------------------------------------------------------------------
22         open       SSH             syn-ack              OpenSSH_8.2p1 Ubuntu
80         open       HTTP            syn-ack              nginx/1.18.0
443        open       HTTPS           syn-ack              nginx/1.18.0
3306       open       MySQL           syn-ack              MySQL 8.0.28
```

## 📁 Project Structure

```
network-toolkit/
├── main.go                          # Application entry point and interactive menu
├── network/
│   ├── listening_ports.go           # Listening ports module
│   ├── port_scanner.go              # CIDR network scanner
│   └── port_scanner_stealthy.go     # Single-host stealth scanner
├── go.mod                           # Dependency management
├── go.sum                           # Dependency checksums
├── .gitignore                       # Files ignored by Git
├── network-toolkit.exe              # Compiled executable
└── README.md                        # This file
```

## 📦 Dependencies

- [`github.com/shirou/gopsutil/v3`](https://github.com/shirou/gopsutil) - Library to get system, process, and network information in a cross-platform manner

## 📝 Important Notes

### Windows
- **Administrator Privileges**: Run the program as Administrator to view complete information for all processes
- **Windows Defender/Antivirus**: Some security solutions may alert about the executable. This is normal for network tools.

### Compatibility
- ✅ Windows 10/11
- ✅ Windows Server 2016+
- ⚠️ Linux (basic functionality - requires testing)
- ⚠️ macOS (basic functionality - requires testing)

### ⚠️ Security Warnings and Ethical Use

**IMPORTANT**: Network scanning features should only be used:
- On networks and systems you own or have explicit authorization for
- For legitimate security auditing purposes
- In your own testing and development environments

**Unauthorized use may:**
- Violate cybercrime laws
- Result in legal action
- Be detected by IDS/IPS systems
- Generate security alerts

**Recommendations:**
- Always obtain written authorization before scanning networks
- Use during low-traffic hours when possible
- Configure appropriate threads and timeouts
- Keep logs of scanning activities
- Respect information security policies

### Known Limitations
- Protected system processes may appear as "Unknown" without administrative privileges
- Performance may vary depending on the number of active connections on the system
- Stealth scanner uses TCP connect scan (not real SYN) due to Go limitations
- OS detection is limited (not fully implemented)
- IPv4 support only at the moment
- Firewalls may block or limit network scans

## 🗺️ Roadmap

### ✅ Version 1.1.0 (Completed)
- [x] Network scanner with CIDR support
- [x] Active host detection
- [x] Parallel TCP port scanning
- [x] Common service identification
- [x] Basic banner grabbing

### ✅ Version 1.2.0 (Completed)
- [x] Single-host stealth scanner
- [x] Aggressive timing (T4)
- [x] Service version detection
- [x] Reason analysis (--reason)
- [x] Port states (open/closed/filtered)
- [x] Real-time progress

### Version 1.3.0 (In Planning)
- [ ] Add UDP port support
- [ ] Implement filters (by port, by process, by address)
- [ ] Add option to export results to CSV/JSON
- [ ] Improve error handling and user messages
- [ ] List all active connections (not just LISTEN)

### Version 2.0.0
- [ ] Connectivity testing (ping, traceroute)
- [ ] Latency and jitter analysis
- [ ] Optional web interface (server mode)
- [ ] Full IPv6 support
- [ ] OS detection (fingerprinting)
- [ ] Continuous monitoring mode

### Future Features
- [ ] Bandwidth monitoring per process
- [ ] Alerts and notifications
- [ ] Connection history
- [ ] Suspicious connection detection
- [ ] Integration with logging tools
- [ ] REST API for integration with other tools
- [ ] Daemon/service mode for continuous monitoring

## 🐛 Known Issues

No critical issues identified at this time.

## 🤝 Contributing

Suggestions and improvements are welcome! This project is under active development.

### How to Contribute
1. Identify a bug or desired feature
2. Implement the solution
3. Test in different scenarios
4. Document the changes

## 📄 License

This project is for internal and educational use.

## 👨‍💻 Development

### Technologies Used
- **Language**: Go 1.21+
- **Libraries**: gopsutil v3
- **Platform**: Multiplatform (run on Windows and Linux, build needed)

### Project Status
🟢 Under active development - v1.2.0

### Last Update
January 8, 2026

### Version History
- **v1.2.0** (01/08/2026) - Single-Host Stealth Scanner
- **v1.1.0** (01/07/2026) - CIDR Network Scanner
- **v1.0.1** (01/07/2026) - Intermediate adjustments
- **v1.0.0** (01/07/2026) - Initial release

---

**Network Toolkit** - Simplifying network management 🚀

