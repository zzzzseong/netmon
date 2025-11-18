package style

import "github.com/charmbracelet/lipgloss"

// 색상 팔레트
var (
	PrimaryColor   = lipgloss.Color("63")   // 보라색
	SecondaryColor = lipgloss.Color("212")  // 핑크색
	SuccessColor   = lipgloss.Color("35")   // 초록색
	WarningColor   = lipgloss.Color("220")  // 노란색
	InfoColor      = lipgloss.Color("39")   // 파란색
	DangerColor    = lipgloss.Color("196")  // 빨간색
	SubtleColor    = lipgloss.Color("241")  // 회색
)

// 타이틀 스타일
var TitleStyle = lipgloss.NewStyle().
	Foreground(PrimaryColor).
	Bold(true).
	MarginBottom(1)

// 사용법 스타일
var UsageStyle = lipgloss.NewStyle().
	Foreground(SubtleColor).
	MarginTop(1)

// 명령어 스타일
var CommandStyle = lipgloss.NewStyle().
	Foreground(SecondaryColor).
	Bold(true)

// 설명 스타일
var DescStyle = lipgloss.NewStyle().
	Foreground(SubtleColor)

// 에러 스타일
var ErrorStyle = lipgloss.NewStyle().
	Foreground(DangerColor).
	Bold(true)

// 성공 스타일
var SuccessStyle = lipgloss.NewStyle().
	Foreground(SuccessColor).
	Bold(true)

// 정보 박스 스타일
var InfoBoxStyle = lipgloss.NewStyle().
	Border(lipgloss.RoundedBorder()).
	BorderForeground(PrimaryColor).
	Padding(1, 2).
	Margin(1, 0)

// 라벨 스타일
var LabelStyle = lipgloss.NewStyle().
	Foreground(SubtleColor).
	Width(12)

// 값 스타일
var ValueStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("252")).
	Bold(true)

// 헤더 스타일
var HeaderStyle = lipgloss.NewStyle().
	Foreground(PrimaryColor).
	Bold(true).
	Align(lipgloss.Center).
	Padding(0, 1)

