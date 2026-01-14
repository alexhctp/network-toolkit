package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"network-toolkit/network"
)

const (
	appTitle   = "Network Toolkit 🔧"
	appVersion = "1.2.0"
)

func main() {
	clearScreen()
	showHeader()

	reader := bufio.NewReader(os.Stdin)

	for {
		showMenu()

		fmt.Print("\nPick a choice: ")
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Error reading input: %v\n", err)
			continue
		}

		choice := strings.TrimSpace(input)

		switch choice {
		case "1":
			handleListeningPorts()
		case "2":
			handleNetworkScan(reader)
		case "3":
			handleStealthyScan(reader)
		case "0":
			fmt.Println("\n👋 Closing Network Toolkit. Goodbye!")
			os.Exit(0)
		default:
			fmt.Println("\n❌ Invalid option! Please choose a valid option.")
		}

		waitForEnter(reader)
	}
}

// showHeader exibe o cabeçalho da aplicação
func showHeader() {
	fmt.Println("=" + strings.Repeat("=", 60))
	fmt.Printf("  %s - v%s\n", appTitle, appVersion)
	fmt.Println("  Swiss army knife for network management activities AKA Jack of all trades")
	fmt.Println("=" + strings.Repeat("=", 60))
}

// showMenu exibe o menu principal
func showMenu() {
	fmt.Println("\n" + strings.Repeat("-", 60))
	fmt.Println("MAIN MENU")
	fmt.Println(strings.Repeat("-", 60))
	fmt.Println("[1] List Listening Ports (netstat -tuln)")
	fmt.Println("[2] Network Scanner (nmap -sS -sV -p-)")
	fmt.Println("[3] Stealth Single-Host Scanner (nmap -sS -sV -p- -T4)")
	fmt.Println("[0] Exit")
	fmt.Println(strings.Repeat("-", 60))
}

// handleListeningPorts trata a opção de listar portas em escuta
func handleListeningPorts() {
	clearScreen()
	fmt.Println("\n🔍 Searching for listening ports...")
	fmt.Println("⚠️  Note: Run as Administrator to see all processes\n")

	err := network.PrintListeningPorts()
	if err != nil {
		fmt.Printf("\n❌ Error listing ports: %v\n", err)
		return
	}

	fmt.Println("\n✅ Operation completed!")
}

// waitForEnter aguarda o usuário pressionar Enter
func waitForEnter(reader *bufio.Reader) {
	fmt.Print("\nPress ENTER to continue...")
	reader.ReadString('\n')
	clearScreen()
	showHeader()
}

// handleNetworkScan trata a opção de scan de rede
func handleNetworkScan(reader *bufio.Reader) {
	clearScreen()
	fmt.Println("\n🔍 NETWORK SCANNER")
	fmt.Println(strings.Repeat("=", 60))
	fmt.Println("\nThis scanner performs an nmap-like scan:")
	fmt.Println("  • Detects active hosts on the network")
	fmt.Println("  • Scans TCP ports")
	fmt.Println("  • Identifies running services")
	fmt.Println("  • Captures service banners\n")

	// Request CIDR network
	fmt.Print("📡 Enter network in CIDR format (e.g., 192.168.1.0/24): ")
	networkInput, _ := reader.ReadString('\n')
	networkInput = strings.TrimSpace(networkInput)

	if networkInput == "" {
		fmt.Println("\n❌ Network cannot be empty!")
		return
	}

	// Request port range
	fmt.Println("\n🔌 Port options:")
	fmt.Println("   [1] Common ports (fast - ~20 ports)")
	fmt.Println("   [2] Specific range (e.g., 1-1024)")
	fmt.Println("   [3] Specific ports (e.g., 80,443,8080)")
	fmt.Print("\nChoose an option [1]: ")
	portOption, _ := reader.ReadString('\n')
	portOption = strings.TrimSpace(portOption)

	if portOption == "" {
		portOption = "1"
	}

	var portRange string
	switch portOption {
	case "1":
		portRange = "all" // Will use common ports
	case "2":
		fmt.Print("Enter range (e.g., 1-1024): ")
		portInput, _ := reader.ReadString('\n')
		portRange = strings.TrimSpace(portInput)
	case "3":
		fmt.Print("Enter ports separated by commas (e.g., 80,443,8080): ")
		portInput, _ := reader.ReadString('\n')
		portRange = strings.TrimSpace(portInput)
	default:
		portRange = "all"
	}

	// Request number of threads
	fmt.Print("\n⚙️  Number of threads [10]: ")
	threadsInput, _ := reader.ReadString('\n')
	threadsInput = strings.TrimSpace(threadsInput)
	threads := 10
	if threadsInput != "" {
		if t, err := strconv.Atoi(threadsInput); err == nil && t > 0 && t <= 100 {
			threads = t
		}
	}

	// Confirmation
	fmt.Println("\n" + strings.Repeat("-", 60))
	fmt.Println("⚠️  WARNING: The network scan may:")
	fmt.Println("   • Take several minutes depending on the network")
	fmt.Println("   • Be detected by security systems")
	fmt.Println("   • Generate significant network traffic")
	fmt.Println(strings.Repeat("-", 60))
	fmt.Print("\nDo you want to continue? (y/N): ")
	confirm, _ := reader.ReadString('\n')
	confirm = strings.ToLower(strings.TrimSpace(confirm))

	if confirm != "y" && confirm != "yes" {
		fmt.Println("\n❌ Scan cancelled.")
		return
	}

	// Configurar scan
	config := network.NetworkScanConfig{
		Network:          networkInput,
		PortRange:        portRange,
		Timeout:          2 * time.Second,
		Threads:          threads,
		ServiceDetection: true,
		OSDetection:      false,
	}

	fmt.Println("\n🚀 Starting scan... Please wait...")
	fmt.Println("")

	// Execute scan
	results, err := network.ScanNetwork(config)
	if err != nil {
		fmt.Printf("\n❌ Error executing scan: %v\n", err)
		return
	}

	// Display results
	network.PrintScanResults(results)

	fmt.Println("\n✅ Scan completed!")
}

