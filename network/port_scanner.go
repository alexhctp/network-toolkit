package network

import (
	"fmt"
	"net"
	"sort"
	"strings"
	"sync"
	"time"
)

// PortScanResult representa o resultado de um scan de porta
type PortScanResult struct {
	IP       string
	Port     int
	IsOpen   bool
	Service  string
	Banner   string
	ScanTime time.Duration
}

// HostScanResult representa o resultado completo de um scan de host
type HostScanResult struct {
	IP         string
	IsAlive    bool
	OpenPorts  []PortScanResult
	OS         string
	Hostname   string
	TotalPorts int
	ScanTime   time.Duration
}

// NetworkScanConfig configurações do scan de rede
type NetworkScanConfig struct {
	Network          string        // CIDR notation (ex: 192.168.1.0/24)
	PortRange        string        // Range de portas (ex: "1-1024" ou "all")
	Timeout          time.Duration // Timeout por porta
	Threads          int           // Número de threads paralelas
	ServiceDetection bool          // Detectar serviços
	OSDetection      bool          // Detectar SO (limitado)
}

// Mapa de serviços comuns por porta
var commonServices = map[int]string{
	20:    "FTP-DATA",
	21:    "FTP",
	22:    "SSH",
	23:    "Telnet",
	25:    "SMTP",
	53:    "DNS",
	80:    "HTTP",
	110:   "POP3",
	143:   "IMAP",
	443:   "HTTPS",
	445:   "SMB",
	3306:  "MySQL",
	3389:  "RDP",
	5432:  "PostgreSQL",
	5900:  "VNC",
	8080:  "HTTP-Proxy",
	8443:  "HTTPS-Alt",
	27017: "MongoDB",
	6379:  "Redis",
}

// ParseCIDR converte CIDR em lista de IPs
func ParseCIDR(cidr string) ([]string, error) {
	ip, ipNet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, fmt.Errorf("CIDR inválido: %v", err)
	}

	var ips []string
	for ip := ip.Mask(ipNet.Mask); ipNet.Contains(ip); incIP(ip) {
		ips = append(ips, ip.String())
	}

	// Remover network address e broadcast address
	if len(ips) > 2 {
		ips = ips[1 : len(ips)-1]
	}

	return ips, nil
}

// incIP incrementa um endereço IP
func incIP(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

// ParsePortRange converte string de range em lista de portas
func ParsePortRange(portRange string) []int {
	if portRange == "all" || portRange == "" {
		// Scan completo seria muito lento, vamos usar portas comuns
		var ports []int
		for port := range commonServices {
			ports = append(ports, port)
		}
		// Adicionar algumas portas adicionais importantes
		additionalPorts := []int{8000, 8008, 8888, 9090, 9200, 9300}
		ports = append(ports, additionalPorts...)
		sort.Ints(ports)
		return ports
	}

	// Parse format: "1-1024" ou "80,443,8080"
	var ports []int

	if strings.Contains(portRange, "-") {
		parts := strings.Split(portRange, "-")
		if len(parts) == 2 {
			var start, end int
			fmt.Sscanf(parts[0], "%d", &start)
			fmt.Sscanf(parts[1], "%d", &end)

			if start > 0 && end <= 65535 && start <= end {
				for i := start; i <= end; i++ {
					ports = append(ports, i)
				}
			}
		}
	} else if strings.Contains(portRange, ",") {
		parts := strings.Split(portRange, ",")
		for _, p := range parts {
			var port int
			if _, err := fmt.Sscanf(strings.TrimSpace(p), "%d", &port); err == nil {
				if port > 0 && port <= 65535 {
					ports = append(ports, port)
				}
			}
		}
	} else {
		var port int
		if _, err := fmt.Sscanf(portRange, "%d", &port); err == nil {
			if port > 0 && port <= 65535 {
				ports = append(ports, port)
			}
		}
	}

	return ports
}

// IsHostAlive verifica se o host está ativo (ping TCP)
func IsHostAlive(ip string, timeout time.Duration) bool {
	// Tenta conectar em portas comuns
	commonPorts := []int{80, 443, 22, 21, 25, 3389}

	for _, port := range commonPorts {
		address := fmt.Sprintf("%s:%d", ip, port)
		conn, err := net.DialTimeout("tcp", address, timeout)
		if err == nil {
			conn.Close()
			return true
		}
	}

	return false
}

// ScanPort escaneia uma porta específica em um IP
func ScanPort(ip string, port int, timeout time.Duration, serviceDetection bool) PortScanResult {
	result := PortScanResult{
		IP:      ip,
		Port:    port,
		IsOpen:  false,
		Service: "Unknown",
	}

	start := time.Now()
	address := fmt.Sprintf("%s:%d", ip, port)

	conn, err := net.DialTimeout("tcp", address, timeout)
	result.ScanTime = time.Since(start)

	if err != nil {
		return result
	}
	defer conn.Close()

	result.IsOpen = true

	// Identificar serviço
	if service, exists := commonServices[port]; exists {
		result.Service = service
	}

	// Detecção de serviço através de banner grabbing
	if serviceDetection {
		conn.SetReadDeadline(time.Now().Add(timeout))
		buffer := make([]byte, 1024)
		n, err := conn.Read(buffer)
		if err == nil && n > 0 {
			result.Banner = strings.TrimSpace(string(buffer[:n]))
			// Tentar identificar serviço pelo banner
			result.Service = identifyServiceByBanner(result.Banner, result.Service)
		}
	}

	return result
}

// identifyServiceByBanner tenta identificar o serviço pelo banner
func identifyServiceByBanner(banner, currentService string) string {
	banner = strings.ToLower(banner)

	if strings.Contains(banner, "ssh") {
		return "SSH"
	} else if strings.Contains(banner, "ftp") {
		return "FTP"
	} else if strings.Contains(banner, "http") || strings.Contains(banner, "html") {
		return "HTTP"
	} else if strings.Contains(banner, "smtp") || strings.Contains(banner, "mail") {
		return "SMTP"
	} else if strings.Contains(banner, "mysql") {
		return "MySQL"
	} else if strings.Contains(banner, "redis") {
		return "Redis"
	}

	return currentService
}

// ScanHost escaneia todos as portas de um host
func ScanHost(ip string, ports []int, config NetworkScanConfig) HostScanResult {
	result := HostScanResult{
		IP:         ip,
		IsAlive:    false,
		OpenPorts:  []PortScanResult{},
		TotalPorts: len(ports),
	}

	start := time.Now()

	// Verificar se host está ativo
	if !IsHostAlive(ip, config.Timeout) {
		result.ScanTime = time.Since(start)
		return result
	}

	result.IsAlive = true

	// Resolver hostname
	names, err := net.LookupAddr(ip)
	if err == nil && len(names) > 0 {
		result.Hostname = names[0]
	}

	// Scan de portas com pool de workers
	var wg sync.WaitGroup
	portChan := make(chan int, len(ports))
	resultChan := make(chan PortScanResult, len(ports))

	// Workers
	for i := 0; i < config.Threads; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for port := range portChan {
				scanResult := ScanPort(ip, port, config.Timeout, config.ServiceDetection)
				if scanResult.IsOpen {
					resultChan <- scanResult
				}
			}
		}()
	}

	// Enviar portas para scan
	go func() {
		for _, port := range ports {
			portChan <- port
		}
		close(portChan)
	}()

	// Aguardar conclusão
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// Coletar resultados
	for scanResult := range resultChan {
		result.OpenPorts = append(result.OpenPorts, scanResult)
	}

	// Ordenar por número de porta
	sort.Slice(result.OpenPorts, func(i, j int) bool {
		return result.OpenPorts[i].Port < result.OpenPorts[j].Port
	})

	result.ScanTime = time.Since(start)
	return result
}

