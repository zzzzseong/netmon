package cmd

import (
	"fmt"
	"strconv"

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

		// 프로세스 정보 가져오기
		proc, err := process.NewProcess(int32(pid))
		if err != nil {
			return fmt.Errorf("process with PID %d not found: %w", pid, err)
		}

		// 프로세스 이름 가져오기
		processName, err := proc.Name()
		if err != nil {
			processName = "N/A"
		}

		// 프로세스 상태 정보 가져오기
		status, err := proc.Status()
		if err != nil {
			status = []string{"N/A"}
		}

		// 프로세스가 사용하는 포트 가져오기
		connections, err := net.ConnectionsPid("inet", int32(pid))
		if err != nil {
			connections = []net.ConnectionStat{}
		}

		// 포맷터를 사용하여 프로세스 정보 출력
		fmtter := formatter.NewProcessInfoFormatter()
		info := fmtter.Format("", pid, processName, status, connections)
		fmt.Println(info)

		// 인터랙티브 확인 프롬프트
		promptTitle := lipgloss.NewStyle().
			Foreground(style.WarningColor).
			Bold(true).
			MarginTop(1).
			Render("⚠️  Do you want to shutdown this process?")
		fmt.Println(promptTitle)

		prompt := promptui.Select{
			Label: "",
			Items: []string{"Shutdown", "Cancel"},
			Size:  2,
		}

		index, result, err := prompt.Run()
		if err != nil {
			return fmt.Errorf("prompt error: %w", err)
		}

		if index == 0 && result == "Shutdown" {
			// 프로세스 종료
			err = proc.Kill()
			if err != nil {
				return fmt.Errorf("failed to shutdown process: %w", err)
			}
			successMsg := style.SuccessStyle.Render(fmt.Sprintf("✓ Process %d (%s) has been successfully shut down.", pid, processName))
			fmt.Printf("\n%s\n", successMsg)
		} else {
			cancelMsg := lipgloss.NewStyle().
				Foreground(style.SubtleColor).
				Render("Cancelled.")
			fmt.Printf("\n%s\n", cancelMsg)
		}

		return nil
	},
	}
}
