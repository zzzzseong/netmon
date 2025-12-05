package cmd

import (
	"fmt"
	"strconv"

	"github.com/charmbracelet/lipgloss"
	"github.com/shirou/gopsutil/v3/process"
	"github.com/spf13/cobra"
	"netmon/formatter"
	"netmon/style"
	"netmon/utils"
)

// newFindCmd creates and returns the find command.
// It finds processes by PID, port, or process name.
// If input is a number, it searches by both PID and port.
// If input is a string, it searches by process name.
func newFindCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "find <pid|port>",
		Short: "Find process by PID or port",
		Long: `Find processes by PID or port number.
Provide a number to search by both PID and port.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			input := args[0]

			// 숫자인지 확인
			_, err := strconv.Atoi(input)
			if err != nil {
				return fmt.Errorf("invalid input: %s (must be a number for PID or port)", input)
			}

			// FindByInput으로 자동 검색 (숫자만 허용)
			results, err := utils.FindByInput(input)
			if err != nil {
				noResultMsg := lipgloss.NewStyle().
					Foreground(style.WarningColor).
					Bold(true).
					Render(fmt.Sprintf("No process found: %s", input))
				fmt.Println(noResultMsg)
				return nil
			}

			// 각 결과를 포맷팅하여 출력
			fmtter := formatter.NewProcessInfoFormatter()
			for _, result := range results {
				// 프로세스 상태 정보 가져오기
				proc, err := process.NewProcess(result.PID)
				var status []string
				if err == nil {
					status, _ = proc.Status()
				}
				if len(status) == 0 {
					status = []string{"N/A"}
				}

				// 결과 타입에 따라 메시지 추가
				var typeLabel string
				switch result.Type {
				case "pid":
					typeLabel = fmt.Sprintf("Found by PID: %d", result.PID)
				case "port":
					typeLabel = fmt.Sprintf("Found by Port: %d", result.Port)
				}

				// 프로세스 정보 출력 (헤더 포함)
				info := fmtter.Format(
					typeLabel,
					int(result.PID),
					result.ProcessInfo.Name,
					status,
					result.Connections,
				)
				fmt.Println(info)
			}

			return nil
		},
	}
}
