package commands

import (
	"fmt"
	"os"
	"strconv"

	"github.com/manifoldco/promptui"
	"github.com/shirou/gopsutil/v3/net"
	"github.com/shirou/gopsutil/v3/process"
	"github.com/charmbracelet/lipgloss"
	"netmon/formatter"
	"netmon/style"
)

// KillCommand는 프로세스를 종료하는 명령어입니다
type KillCommand struct {
	formatter *formatter.ProcessInfoFormatter
}

// NewKillCommand는 새로운 KillCommand를 생성합니다
func NewKillCommand() *KillCommand {
	return &KillCommand{
		formatter: formatter.NewProcessInfoFormatter(),
	}
}

// Name은 명령어 이름을 반환합니다
func (c *KillCommand) Name() string {
	return "kill"
}

// Description은 명령어 설명을 반환합니다
func (c *KillCommand) Description() string {
	return "프로세스 종료"
}

// Usage는 명령어 사용법을 반환합니다
func (c *KillCommand) Usage() string {
	return "kill <pid>"
}

// Execute는 명령어를 실행합니다
func (c *KillCommand) Execute(args []string) error {
	if len(args) < 1 {
		errorMsg := style.ErrorStyle.Render("Error: PID가 필요합니다.")
		usageMsg := style.UsageStyle.Render("사용법: netmon kill <pid>")
		fmt.Fprintf(os.Stderr, "%s\n%s\n", errorMsg, usageMsg)
		os.Exit(1)
	}

	pidStr := args[0]
	pid, err := strconv.Atoi(pidStr)
	if err != nil {
		errorMsg := style.ErrorStyle.Render(fmt.Sprintf("Error: 유효하지 않은 PID: %s", pidStr))
		fmt.Fprintf(os.Stderr, "%s\n", errorMsg)
		os.Exit(1)
	}

	// 프로세스 정보 가져오기
	proc, err := process.NewProcess(int32(pid))
	if err != nil {
		errorMsg := style.ErrorStyle.Render(fmt.Sprintf("Error: PID %d의 프로세스를 찾을 수 없습니다: %v", pid, err))
		fmt.Fprintf(os.Stderr, "%s\n", errorMsg)
		os.Exit(1)
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
	info := c.formatter.Format(pid, processName, status, connections)
	fmt.Println(info)

	// 인터랙티브 확인 프롬프트
	promptTitle := lipgloss.NewStyle().
		Foreground(style.WarningColor).
		Bold(true).
		MarginTop(1).
		Render("⚠️  이 프로세스를 종료하시겠습니까?")
	fmt.Println(promptTitle)

	prompt := promptui.Select{
		Label: "",
		Items: []string{"Kill", "Cancel"},
		Size:  2,
	}

	index, result, err := prompt.Run()
	if err != nil {
		errorMsg := style.ErrorStyle.Render(fmt.Sprintf("Error: 프롬프트 오류: %v", err))
		fmt.Fprintf(os.Stderr, "%s\n", errorMsg)
		os.Exit(1)
	}

	if index == 0 && result == "Kill" {
		// 프로세스 종료
		err = proc.Kill()
		if err != nil {
			errorMsg := style.ErrorStyle.Render(fmt.Sprintf("Error: 프로세스를 종료할 수 없습니다: %v", err))
			fmt.Fprintf(os.Stderr, "\n%s\n", errorMsg)
			os.Exit(1)
		}
		successMsg := style.SuccessStyle.Render(fmt.Sprintf("✓ 프로세스 %d (%s)가 성공적으로 종료되었습니다.", pid, processName))
		fmt.Printf("\n%s\n", successMsg)
	} else {
		cancelMsg := lipgloss.NewStyle().
			Foreground(style.SubtleColor).
			Render("취소되었습니다.")
		fmt.Printf("\n%s\n", cancelMsg)
	}

	return nil
}

