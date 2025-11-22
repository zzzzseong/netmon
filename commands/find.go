package commands

import (
	"fmt"
	"os"
	"strconv"

	"github.com/charmbracelet/lipgloss"
	"github.com/shirou/gopsutil/v3/net"
	"netmon/formatter"
	"netmon/style"
	"netmon/utils"
)

// FindCommand는 특정 포트를 사용하는 프로세스를 찾는 명령어입니다
type FindCommand struct {
	formatter *formatter.PortTableFormatter
}

// NewFindCommand는 새로운 FindCommand를 생성합니다
func NewFindCommand() *FindCommand {
	return &FindCommand{
		formatter: formatter.NewPortTableFormatter(),
	}
}

// Name은 명령어 이름을 반환합니다
func (c *FindCommand) Name() string {
	return "find"
}

// Description은 명령어 설명을 반환합니다
func (c *FindCommand) Description() string {
	return "Find process using a specific port"
}

// Usage는 명령어 사용법을 반환합니다
func (c *FindCommand) Usage() string {
	return "find <port>"
}

// Execute는 명령어를 실행합니다
func (c *FindCommand) Execute(args []string) error {
	if len(args) < 1 {
		errorMsg := style.ErrorStyle.Render("Error: Port number is required.")
		usageMsg := style.UsageStyle.Render("Usage: netmon find <port>")
		fmt.Fprintf(os.Stderr, "%s\n%s\n", errorMsg, usageMsg)
		os.Exit(1)
	}

	portStr := args[0]
	port, err := strconv.Atoi(portStr)
	if err != nil {
		errorMsg := style.ErrorStyle.Render(fmt.Sprintf("Error: Invalid port number: %s", portStr))
		fmt.Fprintf(os.Stderr, "%s\n", errorMsg)
		os.Exit(1)
	}

	if port < 1 || port > 65535 {
		errorMsg := style.ErrorStyle.Render("Error: Port number must be between 1 and 65535.")
		fmt.Fprintf(os.Stderr, "%s\n", errorMsg)
		os.Exit(1)
	}

	// 모든 네트워크 연결 가져오기
	connections, err := net.Connections("inet")
	if err != nil {
		errorMsg := style.ErrorStyle.Render(fmt.Sprintf("Error: Failed to get network connection information: %v", err))
		fmt.Fprintf(os.Stderr, "%s\n", errorMsg)
		os.Exit(1)
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
	table := c.formatter.Format(foundConns)
	fmt.Println(table)

	return nil
}

