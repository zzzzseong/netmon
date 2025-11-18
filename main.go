package main

import (
	"fmt"
	"os"
	"strconv"
	"text/tabwriter"

	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
	"github.com/shirou/gopsutil/v3/net"
	"github.com/shirou/gopsutil/v3/process"
)

func main() {
	// 명령어 인자 확인
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "ls":
		listPorts()
	case "kill":
		if len(os.Args) < 3 {
			fmt.Fprintf(os.Stderr, "Error: PID가 필요합니다.\n")
			fmt.Fprintf(os.Stderr, "사용법: netmon kill <pid>\n")
			os.Exit(1)
		}
		pidStr := os.Args[2]
		pid, err := strconv.Atoi(pidStr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: 유효하지 않은 PID: %s\n", pidStr)
			os.Exit(1)
		}
		killProcess(pid)
	default:
		fmt.Fprintf(os.Stderr, "Error: 알 수 없는 명령어: %s\n", command)
		printUsage()
		os.Exit(1)
	}
}

// 사용법 출력
func printUsage() {
	fmt.Println("사용법:")
	fmt.Println("  netmon ls              - 활성 포트 목록 표시")
	fmt.Println("  netmon kill <pid>      - 프로세스 종료")
}

// 연결 타입을 문자열로 변환 (1=TCP, 2=UDP)
func connectionTypeToString(connType uint32) string {
	switch connType {
	case 1:
		return "TCP"
	case 2:
		return "UDP"
	default:
		return "UNKNOWN"
	}
}

// 연결 타입이 UDP인지 확인
func isUDP(connType uint32) bool {
	return connType == 2
}

// 포트 목록 출력 함수
func listPorts() {
	// 모든 네트워크 연결 가져오기
	connections, err := net.Connections("inet")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: 네트워크 연결 정보를 가져올 수 없습니다: %v\n", err)
		os.Exit(1)
	}

	// 테이블 작성기 생성
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	
	// 헤더 출력
	fmt.Fprintln(w, "PROTOCOL\tLOCAL ADDRESS\tSTATUS\tPID\tPROCESS NAME")
	fmt.Fprintln(w, "--------\t------------\t------\t---\t------------")

	// 색상 정의
	green := color.New(color.FgGreen).SprintFunc()
	blue := color.New(color.FgBlue).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()

	// LISTEN 상태인 연결만 필터링 (UDP는 상태가 없을 수 있으므로 포트가 열려있으면 포함)
	listeningConns := make(map[string]net.ConnectionStat)
	for _, conn := range connections {
		// TCP는 LISTEN 상태만, UDP는 포트가 열려있으면 모두 포함
		if conn.Status == "LISTEN" || (isUDP(conn.Type) && conn.Laddr.Port > 0) {
			key := fmt.Sprintf("%d:%d", conn.Type, conn.Laddr.Port)
			listeningConns[key] = conn
		}
	}

	// 포트별로 정리하여 출력
	for _, conn := range listeningConns {
		protocol := connectionTypeToString(conn.Type)
		localAddr := fmt.Sprintf("%s:%d", conn.Laddr.IP, conn.Laddr.Port)
		
		// 상태에 따른 색상 적용 (UDP는 상태가 없을 수 있음)
		var statusStr string
		if conn.Status == "" {
			statusStr = green("LISTEN")
		} else {
			switch conn.Status {
			case "LISTEN":
				statusStr = green(conn.Status)
			case "ESTABLISHED":
				statusStr = blue(conn.Status)
			default:
				statusStr = yellow(conn.Status)
			}
		}

		// PID와 프로세스 이름 가져오기
		pid := int(conn.Pid)
		processName := "N/A"
		
		if pid > 0 {
			proc, err := process.NewProcess(int32(pid))
			if err == nil {
				name, err := proc.Name()
				if err == nil {
					processName = name
				}
			}
		}

		// 테이블 행 출력
		fmt.Fprintf(w, "%s\t%s\t%s\t%d\t%s\n", 
			protocol, localAddr, statusStr, pid, processName)
	}

	w.Flush()
}

// 프로세스 종료 함수
func killProcess(pid int) {
	// 프로세스 정보 가져오기
	proc, err := process.NewProcess(int32(pid))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: PID %d의 프로세스를 찾을 수 없습니다: %v\n", pid, err)
		os.Exit(1)
	}

	// 프로세스 이름 가져오기
	processName, err := proc.Name()
	if err != nil {
		processName = "N/A"
	}

	// 프로세스 상태 정보 가져오기
	status, err := proc.Status()
	if err != nil {
		status = []string{"N/A"}
	}

	// 프로세스가 사용하는 포트 가져오기
	connections, err := net.ConnectionsPid("inet", int32(pid))
	if err != nil {
		connections = []net.ConnectionStat{}
	}

	// 프로세스 정보 출력
	fmt.Println("\n=== 프로세스 정보 ===")
	fmt.Printf("PID:        %d\n", pid)
	fmt.Printf("이름:       %s\n", processName)
	fmt.Printf("상태:       %v\n", status)
	
	if len(connections) > 0 {
		fmt.Println("\n사용 중인 포트:")
		for _, conn := range connections {
			// LISTEN 상태이거나 UDP인 경우 표시
			if conn.Status == "LISTEN" || (conn.Status == "" && isUDP(conn.Type)) {
				fmt.Printf("  - %s:%d (%s)\n", conn.Laddr.IP, conn.Laddr.Port, connectionTypeToString(conn.Type))
			}
		}
	}

	// 인터랙티브 확인 프롬프트
	fmt.Println("\n이 프로세스를 종료하시겠습니까?")
	
	prompt := promptui.Select{
		Label: "선택",
		Items: []string{"Kill", "Cancel"},
		Size:  2,
	}

	index, result, err := prompt.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: 프롬프트 오류: %v\n", err)
		os.Exit(1)
	}

	if index == 0 && result == "Kill" {
		// 프로세스 종료
		err = proc.Kill()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: 프로세스를 종료할 수 없습니다: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("\n✓ 프로세스 %d (%s)가 성공적으로 종료되었습니다.\n", pid, processName)
	} else {
		fmt.Println("\n취소되었습니다.")
	}
}

