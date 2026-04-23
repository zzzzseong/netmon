package cmd

import (
	"fmt"
	"strconv"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/manifoldco/promptui"
	"github.com/shirou/gopsutil/v3/net"
	"github.com/shirou/gopsutil/v3/process"
	"github.com/spf13/cobra"
	"netmon/formatter"
	"netmon/style"
)

// newShutdownCmd creates and returns the shutdown command.
// It safely shuts down a process with interactive confirmation.
func newShutdownCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "shutdown <pid>",
		Short: "Shutdown a process",
		Long:  `Safely shutdown a process with interactive confirmation.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			pidStr := args[0]
			pid, err := strconv.Atoi(pidStr)
			if err != nil {
				return fmt.Errorf("invalid PID: %s", pidStr)
			}

			// Validate PID is positive
			if pid <= 0 {
				return fmt.Errorf("invalid PID: %d (PID must be greater than 0)", pid)
			}

			// Get process information
			proc, err := process.NewProcess(int32(pid))
			if err != nil {
				return fmt.Errorf("process with PID %d not found: %w", pid, err)
			}

			// Get process name
			processName, err := proc.Name()
			if err != nil {
				processName = "N/A"
			}

			// Get process status information
			status, err := proc.Status()
			if err != nil {
				status = []string{"N/A"}
			}

			// Get ports used by the process
			connections, err := net.ConnectionsPid("inet", int32(pid))
			if err != nil {
				connections = []net.ConnectionStat{}
			}

			// Format and display process information (with warning header)
			fmtter := formatter.NewProcessInfoFormatter()
			header := fmt.Sprintf("⚠️  Shutdown Confirmation for PID %d", pid)
			info := fmtter.Format(header, pid, processName, status, connections)
			fmt.Println(info)

			// Interactive confirmation prompt
			promptTitle := lipgloss.NewStyle().
				Foreground(style.WarningColor).
				Bold(true).
				MarginTop(1).
				Render("Are you sure you want to shutdown this process?")
			fmt.Println(promptTitle)

			// Custom template to remove ?: prompt
			templates := &promptui.SelectTemplates{
				Label:    "{{ . }}",
				Active:   "▸ {{ . | cyan }}",
				Inactive: "  {{ . }}",
				Selected: "{{ . | green }}",
			}

			prompt := promptui.Select{
				Label:     "",
				Items:     []string{"✓ Yes, shutdown", "✗ No, cancel"},
				Templates: templates,
				Size:      2,
				HideHelp:  true, // Hide "Use the arrow keys" help message
			}

			index, result, err := prompt.Run()
			if err != nil {
				return fmt.Errorf("prompt error: %w", err)
			}

			if index == 0 && result == "✓ Yes, shutdown" {
				// Try graceful shutdown first (SIGTERM)
				err = proc.Terminate()
				if err != nil {
					// If SIGTERM fails, try SIGKILL
					err = proc.Kill()
					if err != nil {
						return fmt.Errorf("failed to shutdown process: %w", err)
					}
				} else {
					// Wait up to 5 seconds for graceful shutdown
					terminated := false
					for i := 0; i < 10; i++ {
						time.Sleep(500 * time.Millisecond)
						running, _ := proc.IsRunning()
						if !running {
							terminated = true
							break
						}
					}
					if !terminated {
						// Force kill if still running
						_ = proc.Kill()
					}
				}

				// Success message box
				successBox := lipgloss.NewStyle().
					Border(lipgloss.RoundedBorder()).
					BorderForeground(style.SuccessColor).
					Padding(1, 2).
					Foreground(style.SuccessColor).
					Bold(true).
					Render(fmt.Sprintf("✓ Process %d (%s) has been successfully shut down.", pid, processName))
				fmt.Printf("\n%s\n", successBox)
			} else {
				// Cancel message box
				cancelBox := lipgloss.NewStyle().
					Border(lipgloss.RoundedBorder()).
					BorderForeground(style.SubtleColor).
					Padding(1, 2).
					Foreground(style.SubtleColor).
					Render("✗ Shutdown cancelled.")
				fmt.Printf("\n%s\n", cancelBox)
			}

			return nil
		},
	}
}
