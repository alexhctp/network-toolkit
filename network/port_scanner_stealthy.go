package network

import (
	"fmt"
	"net"
	"sort"
	"strings"
	"sync"
	"time"
)

// StealthyScanResult represents the detailed result of a stealth scan
type StealthyScanResult struct {
	IP           string
	Port         int
	IsOpen       bool
	State        string // open, closed, filtered
	Service      string
	Version      string
	Banner       string
	ResponseTime time.Duration
	Reason       string // Detection reason
}

// StealthyScanReport complete stealth scan report
type StealthyScanReport struct {
	TargetIP      string
	Hostname      string
	TotalPorts    int
	OpenPorts     int
	ClosedPorts   int
	FilteredPorts int
	Results       []StealthyScanResult
	ScanDuration  time.Duration
	ScanDate      time.Time
}

// StealthyScanConfig stealth scan configuration
type StealthyScanConfig struct {
	TargetIP         string
	StartPort        int
	EndPort          int
	Timeout          time.Duration
	Threads          int
	ServiceDetection bool
	AggressiveTiming bool // T4 timing
}

// ScanPortStealthy performs stealth scan on a specific port
func ScanPortStealthy(ip string, port int, timeout time.Duration, serviceDetection bool) StealthyScanResult {
	result := StealthyScanResult{
		IP:      ip,
		Port:    port,
		IsOpen:  false,
		State:   "closed",
		Service: "Unknown",
		Reason:  "no-response",
	}

	start := time.Now()
	address := fmt.Sprintf("%s:%d", ip, port)

	// Try TCP connection
	conn, err := net.DialTimeout("tcp", address, timeout)
	result.ResponseTime = time.Since(start)

	if err != nil {
		// Analyze error type
		if strings.Contains(err.Error(), "timeout") {
			result.State = "filtered"
			result.Reason = "no-response (timeout)"
		} else if strings.Contains(err.Error(), "refused") {
			result.State = "closed"
			result.Reason = "conn-refused"
		} else {
			result.State = "filtered"
			result.Reason = "host-unreachable"
		}
		return result
	}
	defer conn.Close()

	// Port is open
	result.IsOpen = true
	result.State = "open"
	result.Reason = "syn-ack"

	// Identify service by port
	if service, exists := commonServices[port]; exists {
		result.Service = service
	}

	// Version detection through banner grabbing
	if serviceDetection {
		conn.SetReadDeadline(time.Now().Add(timeout))
		buffer := make([]byte, 2048)
		n, err := conn.Read(buffer)
		if err == nil && n > 0 {
			result.Banner = strings.TrimSpace(string(buffer[:n]))
			// Extract version from banner
			result.Version = extractVersionFromBanner(result.Banner)
			// Identify service by banner
			result.Service = identifyServiceByBanner(result.Banner, result.Service)
		}
	}

	return result
}

// extractVersionFromBanner extracts version information from banner
func extractVersionFromBanner(banner string) string {
	banner = strings.TrimSpace(banner)
	lines := strings.Split(banner, "\n")
	if len(lines) > 0 {
		firstLine := strings.TrimSpace(lines[0])
		// Limit version size
		if len(firstLine) > 60 {
			return firstLine[:57] + "..."
		}
		return firstLine
	}
	return ""
}

// ScanHostStealthy performs complete stealth scan on a host
func ScanHostStealthy(config StealthyScanConfig) (*StealthyScanReport, error) {
	report := &StealthyScanReport{
		TargetIP:   config.TargetIP,
		TotalPorts: config.EndPort - config.StartPort + 1,
		ScanDate:   time.Now(),
	}

	// Validate IP
	if net.ParseIP(config.TargetIP) == nil {
		return nil, fmt.Errorf("invalid IP: %s", config.TargetIP)
	}

	// Resolver hostname
	names, err := net.LookupAddr(config.TargetIP)
	if err == nil && len(names) > 0 {
		report.Hostname = names[0]
	}

	start := time.Now()

	fmt.Printf("\n🎯 TARGET: %s", config.TargetIP)
	if report.Hostname != "" {
		fmt.Printf(" (%s)", report.Hostname)
	}
	fmt.Printf("\n")
	fmt.Printf("🔍 Scanning %d ports (range: %d-%d)\n", report.TotalPorts, config.StartPort, config.EndPort)
	fmt.Printf("⚙️  Threads: %d | Timeout: %v | Timing: ", config.Threads, config.Timeout)
	if config.AggressiveTiming {
		fmt.Println("Aggressive (T4)")
	} else {
		fmt.Println("Normal (T3)")
	}
	fmt.Println()

	// Canal para resultados
	resultsChan := make(chan StealthyScanResult, report.TotalPorts)
	portsChan := make(chan int, report.TotalPorts)

	// Worker pool
	var wg sync.WaitGroup
	for i := 0; i < config.Threads; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for port := range portsChan {
				result := ScanPortStealthy(config.TargetIP, port, config.Timeout, config.ServiceDetection)
				resultsChan <- result
			}
		}()
	}

	// Send ports for scanning
	go func() {
		for port := config.StartPort; port <= config.EndPort; port++ {
			portsChan <- port
		}
		close(portsChan)
	}()

	// Wait for workers to complete
	go func() {
		wg.Wait()
		close(resultsChan)
	}()

	// Collect results and show progress
	scanned := 0
	progressInterval := report.TotalPorts / 20 // Show progress every 5%
	if progressInterval < 100 {
		progressInterval = 100
	}

	for result := range resultsChan {
		report.Results = append(report.Results, result)
		scanned++

		// Update counters
		switch result.State {
		case "open":
			report.OpenPorts++
			// Show open ports immediately
			fmt.Printf("✅ Port %d/%s \t%s \t%s\n",
				result.Port,
				"tcp",
				result.State,
				result.Service)
		case "closed":
			report.ClosedPorts++
		case "filtered":
			report.FilteredPorts++
		}

		// Show progress
		if scanned%progressInterval == 0 {
			progress := float64(scanned) / float64(report.TotalPorts) * 100
			fmt.Printf("⏳ Progress: %.0f%% (%d/%d ports scanned)\n", progress, scanned, report.TotalPorts)
		}
	}

	report.ScanDuration = time.Since(start)

	// Sort results by port number
	sort.Slice(report.Results, func(i, j int) bool {
		return report.Results[i].Port < report.Results[j].Port
	})

	return report, nil
}

