package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"netmon/formatter"
	"netmon/provider"
)

// newRouteCmd creates and returns the route command.
// It displays system routing information with smart filtering.
func newRouteCmd() *cobra.Command {
	return &cobra.Command{
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

			// 포맷터를 사용하여 테이블 형식 출력
			fmtter := formatter.NewRouteTableFormatter()
			table := fmtter.Format(routes)
			fmt.Println(table)

			return nil
		},
	}
}
