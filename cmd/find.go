package cmd

import (
	"fmt"
	"strconv"

	"github.com/charmbracelet/lipgloss"
	"github.com/shirou/gopsutil/v3/net"
	"github.com/spf13/cobra"
	"netmon/formatter"
	"netmon/style"
	"netmon/utils"
)

const (
	// minPort is the minimum valid port number
	minPort = 1
	// maxPort is the maximum valid port number
	maxPort = 65535
)

// newFindCmd creates and returns the find command.
// It finds which process is using a specific port.
func newFindCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "find <port>",
		Short: "Find process using a specific port",
		Long:  `Find which process is using a specific port.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			portStr := args[0]
			port, err := strconv.Atoi(portStr)
			if err != nil {
				return fmt.Errorf("invalid port number: %s", portStr)
			}

			if port < minPort || port > maxPort {
				return fmt.Errorf("port number must be between %d and %d", minPort, maxPort)
			}

			// 모든 네트워크 연결 가져오기
			connections, err := net.Connections("inet")
			if err != nil {
				return fmt.Errorf("failed to get network connection information: %w", err)
			}

			// LISTEN 상태이면서 특정 포트를 사용하는 연결만 필터링 (한 번의 순회로 처리)
			foundConns := make(map[string]net.ConnectionStat)
			for _, conn := range connections {
				// LISTEN 상태 또는 UDP이면서 포트가 일치하는 경우
				isListening := conn.Status == "LISTEN" || (utils.IsUDP(conn.Type) && conn.Laddr.Port > 0)
				if isListening && conn.Laddr.Port == uint32(port) {
					key := fmt.Sprintf("%d:%d", conn.Type, conn.Laddr.Port)
					foundConns[key] = conn
				}
			}

			// 결과가 없으면 메시지 출력
			if len(foundConns) == 0 {
				noResultMsg := lipgloss.NewStyle().
					Foreground(style.WarningColor).
					Bold(true).
					Render(fmt.Sprintf("No process found using port %d", port))
				fmt.Println(noResultMsg)
				return nil
			}

			// 포맷터를 사용하여 테이블 생성
			fmtter := formatter.NewPortTableFormatter()
			table := fmtter.Format(foundConns)
			fmt.Println(table)

			return nil
		},
	}
}
