package network

import (
	"fmt"

	"github.com/shirou/gopsutil/v3/net"
	"github.com/shirou/gopsutil/v3/process"
)

// PortInfo representa informações sobre uma porta em escuta
type PortInfo struct {
	LocalAddr   string
	LocalPort   uint32
	State       string
	PID         int32
	ProcessName string
}

// ListListeningPorts lista todas as portas TCP em estado de escuta
func ListListeningPorts() ([]PortInfo, error) {
	var ports []PortInfo

	// Obter todas as conexões TCP usando gopsutil
	connections, err := net.Connections("tcp")
	if err != nil {
		return nil, fmt.Errorf("erro ao obter conexões TCP: %v", err)
	}

	// Filtrar apenas conexões em estado LISTEN
	for _, conn := range connections {
		if conn.Status == "LISTEN" {
			processName := "Unknown"

			// Tentar obter o nome do processo
			if conn.Pid > 0 {
				proc, err := process.NewProcess(conn.Pid)
				if err == nil {
					if name, err := proc.Name(); err == nil {
						processName = name
					}
				}
			}

			ports = append(ports, PortInfo{
				LocalAddr:   conn.Laddr.IP,
				LocalPort:   conn.Laddr.Port,
				State:       conn.Status,
				PID:         conn.Pid,
				ProcessName: processName,
			})
		}
	}

	return ports, nil
}

// PrintListeningPorts imprime as portas em escuta de forma formatada
func PrintListeningPorts() error {
	ports, err := ListListeningPorts()
	if err != nil {
		return err
	}

	if len(ports) == 0 {
		fmt.Println("\nNenhuma porta em escuta encontrada.")
		return nil
	}

	fmt.Println("\n=== PORTAS EM ESCUTA ===")
	fmt.Printf("%-20s %-10s %-15s %-10s %-s\n", "ENDEREÇO", "PORTA", "ESTADO", "PID", "PROCESSO")
	fmt.Println("--------------------------------------------------------------------------------------------")

	for _, port := range ports {
		fmt.Printf("%-20s %-10d %-15s %-10d %-s\n",
			port.LocalAddr,
			port.LocalPort,
			port.State,
			port.PID,
			port.ProcessName,
		)
	}

	fmt.Printf("\nTotal: %d porta(s) em escuta\n", len(ports))
	return nil
}

// GetListeningPortsCount retorna o número de portas em escuta
func GetListeningPortsCount() (int, error) {
	ports, err := ListListeningPorts()
	if err != nil {
		return 0, err
	}
	return len(ports), nil
}

// IsPortListening verifica se uma porta específica está em escuta
func IsPortListening(port uint32) (bool, error) {
	ports, err := ListListeningPorts()
	if err != nil {
		return false, err
	}

	for _, p := range ports {
		if p.LocalPort == port {
			return true, nil
		}
	}

	return false, nil
}

// GetProcessByPort retorna o nome do processo usando uma porta específica
func GetProcessByPort(port uint32) (string, error) {
	ports, err := ListListeningPorts()
	if err != nil {
		return "", err
	}

	for _, p := range ports {
		if p.LocalPort == port {
			return p.ProcessName, nil
		}
	}

	return "", fmt.Errorf("porta %d não está em escuta", port)
}
