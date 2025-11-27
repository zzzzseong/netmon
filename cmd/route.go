package cmd

import (
	"fmt"
	"sort"
	"strconv"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/spf13/cobra"
	"netmon/provider"
	"netmon/style"
)

// routeCmd represents the route command
var routeCmd = &cobra.Command{
	Use:   "route",
	Short: "Show routing table information",
	Long:  `Display system routing information with smart filtering.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// RouteProvider를 통해 라우팅 테이블 가져오기
		routeProvider := provider.NewRouteProvider()
		routes, err := routeProvider.GetRoutes()
		if err != nil {
			return fmt.Errorf("failed to get routing table information: %w", err)
		}

		// 불필요한 라우트 필터링 (단일 호스트, multicast, link-local 등 제외)
		routes = provider.FilterRoutes(routes)

		// 테이블 형식 출력
		table := formatRouteTable(routes)
		fmt.Println(table)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(routeCmd)
}

// formatRouteTable은 라우팅 테이블 정보를 테이블 형식으로 포맷팅합니다
func formatRouteTable(routes []provider.RouteEntry) string {
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
