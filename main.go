package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"network-toolkit/network"
)

const (
	appTitle   = "Network Toolkit 🔧"
	appVersion = "1.0.0"
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
