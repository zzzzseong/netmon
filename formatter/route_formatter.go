package formatter

import (
	"sort"
	"strconv"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"netmon/provider"
	"netmon/style"
)

// RouteTableFormatter formats routing table information.
type RouteTableFormatter struct{}

// NewRouteTableFormatter creates a new RouteTableFormatter instance.
func NewRouteTableFormatter() *RouteTableFormatter {
	return &RouteTableFormatter{}
}

// Format formats routing table entries as a table.
// Routes are sorted with default routes first, followed by other routes.
// Returns a formatted table string.
func (f *RouteTableFormatter) Format(routes []provider.RouteEntry) string {
	// 메모리 사전 할당 최적화
	rows := make([][]string, 0, len(routes))

	// 기본 게이트웨이를 먼저 표시하기 위해 정렬
	sort.Slice(routes, func(i, j int) bool {
		destI, destJ := routes[i].Destination, routes[j].Destination
		// 기본 게이트웨이를 먼저
		isDefaultI := destI == "default" || destI == "0.0.0.0/0"
		isDefaultJ := destJ == "default" || destJ == "0.0.0.0/0"

		if isDefaultI {
			return true
		}
		if isDefaultJ {
			return false
		}
		// 그 다음 목적지로 정렬
		return destI < destJ
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
		destStr = style.DestinationStyle.Render(destStr)

		// 게이트웨이 (Gateway)
		gatewayStr := route.Gateway
		if gatewayStr == "" {
			gatewayStr = "-"
		}
		gatewayStr = style.GatewayStyle.Render(gatewayStr)

		// 인터페이스 (Interface)
		ifaceStr := route.Interface
		if ifaceStr == "" {
			ifaceStr = "-"
		}
		ifaceStr = style.InterfaceStyle.Render(ifaceStr)

		// 메트릭 (Metric)
		metricStr := "-"
		if route.Metric > 0 {
			metricStr = strconv.Itoa(route.Metric)
		}
		metricStr = style.MetricStyle.Render(metricStr)

		// Source (소스 주소)
		sourceStr := route.Source
		if sourceStr == "" {
			sourceStr = "-"
		}
		sourceStr = style.SourceStyle.Render(sourceStr)

		rows = append(rows, []string{
			destStr,
			gatewayStr,
			ifaceStr,
			metricStr,
			sourceStr,
		})
	}

	// 헤더 스타일 적용
	headerStyle := style.HeaderStyle.Copy()
	styledHeaders := []string{
		headerStyle.Render("DESTINATION"),
		headerStyle.Render("GATEWAY"),
		headerStyle.Render("INTERFACE"),
		headerStyle.Render("METRIC"),
		headerStyle.Render("SOURCE"),
	}

	// 테이블 생성 및 스타일링
	t := table.New().
		Border(lipgloss.RoundedBorder()).
		BorderStyle(style.TableBorderStyle).
		StyleFunc(GetTableRowStyle).
		Headers(styledHeaders...).
		Rows(rows...).
		Width(style.TableWidthRoute)

	return t.String()
}
