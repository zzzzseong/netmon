package cmd

import (
	"fmt"
	stdnet "net"
	"sort"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/shirou/gopsutil/v3/net"
	"github.com/spf13/cobra"
	"netmon/style"
)

// ipCmd represents the ip command
var ipCmd = &cobra.Command{
	Use:   "ip",
	Short: "Show network interface information",
	Long:  `Display network interfaces with IP addresses.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// -a 플래그 확인
		showAll, _ := cmd.Flags().GetBool("all")

		// 네트워크 인터페이스 정보 가져오기
		interfaces, err := net.Interfaces()
		if err != nil {
			return fmt.Errorf("failed to get network interface information: %w", err)
		}

		// 인터페이스 이름으로 정렬
		sort.Slice(interfaces, func(i, j int) bool {
			return interfaces[i].Name < interfaces[j].Name
		})

		// 테이블 생성
		table := formatInterfaceTable(interfaces, showAll)
		fmt.Println(table)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(ipCmd)

	// -a, --all 플래그 정의
	ipCmd.Flags().BoolP("all", "a", false, "Show all addresses including IPv6")
}

// formatInterfaceTable은 네트워크 인터페이스 정보를 테이블 형식으로 포맷팅합니다
func formatInterfaceTable(interfaces []net.InterfaceStat, showAll bool) string {
	// 메모리 사전 할당 최적화
	rows := make([][]string, 0, len(interfaces))

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
