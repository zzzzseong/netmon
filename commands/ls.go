package commands

import (
	"fmt"
	"os"

	"github.com/shirou/gopsutil/v3/net"
	"netmon/formatter"
	"netmon/style"
	"netmon/utils"
)

// ListCommand는 포트 목록을 표시하는 명령어입니다
type ListCommand struct {
	formatter *formatter.PortTableFormatter
}

// NewListCommand는 새로운 ListCommand를 생성합니다
func NewListCommand() *ListCommand {
	return &ListCommand{
		formatter: formatter.NewPortTableFormatter(),
	}
}

// Name은 명령어 이름을 반환합니다
func (c *ListCommand) Name() string {
	return "ls"
}

// Description은 명령어 설명을 반환합니다
func (c *ListCommand) Description() string {
	return "활성 포트 목록 표시"
}

// Usage는 명령어 사용법을 반환합니다
func (c *ListCommand) Usage() string {
	return "ls"
}

// Execute는 명령어를 실행합니다
func (c *ListCommand) Execute(args []string) error {
	// 모든 네트워크 연결 가져오기
	connections, err := net.Connections("inet")
	if err != nil {
		errorMsg := style.ErrorStyle.Render(fmt.Sprintf("Error: 네트워크 연결 정보를 가져올 수 없습니다: %v", err))
		fmt.Fprintf(os.Stderr, "%s\n", errorMsg)
		os.Exit(1)
	}

	// LISTEN 상태인 연결만 필터링
	listeningConns := utils.FilterListeningConnections(connections)

	// 포맷터를 사용하여 테이블 생성
	table := c.formatter.Format(listeningConns)
	fmt.Println(table)

	return nil
}

