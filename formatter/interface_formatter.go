package formatter

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/shirou/gopsutil/v3/net"
	stdnet "net"
	"netmon/style"
)

// InterfaceTableFormatter formats network interface information.
type InterfaceTableFormatter struct{}

// NewInterfaceTableFormatter creates a new InterfaceTableFormatter instance.
func NewInterfaceTableFormatter() *InterfaceTableFormatter {
	return &InterfaceTableFormatter{}
}

// Format formats network interface information as a table.
// If showAll is true, IPv6 addresses are included; otherwise only IPv4 addresses are shown.
// Returns a formatted table string.
func (f *InterfaceTableFormatter) Format(interfaces []net.InterfaceStat, showAll bool) string {
	// 메모리 사전 할당 최적화
	rows := make([][]string, 0, len(interfaces))

	for _, iface := range interfaces {
		// 인터페이스 이름
		nameStr := style.DestinationStyle.Render(iface.Name)

		// MAC 주소
		macStr := iface.HardwareAddr
		if macStr == "" {
			macStr = "N/A"
		}
		macStr = style.MACAddressStyle.Render(macStr)

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
		ipStr = style.IPAddressStyle.Render(ipStr)

		// 상태 (UP/DOWN)
		var statusStr string
		if len(iface.Flags) > 0 {
			// flags에서 UP 상태 확인
			flagsStr := strings.Join(iface.Flags, ",")
			if strings.Contains(flagsStr, "up") || strings.Contains(flagsStr, "UP") {
				statusStr = style.StatusUPStyle.Render("UP")
			} else {
				statusStr = style.StatusDOWNStyle.Render("DOWN")
			}
		} else {
			statusStr = style.StatusUNKNOWNStyle.Render("UNKNOWN")
		}

		// MTU
		mtuStr := fmt.Sprintf("%d", iface.MTU)
		mtuStr = style.MTUStyle.Render(mtuStr)

		rows = append(rows, []string{
			nameStr,
			ipStr,
			macStr,
			statusStr,
			mtuStr,
		})
	}

	// 헤더 생성
	headerStyle := style.HeaderStyle.Copy().Align(lipgloss.Center)
	styledHeaders := make([]string, len(InterfaceTableColumns))
	for i, col := range InterfaceTableColumns {
		styledHeaders[i] = headerStyle.Width(col.Width).Render(col.Title)
	}

	// 테이블 생성 및 스타일링
	t := table.New().
		Border(lipgloss.RoundedBorder()).
		BorderStyle(style.TableBorderStyle).
		StyleFunc(GetTableRowStyle).
		Headers(styledHeaders...).
		Rows(rows...).
		Width(style.TableWidthInterface)

	return t.String()
}
