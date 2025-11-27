package commands

import (
	"fmt"

	"github.com/shirou/gopsutil/v3/net"
	"netmon/formatter"
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
	return "List active ports"
}

// Usage는 명령어 사용법을 반환합니다
func (c *ListCommand) Usage() string {
	return "ls [-a]"
}

// Execute는 명령어를 실행합니다
func (c *ListCommand) Execute(args []string) error {
	// -a 옵션 파싱
	includeUDP := false
	for _, arg := range args {
		if arg == "-a" {
			includeUDP = true
			break
		}
	}

	// 모든 네트워크 연결 가져오기
	connections, err := net.Connections("inet")
	if err != nil {
		return fmt.Errorf("failed to get network connection information: %w", err)
	}

	// LISTEN 상태인 연결만 필터링 (-a 옵션이 있을 때만 UDP 포함)
	listeningConns := utils.FilterListeningConnections(connections, includeUDP)

	// 포맷터를 사용하여 테이블 생성
	table := c.formatter.Format(listeningConns)
	fmt.Println(table)

	return nil
}