// PrintStealthyScanReport prints the detailed scan report
func PrintStealthyScanReport(report *StealthyScanReport) {
	fmt.Println("\n" + strings.Repeat("=", 90))
	fmt.Println("🎯 STEALTH SCAN REPORT (NMAP-LIKE)")
	fmt.Println(strings.Repeat("=", 90))

	fmt.Printf("\n📍 TARGET: %s", report.TargetIP)
	if report.Hostname != "" {
		fmt.Printf(" (%s)", report.Hostname)
	}
	fmt.Println()
	fmt.Printf("📅 Scan Date: %s\n", report.ScanDate.Format("2006-01-02 15:04:05"))
	fmt.Printf("⏱️  Duration: %v\n", report.ScanDuration.Round(time.Millisecond))

	fmt.Println("\n" + strings.Repeat("-", 90))
	fmt.Println("📊 STATISTICS")
	fmt.Println(strings.Repeat("-", 90))
	fmt.Printf("Total ports scanned: %d\n", report.TotalPorts)
	fmt.Printf("   🟢 Open:     %d\n", report.OpenPorts)
	fmt.Printf("   🔴 Closed:   %d\n", report.ClosedPorts)
	fmt.Printf("   🟡 Filtered: %d\n", report.FilteredPorts)

	// Show only open ports in final report
	if report.OpenPorts > 0 {
		fmt.Println("\n" + strings.Repeat("-", 90))
		fmt.Println("🔓 DETECTED OPEN PORTS")
		fmt.Println(strings.Repeat("-", 90))
		fmt.Printf("%-10s %-10s %-15s %-20s %-30s\n", "PORT", "STATE", "SERVICE", "REASON", "VERSION/BANNER")
		fmt.Println(strings.Repeat("-", 90))

		for _, result := range report.Results {
			if result.State == "open" {
				version := result.Version
				if version == "" && result.Banner != "" {
					version = result.Banner
				}
				if len(version) > 28 {
					version = version[:25] + "..."
				}

				fmt.Printf("%-10d %-10s %-15s %-20s %-30s\n",
					result.Port,
					result.State,
					result.Service,
					result.Reason,
					version,
				)
			}
		}
	}

	// Show filtered ports if any
	if report.FilteredPorts > 0 && report.FilteredPorts <= 50 {
		fmt.Println("\n" + strings.Repeat("-", 90))
		fmt.Println("🟡 FILTERED PORTS (Possible Firewall)")
		fmt.Println(strings.Repeat("-", 90))
		fmt.Printf("%-10s %-10s %-20s\n", "PORT", "STATE", "REASON")
		fmt.Println(strings.Repeat("-", 90))

		count := 0
		for _, result := range report.Results {
			if result.State == "filtered" && count < 20 {
				fmt.Printf("%-10d %-10s %-20s\n", result.Port, result.State, result.Reason)
				count++
			}
		}
		if report.FilteredPorts > 20 {
			fmt.Printf("\n... and %d more filtered port(s)\n", report.FilteredPorts-20)
		}
	}

	fmt.Println("\n" + strings.Repeat("=", 90))
	fmt.Println("✅ Scan completed successfully!")
	fmt.Println(strings.Repeat("=", 90) + "\n")
}

// GetCommonPortsRange returns common port ranges for quick scan
func GetCommonPortsRange() []string {
	return []string{
		"1-1024",      // Well-known ports
		"1025-49151",  // Registered ports
		"49152-65535", // Dynamic/private ports
	}
}

// QuickScanHost performs a quick scan of common ports only
func QuickScanHost(targetIP string) (*StealthyScanReport, error) {
	config := StealthyScanConfig{
		TargetIP:         targetIP,
		StartPort:        1,
		EndPort:          1024,
		Timeout:          1 * time.Second,
		Threads:          50,
		ServiceDetection: true,
		AggressiveTiming: true,
	}

	return ScanHostStealthy(config)
}

// FullScanHost performs full scan of all ports
func FullScanHost(targetIP string, threads int) (*StealthyScanReport, error) {
	config := StealthyScanConfig{
		TargetIP:         targetIP,
		StartPort:        1,
		EndPort:          65535,
		Timeout:          1 * time.Second,
		Threads:          threads,
		ServiceDetection: true,
		AggressiveTiming: true,
	}

	return ScanHostStealthy(config)
}
