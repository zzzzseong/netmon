package commands

import (
	"fmt"
	"os"
	"sort"
	"strconv"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"netmon/provider"
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

// RouteEntry는 provider.RouteEntry의 별칭입니다 (하위 호환성 유지)
type RouteEntry = provider.RouteEntry

// Execute는 명령어를 실행합니다
func (c *RouteCommand) Execute(args []string) error {
	// RouteProvider를 통해 라우팅 테이블 가져오기
	routeProvider := provider.NewRouteProvider()
	routes, err := routeProvider.GetRoutes()
	if err != nil {
		errorMsg := style.ErrorStyle.Render(fmt.Sprintf("Error: Failed to get routing table information: %v", err))
		fmt.Fprintf(os.Stderr, "%s\n", errorMsg)
		os.Exit(1)
	}

	// 불필요한 라우트 필터링 (단일 호스트, multicast, link-local 등 제외)
	routes = provider.FilterRoutes(routes)

	// 테이블 형식 출력
	table := formatRouteTable(routes)
	fmt.Println(table)

	return nil
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
		metricStr := "-"
		if route.Metric > 0 {
			metricStr = strconv.Itoa(route.Metric)
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
