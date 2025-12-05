package cmd

import (
	"fmt"
	"sort"

	"github.com/shirou/gopsutil/v3/net"
	"github.com/spf13/cobra"
	"netmon/formatter"
)

// newIPCmd creates and returns the ip command.
// It displays network interfaces with IP addresses.
func newIPCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ip [-a]",
		Short: "Show network interfaces (-a: include IPv6)",
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

			// 포맷터를 사용하여 테이블 생성
			fmtter := formatter.NewInterfaceTableFormatter()
			table := fmtter.Format(interfaces, showAll)
			fmt.Println(table)

			return nil
		},
	}

	// -a, --all 플래그 정의
	cmd.Flags().BoolP("all", "a", false, "Show all addresses including IPv6")

	return cmd
}
