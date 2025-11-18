package commands

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"netmon/style"
)

// HelpCommand는 도움말을 표시하는 명령어입니다
type HelpCommand struct {
	registry *Registry
}

// NewHelpCommand는 새로운 HelpCommand를 생성합니다
func NewHelpCommand(registry *Registry) *HelpCommand {
	return &HelpCommand{
		registry: registry,
	}
}

// Name은 명령어 이름을 반환합니다
func (c *HelpCommand) Name() string {
	return "help"
}

// Description은 명령어 설명을 반환합니다
func (c *HelpCommand) Description() string {
	return "Show help information"
}

// Usage는 명령어 사용법을 반환합니다
func (c *HelpCommand) Usage() string {
	return "help"
}

// Execute는 명령어를 실행합니다
func (c *HelpCommand) Execute(args []string) error {
	PrintUsage(c.registry)
	return nil
}

// PrintUsage는 사용법을 출력합니다
func PrintUsage(registry *Registry) {
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

	// Description
	descText := "A powerful CLI tool for monitoring network connections and managing processes."
	descStyle := lipgloss.NewStyle().
		Foreground(style.SubtleColor).
		Italic(true).
		MarginTop(1).
		MarginBottom(2)

	// ASCII Art 출력
	fmt.Print(asciiStyle.Render(asciiArt))
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

	// 명령어 목록을 박스로 감싸기
	var commandsList string
	for i, cmd := range registry.List() {
		cmdLine := style.CommandStyle.Render(fmt.Sprintf("  netmon %s", cmd.Usage()))
		desc := style.DescStyle.Render(cmd.Description())
		commandsList += fmt.Sprintf("%-32s %s", cmdLine, desc)
		if i < len(registry.List())-1 {
			commandsList += "\n"
		}
	}

	commandsBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(style.SecondaryColor).
		Padding(1, 2).
		Render(commandsList)
	fmt.Println(commandsBox)

	// Footer
	footerText := "For more information, visit: https://github.com/zzzzseong/netmon"
	footerStyle := lipgloss.NewStyle().
		Foreground(style.SubtleColor).
		Italic(true).
		MarginTop(2)
	fmt.Println(footerStyle.Render(footerText))
}

