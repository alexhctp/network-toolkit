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
	appVersion = "1.1.0"
)

func main() {
	clearScreen()
	showHeader()

	reader := bufio.NewReader(os.Stdin)

	for {
		showMenu()

		fmt.Print("\nEscolha uma opção: ")
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Erro ao ler entrada: %v\n", err)
			continue
		}

		choice := strings.TrimSpace(input)

		switch choice {
		case "1":
			handleListeningPorts()
		case "2":
			handleNetworkScan(reader)
		case "0":
			fmt.Println("\n👋 Encerrando Network Toolkit. Até logo!")
			os.Exit(0)
		default:
			fmt.Println("\n❌ Opção inválida! Por favor, escolha uma opção válida.")
		}

		waitForEnter(reader)
	}
}

// showHeader exibe o cabeçalho da aplicação
func showHeader() {
	fmt.Println("=" + strings.Repeat("=", 60))
	fmt.Printf("  %s - v%s\n", appTitle, appVersion)
	fmt.Println("  Canivete suíço para atividades de gerenciamento de redes")
	fmt.Println("=" + strings.Repeat("=", 60))
}

// showMenu exibe o menu principal
func showMenu() {
	fmt.Println("\n" + strings.Repeat("-", 60))
	fmt.Println("MENU PRINCIPAL")
	fmt.Println(strings.Repeat("-", 60))
	fmt.Println("[1] Listar Portas em Escuta (netstat -tuln)")
	fmt.Println("[2] Scanner de Rede (nmap -sS -sV -p-)")
	fmt.Println("[0] Sair")
	fmt.Println(strings.Repeat("-", 60))
}

// handleListeningPorts trata a opção de listar portas em escuta
func handleListeningPorts() {
	clearScreen()
	fmt.Println("\n🔍 Buscando portas em escuta...")
	fmt.Println("⚠️  Nota: Execute como Administrador para ver todos os processos\n")

	err := network.PrintListeningPorts()
	if err != nil {
		fmt.Printf("\n❌ Erro ao listar portas: %v\n", err)
		return
	}

	fmt.Println("\n✅ Operação concluída!")
}

// waitForEnter aguarda o usuário pressionar Enter
func waitForEnter(reader *bufio.Reader) {
	fmt.Print("\nPressione ENTER para continuar...")
	reader.ReadString('\n')
	clearScreen()
	showHeader()
}

// handleNetworkScan trata a opção de scan de rede
func handleNetworkScan(reader *bufio.Reader) {
	clearScreen()
	fmt.Println("\n🔍 SCANNER DE REDE")
	fmt.Println(strings.Repeat("=", 60))
	fmt.Println("\nEste scanner realiza uma varredura similar ao nmap:")
	fmt.Println("  • Detecta hosts ativos na rede")
	fmt.Println("  • Escaneia portas TCP")
	fmt.Println("  • Identifica serviços em execução")
	fmt.Println("  • Captura banners de serviços\n")

	// Solicitar rede CIDR
	fmt.Print("📡 Digite a rede em formato CIDR (ex: 192.168.1.0/24): ")
	networkInput, _ := reader.ReadString('\n')
	networkInput = strings.TrimSpace(networkInput)

	if networkInput == "" {
		fmt.Println("\n❌ Rede não pode ser vazia!")
		return
	}

	// Solicitar range de portas
	fmt.Println("\n🔌 Opções de portas:")
	fmt.Println("   [1] Portas comuns (rápido - ~20 portas)")
	fmt.Println("   [2] Range específico (ex: 1-1024)")
	fmt.Println("   [3] Portas específicas (ex: 80,443,8080)")
	fmt.Print("\nEscolha uma opção [1]: ")
	portOption, _ := reader.ReadString('\n')
	portOption = strings.TrimSpace(portOption)

	if portOption == "" {
		portOption = "1"
	}

	var portRange string
	switch portOption {
	case "1":
		portRange = "all" // Usará portas comuns
	case "2":
		fmt.Print("Digite o range (ex: 1-1024): ")
		portInput, _ := reader.ReadString('\n')
		portRange = strings.TrimSpace(portInput)
	case "3":
		fmt.Print("Digite as portas separadas por vírgula (ex: 80,443,8080): ")
		portInput, _ := reader.ReadString('\n')
		portRange = strings.TrimSpace(portInput)
	default:
		portRange = "all"
	}

	// Solicitar número de threads
	fmt.Print("\n⚙️  Número de threads [10]: ")
	threadsInput, _ := reader.ReadString('\n')
	threadsInput = strings.TrimSpace(threadsInput)
	threads := 10
	if threadsInput != "" {
		if t, err := strconv.Atoi(threadsInput); err == nil && t > 0 && t <= 100 {
			threads = t
		}
	}

	// Confirmação
	fmt.Println("\n" + strings.Repeat("-", 60))
	fmt.Println("⚠️  AVISO: O scan de rede pode:")
	fmt.Println("   • Demorar vários minutos dependendo da rede")
	fmt.Println("   • Ser detectado por sistemas de segurança")
	fmt.Println("   • Gerar tráfego de rede significativo")
	fmt.Println(strings.Repeat("-", 60))
	fmt.Print("\nDeseja continuar? (s/N): ")
	confirm, _ := reader.ReadString('\n')
	confirm = strings.ToLower(strings.TrimSpace(confirm))

	if confirm != "s" && confirm != "sim" {
		fmt.Println("\n❌ Scan cancelado.")
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

	fmt.Println("\n🚀 Iniciando scan... Por favor, aguarde...")
	fmt.Println("")

	// Executar scan
	results, err := network.ScanNetwork(config)
	if err != nil {
		fmt.Printf("\n❌ Erro ao executar scan: %v\n", err)
		return
	}

	// Exibir resultados
	network.PrintScanResults(results)

	fmt.Println("\n✅ Scan concluído!")
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
