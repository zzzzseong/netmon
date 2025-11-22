package commands

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"netmon/style"
)

// VersionCommand는 버전 정보를 표시하는 명령어입니다
type VersionCommand struct{}

// NewVersionCommand는 새로운 VersionCommand를 생성합니다
func NewVersionCommand() *VersionCommand {
	return &VersionCommand{}
}

// Name은 명령어 이름을 반환합니다
func (c *VersionCommand) Name() string {
	return "version"
}

// Description은 명령어 설명을 반환합니다
func (c *VersionCommand) Description() string {
	return "Show version information"
}

// Usage는 명령어 사용법을 반환합니다
func (c *VersionCommand) Usage() string {
	return "version"
}

// Execute는 명령어를 실행합니다
func (c *VersionCommand) Execute(args []string) error {
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

	return nil
}

