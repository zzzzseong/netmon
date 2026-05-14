package cmd

import (
	"fmt"

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
		Use:   "find <pid|port|name>",
		Short: "Find process by PID, port, or name",
		Long: `Find processes by PID, port number, or process name.
Provide a number to search by both PID and port.
Provide a string to search by process name (partial match, case-insensitive).`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			input := args[0]

			// Auto-search using FindByInput (supports both numbers and names)
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
				case "name":
					typeLabel = fmt.Sprintf("Found by Name: %s (PID: %d)", result.ProcessInfo.Name, result.PID)
				}

				// Display process information (with header)
				info := fmtter.FormatWithCmdline(
					typeLabel,
					int(result.PID),
					result.ProcessInfo.Name,
					status,
					result.Connections,
					result.ProcessInfo.Cmdline,
				)
				fmt.Println(info)
			}

			return nil
		},
	}
}
