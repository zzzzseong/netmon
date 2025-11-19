package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/lipgloss"
	"netmon/commands"
	"netmon/style"
)

// printVersion 버전 정보 출력
func printVersion() {
	asciiArt := `███╗   ██╗███████╗████████╗███╗   ███╗ ██████╗ ███╗   ██╗
████╗  ██║██╔════╝╚══██╔══╝████╗ ████║██╔═══██╗████╗  ██║
██╔██╗ ██║█████╗     ██║   ██╔████╔██║██║   ██║██╔██╗ ██║
██║╚██╗██║██╔══╝     ██║   ██║╚██╔╝██║██║   ██║██║╚██╗██║
██║ ╚████║███████╗   ██║   ██║ ╚═╝ ██║╚██████╔╝██║ ╚████║
╚═╝  ╚═══╝╚══════╝   ╚═╝   ╚═╝     ╚═╝ ╚═════╝ ╚═╝  ╚═══╝`

	asciiStyle := lipgloss.NewStyle().
		Foreground(style.PrimaryColor).
		Bold(true)

	versionText := fmt.Sprintf("Version: %s", Version)
	versionStyle := lipgloss.NewStyle().
		Foreground(style.SecondaryColor).
		Bold(true).
		MarginTop(1)

	buildText := fmt.Sprintf("Build Date: %s", BuildDate)
	buildStyle := lipgloss.NewStyle().
		Foreground(style.SubtleColor).
		MarginBottom(1)

	fmt.Print(asciiStyle.Render(asciiArt))
	fmt.Println(versionStyle.Render(versionText))
	fmt.Println(buildStyle.Render(buildText))
	
	footerText := "For more information, visit: https://github.com/zzzzseong/netmon"
	footerStyle := lipgloss.NewStyle().
		Foreground(style.SubtleColor).
		Italic(true)
	fmt.Println(footerStyle.Render(footerText))
}

func main() {
	// 명령어 레지스트리 생성 및 등록 (순서 유지)
	registry := commands.NewRegistry()
	registry.Register(commands.NewListCommand())
	registry.Register(commands.NewIPCommand())
	registry.Register(commands.NewRouteCommand())
	registry.Register(commands.NewFindCommand())
	registry.Register(commands.NewShutdownCommand())
	registry.Register(commands.NewHelpCommand(registry))

	// 명령어 인자 확인
	if len(os.Args) < 2 {
		commands.PrintUsage(registry)
		os.Exit(1)
	}

	commandName := os.Args[1]
	args := os.Args[2:]

	// version 플래그 처리 (-v, --version, version)
	if commandName == "-v" || commandName == "--version" || commandName == "version" {
		printVersion()
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

