package network

import (
	"fmt"

	"github.com/shirou/gopsutil/v3/net"
	"github.com/shirou/gopsutil/v3/process"
)

// PortInfo represents information about a listening port
type PortInfo struct {
	LocalAddr   string
	LocalPort   uint32
	State       string
	PID         int32
	ProcessName string
}

// ListListeningPorts lists all TCP ports in listening state
func ListListeningPorts() ([]PortInfo, error) {
	var ports []PortInfo

	// Get all TCP connections using gopsutil
	connections, err := net.Connections("tcp")
	if err != nil {
		return nil, fmt.Errorf("error getting TCP connections: %v", err)
	}

	// Filter only LISTEN state connections
	for _, conn := range connections {
		if conn.Status == "LISTEN" {
			processName := "Unknown"

			// Try to get process name
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

// PrintListeningPorts prints listening ports in a formatted way
func PrintListeningPorts() error {
	ports, err := ListListeningPorts()
	if err != nil {
		return err
	}

	if len(ports) == 0 {
		fmt.Println("\nNo listening ports found.")
		return nil
	}

	fmt.Println("\n=== LISTENING PORTS ===")
	fmt.Printf("%-20s %-10s %-15s %-10s %-s\n", "ADDRESS", "PORT", "STATE", "PID", "PROCESS")
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

	fmt.Printf("\nTotal: %d listening port(s)\n", len(ports))
	return nil
}

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
