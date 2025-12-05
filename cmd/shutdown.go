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

		// 포맷터를 사용하여 프로세스 정보 출력 (경고 헤더 포함)
		fmtter := formatter.NewProcessInfoFormatter()
		header := fmt.Sprintf("⚠️  Shutdown Confirmation for PID %d", pid)
		info := fmtter.Format(header, pid, processName, status, connections)
		fmt.Println(info)

		// 인터랙티브 확인 프롬프트
		promptTitle := lipgloss.NewStyle().
			Foreground(style.WarningColor).
			Bold(true).
			MarginTop(1).
			Render("Are you sure you want to shutdown this process?")
		fmt.Println(promptTitle)

		// 커스텀 템플릿으로 ?: 제거
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
			HideHelp:  true, // "Use the arrow keys" 도움말 숨기기
		}

		index, result, err := prompt.Run()
		if err != nil {
			return fmt.Errorf("prompt error: %w", err)
		}

		if index == 0 && result == "✓ Yes, shutdown" {
			// 프로세스 종료
			err = proc.Kill()
			if err != nil {
				return fmt.Errorf("failed to shutdown process: %w", err)
			}
			
			// 성공 메시지 박스
			successBox := lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(style.SuccessColor).
				Padding(1, 2).
				Foreground(style.SuccessColor).
				Bold(true).
				Render(fmt.Sprintf("✓ Process %d (%s) has been successfully shut down.", pid, processName))
			fmt.Printf("\n%s\n", successBox)
		} else {
			// 취소 메시지 박스
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
