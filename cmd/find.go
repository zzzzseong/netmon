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

			// Check if input is a number
			num, err := strconv.Atoi(input)
			if err != nil {
				return fmt.Errorf("invalid input: %s (must be a number for PID or port)", input)
			}

			// Validate PID/port range
			if num <= 0 {
				return fmt.Errorf("invalid input: %d (must be greater than 0)", num)
			}

			// Auto-search using FindByInput (numbers only)
			results, err := utils.FindByInput(input)
			if err != nil {
				noResultMsg := lipgloss.NewStyle().
					Foreground(style.WarningColor).
					Bold(true).
					Render(fmt.Sprintf("No process found: %s", input))
				fmt.Println(noResultMsg)
				return nil
			}

			// Format and display each result
			fmtter := formatter.NewProcessInfoFormatter()
			for _, result := range results {
				// Get process status information
				proc, err := process.NewProcess(result.PID)
				var status []string
				if err == nil {
					status, _ = proc.Status()
				}
				if len(status) == 0 {
					status = []string{"N/A"}
				}

				// Add message based on result type
				var typeLabel string
				switch result.Type {
				case "pid":
					typeLabel = fmt.Sprintf("Found by PID: %d", result.PID)
				case "port":
					typeLabel = fmt.Sprintf("Found by Port: %d", result.Port)
				}

				// Display process information (with header)
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
