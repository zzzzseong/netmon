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
	registry.Register(commands.NewIPCommand())
	registry.Register(commands.NewShutdownCommand())
	registry.Register(commands.NewFindCommand())
	registry.Register(commands.NewHelpCommand(registry))

	// 명령어 인자 확인
	if len(os.Args) < 2 {
		commands.PrintUsage(registry)
		os.Exit(1)
	}

	commandName := os.Args[1]
	args := os.Args[2:]

	// help 별칭 처리 (--help, -h)
	if commandName == "--help" || commandName == "-h" {
		commands.PrintUsage(registry)
		os.Exit(0)
	}

	// 명령어 실행
	err := registry.Execute(commandName, args)
	if err != nil {
		errorMsg := style.ErrorStyle.Render(fmt.Sprintf("Error: %v", err))
		fmt.Fprintf(os.Stderr, "%s\n\n", errorMsg)
		commands.PrintUsage(registry)
		os.Exit(1)
	}
}

