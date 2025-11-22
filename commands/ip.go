package commands

import (
	"fmt"
	stdnet "net"
	"os"
	"sort"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/shirou/gopsutil/v3/net"
	"netmon/style"
)

// IPCommand는 네트워크 인터페이스 정보를 표시하는 명령어입니다
type IPCommand struct{}

// NewIPCommand는 새로운 IPCommand를 생성합니다
func NewIPCommand() *IPCommand {
	return &IPCommand{}
}

// Name은 명령어 이름을 반환합니다
func (c *IPCommand) Name() string {
	return "ip"
}

// Description은 명령어 설명을 반환합니다
func (c *IPCommand) Description() string {
	return "Show network interface information"
}

// Usage는 명령어 사용법을 반환합니다
func (c *IPCommand) Usage() string {
	return "ip [-a]"
}

// Execute는 명령어를 실행합니다
func (c *IPCommand) Execute(args []string) error {
	// -a 플래그 확인
	showAll := false
	for _, arg := range args {
		if arg == "-a" {
			showAll = true
			break
		}
	}

	// 네트워크 인터페이스 정보 가져오기
	interfaces, err := net.Interfaces()
	if err != nil {
		errorMsg := style.ErrorStyle.Render(fmt.Sprintf("Error: Failed to get network interface information: %v", err))
		fmt.Fprintf(os.Stderr, "%s\n", errorMsg)
		os.Exit(1)
	}

	// 인터페이스 이름으로 정렬
	sort.Slice(interfaces, func(i, j int) bool {
		return interfaces[i].Name < interfaces[j].Name
	})

	// 테이블 생성
	table := formatInterfaceTable(interfaces, showAll)
	fmt.Println(table)

	return nil
}

// formatInterfaceTable은 네트워크 인터페이스 정보를 테이블 형식으로 포맷팅합니다
func formatInterfaceTable(interfaces []net.InterfaceStat, showAll bool) string {
	var rows [][]string

	for _, iface := range interfaces {
		// 인터페이스 이름
		nameStr := lipgloss.NewStyle().Foreground(style.SecondaryColor).Bold(true).Render(iface.Name)

		// MAC 주소
		macStr := iface.HardwareAddr
		if macStr == "" {
			macStr = "N/A"
		}
		macStr = lipgloss.NewStyle().Foreground(style.SubtleColor).Render(macStr)

		// IP 주소들
		var ipAddrs []string
		for _, addr := range iface.Addrs {
			// IP 주소와 서브넷 마스크 파싱
			ip, ipNet, err := stdnet.ParseCIDR(addr.Addr)
			if err != nil {
				continue
			}

			// IPv4와 IPv6 구분
			if ip.To4() != nil {
				// IPv4 - 항상 포함
				mask := ipNet.Mask
				ones, _ := mask.Size()
				ipAddrs = append(ipAddrs, fmt.Sprintf("%s/%d", ip.String(), ones))
			} else if showAll {
				// IPv6 - showAll이 true일 때만 포함
				ones, _ := ipNet.Mask.Size()
				ipAddrs = append(ipAddrs, fmt.Sprintf("%s/%d", ip.String(), ones))
			}
		}

		// IP 주소 문자열 생성
		ipStr := strings.Join(ipAddrs, ", ")
		if ipStr == "" {
			// showAll이 false면 IP 주소가 없는 인터페이스는 건너뛰기
			if !showAll {
				continue
			}
			ipStr = "N/A"
		}
		ipStr = lipgloss.NewStyle().Foreground(style.InfoColor).Render(ipStr)

		// 상태 (UP/DOWN)
		var statusStr string
		if len(iface.Flags) > 0 {
			// flags에서 UP 상태 확인
			flagsStr := strings.Join(iface.Flags, ",")
			if strings.Contains(flagsStr, "up") || strings.Contains(flagsStr, "UP") {
				statusStr = lipgloss.NewStyle().Foreground(style.SuccessColor).Bold(true).Render("UP")
			} else {
				statusStr = lipgloss.NewStyle().Foreground(style.WarningColor).Bold(true).Render("DOWN")
			}
		} else {
			statusStr = lipgloss.NewStyle().Foreground(style.SubtleColor).Render("UNKNOWN")
		}

		// MTU
		mtuStr := fmt.Sprintf("%d", iface.MTU)
		mtuStr = lipgloss.NewStyle().Foreground(style.SubtleColor).Render(mtuStr)

		rows = append(rows, []string{
			nameStr,
			ipStr,
			macStr,
			statusStr,
			mtuStr,
		})
	}

	// 헤더 스타일 적용 (가운데 정렬 및 너비 설정)
	headerStyle := style.HeaderStyle.Copy().Align(lipgloss.Center)
	styledHeaders := []string{
		headerStyle.Width(12).Render("INTERFACE"),
		headerStyle.Width(40).Render("IP ADDRESS"),
		headerStyle.Width(20).Render("MAC ADDRESS"),
		headerStyle.Width(8).Render("STATUS"),
		headerStyle.Width(6).Render("MTU"),
	}

	// 테이블 생성 및 스타일링
	t := table.New().
		Border(lipgloss.RoundedBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(style.PrimaryColor)).
		StyleFunc(func(row, col int) lipgloss.Style {
			// 모든 데이터 행에 왼쪽 정렬 적용
			if row%2 == 0 {
				return lipgloss.NewStyle().
					Foreground(lipgloss.Color("252")).
					Align(lipgloss.Left).
					Padding(0, 1)
			} else {
				return lipgloss.NewStyle().
					Foreground(lipgloss.Color("248")).
					Align(lipgloss.Left).
					Padding(0, 1)
			}
		}).
		Headers(styledHeaders...).
		Rows(rows...).
		Width(110)

	return t.String()
}

