package cmd

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
	"netmon/style"
)

// printCustomHelp prints a beautifully styled help message
func printCustomHelp(cmd *cobra.Command, args []string) {
	// ASCII Art
	asciiArt := `███╗   ██╗███████╗████████╗███╗   ███╗ ██████╗ ███╗   ██╗
████╗  ██║██╔════╝╚══██╔══╝████╗ ████║██╔═══██╗████╗  ██║
██╔██╗ ██║█████╗     ██║   ██╔████╔██║██║   ██║██╔██╗ ██║
██║╚██╗██║██╔══╝     ██║   ██║╚██╔╝██║██║   ██║██║╚██╗██║
██║ ╚████║███████╗   ██║   ██║ ╚═╝ ██║╚██████╔╝██║ ╚████║
╚═╝  ╚═══╝╚══════╝   ╚═╝   ╚═╝     ╚═╝ ╚═════╝ ╚═╝  ╚═══╝`

	// ASCII Art 스타일링
	asciiStyle := lipgloss.NewStyle().
		Foreground(style.PrimaryColor).
		Bold(true)

	// Version 정보
	versionText := fmt.Sprintf("Version %s", Version)
	versionStyle := lipgloss.NewStyle().
		Foreground(style.SecondaryColor).
		Bold(true).
		MarginTop(1)

	// Description
	descText := "A powerful CLI tool for monitoring network connections and managing processes."
	descStyle := lipgloss.NewStyle().
		Foreground(style.SubtleColor).
		Italic(true).
		MarginBottom(2)

	// ASCII Art 출력
	fmt.Print(asciiStyle.Render(asciiArt))
	fmt.Println(versionStyle.Render(versionText))
	fmt.Println(descStyle.Render(descText))

	// Usage 섹션
	usageHeader := lipgloss.NewStyle().
		Foreground(style.PrimaryColor).
		Bold(true).
		MarginBottom(1).
		Render("Usage:")
	fmt.Println(usageHeader)

	// Usage 내용 박스
	usageContent := "netmon <command> [arguments]"
	usageBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(style.PrimaryColor).
		Padding(1, 2).
		MarginBottom(2).
		Render(
			lipgloss.NewStyle().
				Foreground(lipgloss.Color("252")).
				Render(usageContent),
		)
	fmt.Println(usageBox)

	// Commands 섹션
	commandsHeader := lipgloss.NewStyle().
		Foreground(style.SecondaryColor).
		Bold(true).
		MarginBottom(1).
		Render("Commands:")
	fmt.Println(commandsHeader)

	// 명령어 목록 수집
	var commands []*cobra.Command
	for _, c := range cmd.Commands() {
		// completion과 help 명령어는 제외
		if c.Name() != "completion" && c.Name() != "help" {
			commands = append(commands, c)
		}
	}

	// 가장 긴 명령어 길이 계산
	maxUsageLen := 0
	for _, c := range commands {
		usageLen := len(fmt.Sprintf("  netmon %s", c.Use))
		if usageLen > maxUsageLen {
			maxUsageLen = usageLen
		}
	}

	// 명령어 목록 생성
	var commandsList strings.Builder
	for i, c := range commands {
		usage := fmt.Sprintf("  netmon %s", c.Use)
		// 패딩 추가
		padding := maxUsageLen - len(usage) + 4
		usage += strings.Repeat(" ", padding)

		cmdLine := style.CommandStyle.Render(usage)
		desc := style.DescStyle.Render(c.Short)
		commandsList.WriteString(cmdLine + desc)
		if i < len(commands)-1 {
			commandsList.WriteString("\n")
		}
	}

	commandsBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(style.SecondaryColor).
		Padding(1, 2).
		Render(commandsList.String())
	fmt.Println(commandsBox)

	// Footer
	footerText := "For more information, visit: https://github.com/zzzzseong/netmon"
	footerStyle := lipgloss.NewStyle().
		Foreground(style.SubtleColor).
		Italic(true).
		MarginTop(2)
	fmt.Println(footerStyle.Render(footerText))
}

// printCustomUsage prints a custom usage message
func printCustomUsage(cmd *cobra.Command) error {
	usageText := fmt.Sprintf("Usage: %s [command] [arguments]", cmd.Use)
	usageStyle := lipgloss.NewStyle().
		Foreground(style.PrimaryColor).
		Bold(true)
	fmt.Println(usageStyle.Render(usageText))
	return nil
}

