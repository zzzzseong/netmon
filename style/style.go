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

// 테이블 행 색상
var (
	TableRowEvenColor = lipgloss.Color("252") // 짝수 행 색상
	TableRowOddColor  = lipgloss.Color("248") // 홀수 행 색상
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

// 공통 스타일 (재사용)
var (
	// Destination 스타일
	DestinationStyle = lipgloss.NewStyle().Foreground(SecondaryColor).Bold(true)
	
	// Gateway 스타일
	GatewayStyle = lipgloss.NewStyle().Foreground(InfoColor)
	
	// Interface 스타일
	InterfaceStyle = lipgloss.NewStyle().Foreground(SubtleColor)
	
	// Metric 스타일
	MetricStyle = lipgloss.NewStyle().Foreground(SubtleColor)
	
	// Source 스타일
	SourceStyle = lipgloss.NewStyle().Foreground(SubtleColor)
	
	// IP Address 스타일
	IPAddressStyle = lipgloss.NewStyle().Foreground(InfoColor)
	
	// MAC Address 스타일
	MACAddressStyle = lipgloss.NewStyle().Foreground(SubtleColor)
	
	// Status UP 스타일
	StatusUPStyle = lipgloss.NewStyle().Foreground(SuccessColor).Bold(true)
	
	// Status DOWN 스타일
	StatusDOWNStyle = lipgloss.NewStyle().Foreground(WarningColor).Bold(true)
	
	// Status UNKNOWN 스타일
	StatusUNKNOWNStyle = lipgloss.NewStyle().Foreground(SubtleColor)
	
	// MTU 스타일
	MTUStyle = lipgloss.NewStyle().Foreground(SubtleColor)
	
	// Hop 스타일
	HopStyle = lipgloss.NewStyle().Foreground(SecondaryColor).Bold(true)
	
	// Host 스타일
	HostStyle = lipgloss.NewStyle().Foreground(InfoColor)
	
	// Host Timeout 스타일
	HostTimeoutStyle = lipgloss.NewStyle().Foreground(WarningColor)
	
	// RTT 스타일
	RTTStyle = lipgloss.NewStyle().Foreground(SubtleColor)
	
	// 테이블 테두리 스타일
	TableBorderStyle = lipgloss.NewStyle().Foreground(PrimaryColor)
	
	// ASCII Art 스타일 (help, version 명령어에서 재사용)
	ASCIIStyle = lipgloss.NewStyle().
		Foreground(PrimaryColor).
		Bold(true)
	
	// Version 텍스트 스타일
	VersionTextStyle = lipgloss.NewStyle().
		Foreground(SecondaryColor).
		Bold(true).
		MarginTop(1)
	
	// Build 텍스트 스타일
	BuildTextStyle = lipgloss.NewStyle().
		Foreground(SubtleColor).
		MarginBottom(1)
	
	// Footer 텍스트 스타일
	FooterTextStyle = lipgloss.NewStyle().
		Foreground(SubtleColor).
		Italic(true)
	
	// Description 텍스트 스타일
	DescTextStyle = lipgloss.NewStyle().
		Foreground(SubtleColor).
		Italic(true).
		MarginBottom(2)
)

// 테이블 너비 상수
const (
	TableWidthRoute      = 140
	TableWidthInterface  = 110
	TableWidthPort       = 130
	TableWidthTraceroute = 100
)

