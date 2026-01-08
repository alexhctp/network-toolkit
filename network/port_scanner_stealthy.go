package network

import (
	"fmt"
	"net"
	"sort"
	"strings"
	"sync"
	"time"
)

// StealthyScanResult representa o resultado detalhado de um scan stealth
type StealthyScanResult struct {
	IP           string
	Port         int
	IsOpen       bool
	State        string // open, closed, filtered
	Service      string
	Version      string
	Banner       string
	ResponseTime time.Duration
	Reason       string // Motivo da detecção
}

// StealthyScanReport relatório completo do scan stealth
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

// StealthyScanConfig configuração do scan stealth
type StealthyScanConfig struct {
	TargetIP         string
	StartPort        int
	EndPort          int
	Timeout          time.Duration
	Threads          int
	ServiceDetection bool
	AggressiveTiming bool // T4 timing
}

// ScanPortStealthy realiza scan stealth em uma porta específica
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

	// Tentar conexão TCP
	conn, err := net.DialTimeout("tcp", address, timeout)
	result.ResponseTime = time.Since(start)

	if err != nil {
		// Analisar o tipo de erro
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

	// Porta está aberta
	result.IsOpen = true
	result.State = "open"
	result.Reason = "syn-ack"

	// Identificar serviço pela porta
	if service, exists := commonServices[port]; exists {
		result.Service = service
	}

	// Detecção de versão através de banner grabbing
	if serviceDetection {
		conn.SetReadDeadline(time.Now().Add(timeout))
		buffer := make([]byte, 2048)
		n, err := conn.Read(buffer)
		if err == nil && n > 0 {
			result.Banner = strings.TrimSpace(string(buffer[:n]))
			// Extrair versão do banner
			result.Version = extractVersionFromBanner(result.Banner)
			// Identificar serviço pelo banner
			result.Service = identifyServiceByBanner(result.Banner, result.Service)
		}
	}

	return result
}

// extractVersionFromBanner extrai informação de versão do banner
func extractVersionFromBanner(banner string) string {
	banner = strings.TrimSpace(banner)
	lines := strings.Split(banner, "\n")
	if len(lines) > 0 {
		firstLine := strings.TrimSpace(lines[0])
		// Limitar tamanho da versão
		if len(firstLine) > 60 {
			return firstLine[:57] + "..."
		}
		return firstLine
	}
	return ""
}

// ScanHostStealthy realiza scan stealth completo em um host
func ScanHostStealthy(config StealthyScanConfig) (*StealthyScanReport, error) {
	report := &StealthyScanReport{
		TargetIP:   config.TargetIP,
		TotalPorts: config.EndPort - config.StartPort + 1,
		ScanDate:   time.Now(),
	}

	// Validar IP
	if net.ParseIP(config.TargetIP) == nil {
		return nil, fmt.Errorf("IP inválido: %s", config.TargetIP)
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

	// Enviar portas para scan
	go func() {
		for port := config.StartPort; port <= config.EndPort; port++ {
			portsChan <- port
		}
		close(portsChan)
	}()

	// Aguardar conclusão dos workers
	go func() {
		wg.Wait()
		close(resultsChan)
	}()

	// Coletar resultados e mostrar progresso
	scanned := 0
	progressInterval := report.TotalPorts / 20 // Mostrar progresso a cada 5%
	if progressInterval < 100 {
		progressInterval = 100
	}

	for result := range resultsChan {
		report.Results = append(report.Results, result)
		scanned++

		// Atualizar contadores
		switch result.State {
		case "open":
			report.OpenPorts++
			// Mostrar portas abertas imediatamente
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

		// Mostrar progresso
		if scanned%progressInterval == 0 {
			progress := float64(scanned) / float64(report.TotalPorts) * 100
			fmt.Printf("⏳ Progresso: %.0f%% (%d/%d portas escaneadas)\n", progress, scanned, report.TotalPorts)
		}
	}

	report.ScanDuration = time.Since(start)

	// Ordenar resultados por número de porta
	sort.Slice(report.Results, func(i, j int) bool {
		return report.Results[i].Port < report.Results[j].Port
	})

	return report, nil
}

// PrintStealthyScanReport imprime o relatório detalhado do scan
func PrintStealthyScanReport(report *StealthyScanReport) {
	fmt.Println("\n" + strings.Repeat("=", 90))
	fmt.Println("🎯 RELATÓRIO DE SCAN STEALTH (NMAP-LIKE)")
	fmt.Println(strings.Repeat("=", 90))

	fmt.Printf("\n📍 TARGET: %s", report.TargetIP)
	if report.Hostname != "" {
		fmt.Printf(" (%s)", report.Hostname)
	}
	fmt.Println()
	fmt.Printf("📅 Data do Scan: %s\n", report.ScanDate.Format("2006-01-02 15:04:05"))
	fmt.Printf("⏱️  Duração: %v\n", report.ScanDuration.Round(time.Millisecond))

	fmt.Println("\n" + strings.Repeat("-", 90))
	fmt.Println("📊 ESTATÍSTICAS")
	fmt.Println(strings.Repeat("-", 90))
	fmt.Printf("Total de portas escaneadas: %d\n", report.TotalPorts)
	fmt.Printf("   🟢 Abertas:   %d\n", report.OpenPorts)
	fmt.Printf("   🔴 Fechadas:  %d\n", report.ClosedPorts)
	fmt.Printf("   🟡 Filtradas: %d\n", report.FilteredPorts)

	// Mostrar apenas portas abertas no relatório final
	if report.OpenPorts > 0 {
		fmt.Println("\n" + strings.Repeat("-", 90))
		fmt.Println("🔓 PORTAS ABERTAS DETECTADAS")
		fmt.Println(strings.Repeat("-", 90))
		fmt.Printf("%-10s %-10s %-15s %-20s %-30s\n", "PORTA", "ESTADO", "SERVIÇO", "RAZÃO", "VERSÃO/BANNER")
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

	// Mostrar portas filtradas se houver
	if report.FilteredPorts > 0 && report.FilteredPorts <= 50 {
		fmt.Println("\n" + strings.Repeat("-", 90))
		fmt.Println("🟡 PORTAS FILTRADAS (Possível Firewall)")
		fmt.Println(strings.Repeat("-", 90))
		fmt.Printf("%-10s %-10s %-20s\n", "PORTA", "ESTADO", "RAZÃO")
		fmt.Println(strings.Repeat("-", 90))

		count := 0
		for _, result := range report.Results {
			if result.State == "filtered" && count < 20 {
				fmt.Printf("%-10d %-10s %-20s\n", result.Port, result.State, result.Reason)
				count++
			}
		}
		if report.FilteredPorts > 20 {
			fmt.Printf("\n... e mais %d porta(s) filtrada(s)\n", report.FilteredPorts-20)
		}
	}

	fmt.Println("\n" + strings.Repeat("=", 90))
	fmt.Println("✅ Scan concluído com sucesso!")
	fmt.Println(strings.Repeat("=", 90) + "\n")
}

// GetCommonPortsRange retorna ranges de portas comuns para scan rápido
func GetCommonPortsRange() []string {
	return []string{
		"1-1024",      // Portas well-known
		"1025-49151",  // Portas registradas
		"49152-65535", // Portas dinâmicas/privadas
	}
}

// QuickScanHost realiza um scan rápido apenas de portas comuns
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

// FullScanHost realiza scan completo de todas as portas
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
