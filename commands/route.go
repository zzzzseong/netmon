package commands

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"sort"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"netmon/style"
)

// RouteCommand는 라우팅 테이블 정보를 표시하는 명령어입니다
type RouteCommand struct{}

// NewRouteCommand는 새로운 RouteCommand를 생성합니다
func NewRouteCommand() *RouteCommand {
	return &RouteCommand{}
}

// Name은 명령어 이름을 반환합니다
func (c *RouteCommand) Name() string {
	return "route"
}

// Description은 명령어 설명을 반환합니다
func (c *RouteCommand) Description() string {
	return "Show routing table information"
}

// Usage는 명령어 사용법을 반환합니다
func (c *RouteCommand) Usage() string {
	return "route"
}

// RouteEntry는 라우팅 테이블 항목을 나타냅니다
type RouteEntry struct {
	Destination string
	Gateway     string
	Interface   string
	Metric      string
	Source      string
}

// Execute는 명령어를 실행합니다
func (c *RouteCommand) Execute(args []string) error {
	// 시스템 명령어로 라우팅 테이블 가져오기
	routes, err := getRoutes()
	if err != nil {
		errorMsg := style.ErrorStyle.Render(fmt.Sprintf("Error: Failed to get routing table information: %v", err))
		fmt.Fprintf(os.Stderr, "%s\n", errorMsg)
		os.Exit(1)
	}

	// 테이블 생성
	table := formatRouteTable(routes)
	fmt.Println(table)

	return nil
}

// getRoutes는 시스템 명령어를 사용하여 라우팅 테이블을 가져옵니다
func getRoutes() ([]RouteEntry, error) {
	var cmd *exec.Cmd
	if runtime.GOOS == "darwin" {
		// macOS: netstat -rn 사용
		cmd = exec.Command("netstat", "-rn")
	} else {
		// Linux: ip route show 사용
		cmd = exec.Command("ip", "route", "show")
	}

	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	return parseRoutes(string(output), runtime.GOOS == "darwin"), nil
}

// parseRoutes는 명령어 출력을 파싱하여 RouteEntry 슬라이스로 변환합니다
func parseRoutes(output string, isMacOS bool) []RouteEntry {
	var routes []RouteEntry
	scanner := bufio.NewScanner(strings.NewReader(output))

	// 첫 번째 줄(헤더) 건너뛰기
	firstLine := true

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		if firstLine {
			firstLine = false
			// macOS의 경우 헤더가 있으면 건너뛰기
			if isMacOS && (strings.Contains(line, "Destination") || strings.Contains(line, "Routing tables")) {
				continue
			}
			// Linux의 경우 헤더가 없으므로 바로 파싱
		}

		if isMacOS {
			route := parseMacOSRoute(line)
			if route != nil {
				routes = append(routes, *route)
			}
		} else {
			route := parseLinuxRoute(line)
			if route != nil {
				routes = append(routes, *route)
			}
		}
	}

	return routes
}

// parseMacOSRoute는 macOS netstat -rn 출력을 파싱하여 Linux ip route 스타일과 유사하게 변환합니다
func parseMacOSRoute(line string) *RouteEntry {
	fields := strings.Fields(line)
	if len(fields) < 4 {
		return nil
	}

	dest := fields[0]
	gateway := fields[1]
	iface := fields[3]

	if strings.EqualFold(dest, "destination") ||
		isIPv6Token(dest) || isIPv6Token(gateway) ||
		strings.Contains(gateway, "%") ||
		strings.HasPrefix(iface, "utun") ||
		strings.HasPrefix(iface, "llw") ||
		iface == "lo0" {
		return nil
	}

	if dest != "default" {
		if isSystemRouteDestination(dest) || isHostRouteDestination(dest) {
			return nil
		}
	}

	if dest == "default" {
		dest = "default"
	}

	normalizedGateway := normalizeGatewayToken(gateway)

	return &RouteEntry{
		Destination: dest,
		Gateway:     normalizedGateway,
		Interface:   iface,
		Metric:      "",
		Source:      "",
	}
}