// handleStealthyScan trata a opção de scan stealth de host único
func handleStealthyScan(reader *bufio.Reader) {
	clearScreen()
	fmt.Println("\n🎯 STEALTH SINGLE-HOST SCANNER")
	fmt.Println(strings.Repeat("=", 60))
	fmt.Println("\nThis scanner performs a detailed scan on a single target:")
	fmt.Println("  • TCP SYN Scan (stealth)")
	fmt.Println("  • Service version detection")
	fmt.Println("  • Full port scan (1-65535)")
	fmt.Println("  • Aggressive timing (T4)")
	fmt.Println("  • Detection reason (--reason)\n")

	// Request target IP
	fmt.Print("🎯 Enter target IP (e.g., 192.168.1.20): ")
	ipInput, _ := reader.ReadString('\n')
	ipInput = strings.TrimSpace(ipInput)

	if ipInput == "" {
		fmt.Println("\n❌ IP cannot be empty!")
		return
	}

	// Request scan type
	fmt.Println("\n🔍 Scan type:")
	fmt.Println("   [1] Quick - Common ports (1-1024)")
	fmt.Println("   [2] Full - All ports (1-65535)")
	fmt.Println("   [3] Custom - Specific range")
	fmt.Print("\nChoose an option [1]: ")
	scanOption, _ := reader.ReadString('\n')
	scanOption = strings.TrimSpace(scanOption)

	if scanOption == "" {
		scanOption = "1"
	}

	var startPort, endPort, threads int

	switch scanOption {
	case "1":
		startPort = 1
		endPort = 1024
		threads = 50
	case "2":
		startPort = 1
		endPort = 65535
		threads = 100
	case "3":
		fmt.Print("Enter starting port (e.g., 1): ")
		startInput, _ := reader.ReadString('\n')
		startInput = strings.TrimSpace(startInput)
		if s, err := strconv.Atoi(startInput); err == nil && s > 0 && s <= 65535 {
			startPort = s
		} else {
			startPort = 1
		}

		fmt.Print("Enter ending port (e.g., 1000): ")
		endInput, _ := reader.ReadString('\n')
		endInput = strings.TrimSpace(endInput)
		if e, err := strconv.Atoi(endInput); err == nil && e > 0 && e <= 65535 && e >= startPort {
			endPort = e
		} else {
			endPort = 1024
		}

		fmt.Print("Enter number of threads [50]: ")
		threadsInput, _ := reader.ReadString('\n')
		threadsInput = strings.TrimSpace(threadsInput)
		if t, err := strconv.Atoi(threadsInput); err == nil && t > 0 && t <= 200 {
			threads = t
		} else {
			threads = 50
		}
	default:
		startPort = 1
		endPort = 1024
		threads = 50
	}

	// Confirmation
	totalPorts := endPort - startPort + 1
	fmt.Println("\n" + strings.Repeat("-", 60))
	fmt.Printf("⚙️  Scan Configuration:\n")
	fmt.Printf("   Target: %s/32\n", ipInput)
	fmt.Printf("   Range: %d-%d (%d ports)\n", startPort, endPort, totalPorts)
	fmt.Printf("   Threads: %d\n", threads)
	fmt.Printf("   Estimated time: ")

	// Estimate time based on number of ports and threads
	estimatedSeconds := float64(totalPorts) / float64(threads) * 0.5
	if estimatedSeconds < 60 {
		fmt.Printf("~%.0f seconds\n", estimatedSeconds)
	} else {
		fmt.Printf("~%.1f minutes\n", estimatedSeconds/60)
	}

	fmt.Println(strings.Repeat("-", 60))
	fmt.Println("\n⚠️  WARNING:")
	fmt.Println("   • This scan may be detected by IDS/IPS")
	fmt.Println("   • Use only on networks you have authorization for")
	fmt.Println("   • The scan may take time depending on target's firewall")
	fmt.Println(strings.Repeat("-", 60))
	fmt.Print("\nDo you want to continue? (y/N): ")
	confirm, _ := reader.ReadString('\n')
	confirm = strings.ToLower(strings.TrimSpace(confirm))

	if confirm != "y" && confirm != "yes" {
		fmt.Println("\n❌ Scan cancelled.")
		return
	}

	// Configurar e executar scan
	config := network.StealthyScanConfig{
		TargetIP:         ipInput,
		StartPort:        startPort,
		EndPort:          endPort,
		Timeout:          1 * time.Second,
		Threads:          threads,
		ServiceDetection: true,
		AggressiveTiming: true,
	}

	fmt.Println("\n🚀 Starting stealth scan... Please wait...")
	fmt.Println(strings.Repeat("=", 90))

	// Execute scan
	report, err := network.ScanHostStealthy(config)
	if err != nil {
		fmt.Printf("\n❌ Error executing scan: %v\n", err)
		return
	}

	// Display report
	network.PrintStealthyScanReport(report)
}

// clearScreen limpa a tela do terminal
func clearScreen() {
	// Windows
	if os.Getenv("OS") == "Windows_NT" {
		fmt.Print("\033[H\033[2J")
	} else {
		// Unix/Linux/Mac
		fmt.Print("\033[H\033[2J")
	}
}
