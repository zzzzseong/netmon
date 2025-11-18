package main

import (
	"fmt"
	"os"
	
	"netmon/commands"
	"netmon/style"
)

func main() {
	// 명령어 레지스트리 생성 및 등록
	registry := commands.NewRegistry()
	registry.Register(commands.NewListCommand())
	registry.Register(commands.NewKillCommand())

	// 명령어 인자 확인
	if len(os.Args) < 2 {
		printUsage(registry)
		os.Exit(1)
	}

	commandName := os.Args[1]
	args := os.Args[2:]

	// 명령어 실행
	err := registry.Execute(commandName, args)
	if err != nil {
		errorMsg := style.ErrorStyle.Render(fmt.Sprintf("Error: %v", err))
		fmt.Fprintf(os.Stderr, "%s\n\n", errorMsg)
		printUsage(registry)
		os.Exit(1)
	}
}

// 사용법 출력
func printUsage(registry *commands.Registry) {
	title := style.TitleStyle.Render("📡 netmon - 네트워크 모니터링 도구")
	usage := style.UsageStyle.Render("사용법:")

	fmt.Println(title)
	fmt.Println(usage)

	for _, cmd := range registry.List() {
		cmdLine := style.CommandStyle.Render(fmt.Sprintf("  netmon %s", cmd.Usage()))
		desc := style.DescStyle.Render(fmt.Sprintf("- %s", cmd.Description()))
		fmt.Printf("%-30s %s\n", cmdLine, desc)
	}
}

