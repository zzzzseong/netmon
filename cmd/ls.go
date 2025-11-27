package cmd

import (
	"fmt"

	"github.com/shirou/gopsutil/v3/net"
	"github.com/spf13/cobra"
	"netmon/formatter"
	"netmon/utils"
)

// lsCmd represents the ls command
var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List active ports",
	Long:  `Display all active listening ports with detailed process information.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// -a 옵션 값 가져오기
		includeUDP, _ := cmd.Flags().GetBool("all")

		// 포맷터 생성
		fmtter := formatter.NewPortTableFormatter()

		// 모든 네트워크 연결 가져오기
		connections, err := net.Connections("inet")
		if err != nil {
			return fmt.Errorf("failed to get network connection information: %w", err)
		}

		// LISTEN 상태인 연결만 필터링 (-a 옵션이 있을 때만 UDP 포함)
		listeningConns := utils.FilterListeningConnections(connections, includeUDP)

		// 포맷터를 사용하여 테이블 생성
		table := fmtter.Format(listeningConns)
		fmt.Println(table)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(lsCmd)

	// -a, --all 플래그 정의
	lsCmd.Flags().BoolP("all", "a", false, "Include UDP connections")
}