// ScanNetwork escaneia toda a rede
func ScanNetwork(config NetworkScanConfig) ([]HostScanResult, error) {
	// Parse CIDR
	ips, err := ParseCIDR(config.Network)
	if err != nil {
		return nil, err
	}

	// Parse portas
	ports := ParsePortRange(config.PortRange)
	if len(ports) == 0 {
		return nil, fmt.Errorf("nenhuma porta válida especificada")
	}

	fmt.Printf("\n🔍 Iniciando scan de rede: %s\n", config.Network)
	fmt.Printf("📊 Hosts a escanear: %d\n", len(ips))
	fmt.Printf("🔌 Portas por host: %d\n", len(ports))
	fmt.Printf("⚙️  Threads: %d\n", config.Threads)
	fmt.Printf("⏱️  Timeout: %v\n\n", config.Timeout)

	var results []HostScanResult
	var resultsMutex sync.Mutex
	var wg sync.WaitGroup

	// Semáforo para limitar hosts simultâneos
	semaphore := make(chan struct{}, 10)

	for _, ip := range ips {
		wg.Add(1)
		semaphore <- struct{}{} // Adquirir

		go func(targetIP string) {
			defer wg.Done()
			defer func() { <-semaphore }() // Liberar

			result := ScanHost(targetIP, ports, config)

			if result.IsAlive {
				resultsMutex.Lock()
				results = append(results, result)
				resultsMutex.Unlock()

				fmt.Printf("✅ %s - %d porta(s) aberta(s)\n", targetIP, len(result.OpenPorts))
			}
		}(ip)
	}

	wg.Wait()

	return results, nil
}

// PrintScanResults imprime os resultados do scan de forma formatada
func PrintScanResults(results []HostScanResult) {
	if len(results) == 0 {
		fmt.Println("\n❌ Nenhum host ativo encontrado na rede.")
		return
	}

	fmt.Printf("\n\n" + strings.Repeat("=", 80) + "\n")
	fmt.Printf("📊 RELATÓRIO DE SCAN DE REDE\n")
	fmt.Printf(strings.Repeat("=", 80) + "\n\n")

	totalOpenPorts := 0

	for _, host := range results {
		fmt.Printf("🖥️  HOST: %s", host.IP)
		if host.Hostname != "" {
			fmt.Printf(" (%s)", host.Hostname)
		}
		fmt.Printf("\n")
		fmt.Printf("   Tempo de scan: %v\n", host.ScanTime.Round(time.Millisecond))

		if len(host.OpenPorts) == 0 {
			fmt.Printf("   ⚠️  Nenhuma porta aberta encontrada\n\n")
			continue
		}

		fmt.Printf("   🔓 Portas abertas: %d\n\n", len(host.OpenPorts))
		fmt.Printf("   %-10s %-20s %-30s\n", "PORTA", "SERVIÇO", "BANNER")
		fmt.Printf("   " + strings.Repeat("-", 70) + "\n")

		for _, port := range host.OpenPorts {
			banner := port.Banner
			if len(banner) > 28 {
				banner = banner[:25] + "..."
			}
			fmt.Printf("   %-10d %-20s %-30s\n", port.Port, port.Service, banner)
			totalOpenPorts++
		}
		fmt.Printf("\n")
	}

	fmt.Printf(strings.Repeat("=", 80) + "\n")
	fmt.Printf("📈 RESUMO:\n")
	fmt.Printf("   Hosts ativos: %d\n", len(results))
	fmt.Printf("   Total de portas abertas: %d\n", totalOpenPorts)
	fmt.Printf(strings.Repeat("=", 80) + "\n\n")
}
