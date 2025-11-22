package main

import (
	"fmt"
	"os"

	"netmon/commands"
	"netmon/style"
)

// Version 정보
const (
	Version   = "1.1.2"
	BuildDate = "2025-11-20"
)

func main() {
	// commands 패키지에 버전 정보 설정
	commands.Version = Version
	commands.BuildDate = BuildDate

	// 명령어 레지스트리 생성 및 등록 (순서 유지)
	registry := commands.NewRegistry()
	registry.Register(commands.NewListCommand())
	registry.Register(commands.NewIPCommand())
	registry.Register(commands.NewRouteCommand())
	registry.Register(commands.NewFindCommand())
	registry.Register(commands.NewShutdownCommand())
	registry.Register(commands.NewVersionCommand())
	registry.Register(commands.NewHelpCommand(registry))

	// 명령어 인자 확인
	if len(os.Args) < 2 {
		commands.PrintUsage(registry)
		os.Exit(1)
	}

	commandName := os.Args[1]
	args := os.Args[2:]

	// version 플래그 처리 (-v, --version)
	if commandName == "-v" || commandName == "--version" {
		err := registry.Execute("version", []string{})
		if err != nil {
			os.Exit(1)
		}
		os.Exit(0)
	}

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