// parseLinuxRoute는 Linux ip route show 출력을 파싱합니다
func parseLinuxRoute(line string) *RouteEntry {
	fields := strings.Fields(line)
	if len(fields) < 1 {
		return nil
	}

	route := &RouteEntry{
		Gateway:   "",
		Interface: "",
		Metric:    "",
		Source:    "",
	}

	if fields[0] == "default" {
		route.Destination = "default"
	}

	// ip route show 형식 파싱
	for i := 0; i < len(fields); i++ {
		switch fields[i] {
		case "default":
			route.Destination = "default"
		case "via":
			if i+1 < len(fields) {
				route.Gateway = fields[i+1]
				i++
			}
		case "dev":
			if i+1 < len(fields) {
				route.Interface = fields[i+1]
				i++
			}
		case "metric":
			if i+1 < len(fields) {
				route.Metric = fields[i+1]
				i++
			}
		case "src":
			if i+1 < len(fields) {
				route.Source = fields[i+1]
				i++
			}
		case "proto":
			if i+1 < len(fields) {
				i++
			}
		case "scope":
			if i+1 < len(fields) {
				i++
			}
		default:
			// 목적지 네트워크 (CIDR 형식)
			if strings.Contains(fields[i], "/") || strings.Contains(fields[i], ".") {
				if route.Destination == "" {
					route.Destination = fields[i]
				}
			}
		}
	}

	if route.Destination == "" {
		return nil
	}

	return route
}

// formatRouteTable은 라우팅 테이블 정보를 테이블 형식으로 포맷팅합니다
func formatRouteTable(routes []RouteEntry) string {
	var rows [][]string

	// 기본 게이트웨이를 먼저 표시하기 위해 정렬
	sort.Slice(routes, func(i, j int) bool {
		// 기본 게이트웨이를 먼저
		if routes[i].Destination == "default" || routes[i].Destination == "0.0.0.0/0" {
			return true
		}
		if routes[j].Destination == "default" || routes[j].Destination == "0.0.0.0/0" {
			return false
		}
		// 그 다음 목적지로 정렬
		return routes[i].Destination < routes[j].Destination
	})

	for _, route := range routes {
		// 목적지 (Destination)
		destStr := route.Destination
		if destStr == "" {
			destStr = "default"
		}
		if destStr == "0.0.0.0/0" {
			destStr = "default"
		}
		destStr = lipgloss.NewStyle().Foreground(style.SecondaryColor).Bold(true).Render(destStr)

		// 게이트웨이 (Gateway)
		gatewayStr := route.Gateway
		if gatewayStr == "" {
			gatewayStr = "-"
		}
		gatewayStr = lipgloss.NewStyle().Foreground(style.InfoColor).Render(gatewayStr)

		// 인터페이스 (Interface)
		ifaceStr := route.Interface
		if ifaceStr == "" {
			ifaceStr = "-"
		}
		ifaceStr = lipgloss.NewStyle().Foreground(style.SubtleColor).Render(ifaceStr)

		// 메트릭 (Metric)
		metricStr := route.Metric
		if metricStr == "" {
			metricStr = "-"
		}
		metricStr = lipgloss.NewStyle().Foreground(style.SubtleColor).Render(metricStr)

		// Source (소스 주소)
		sourceStr := route.Source
		if sourceStr == "" {
			sourceStr = "-"
		}
		sourceStr = lipgloss.NewStyle().Foreground(style.SubtleColor).Render(sourceStr)

		rows = append(rows, []string{
			destStr,
			gatewayStr,
			ifaceStr,
			metricStr,
			sourceStr,
		})
	}

	// 헤더 스타일 적용
	styledHeaders := []string{
		style.HeaderStyle.Render("DESTINATION"),
		style.HeaderStyle.Render("GATEWAY"),
		style.HeaderStyle.Render("INTERFACE"),
		style.HeaderStyle.Render("METRIC"),
		style.HeaderStyle.Render("SOURCE"),
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
		Width(140)

	return t.String()
}

func isIPv6Token(value string) bool {
	return strings.Contains(value, ":")
}

func isMacAddressToken(value string) bool {
	parts := strings.Split(value, ":")
	if len(parts) != 6 {
		return false
	}
	for _, part := range parts {
		if len(part) == 0 || len(part) > 2 {
			return false
		}
	}
	return true
}

func normalizeGatewayToken(gateway string) string {
	if gateway == "" {
		return ""
	}
	if strings.HasPrefix(gateway, "link#") || isMacAddressToken(gateway) || isIPv6Token(gateway) {
		return ""
	}
	return gateway
}

func isSystemRouteDestination(dest string) bool {
	base := dest
	if slash := strings.Index(base, "/"); slash != -1 {
		base = base[:slash]
	}
	prefixes := []string{
		"127.",
		"169.254.",
		"224.",
		"255.",
	}
	for _, prefix := range prefixes {
		if strings.HasPrefix(base, prefix) {
			return true
		}
		trimmed := strings.TrimSuffix(prefix, ".")
		if trimmed != prefix && base == trimmed {
			return true
		}
	}
	return false
}

func isHostRouteDestination(dest string) bool {
	if strings.HasSuffix(dest, "/32") {
		return true
	}
	if strings.Count(dest, ".") == 3 && !strings.Contains(dest, "/") {
		return true
	}
	return false
}
