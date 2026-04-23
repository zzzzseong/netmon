package formatter

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"netmon/style"
)

// TraceHop contains information about a single hop in a traceroute.
type TraceHop struct {
	Hop    int    // Hop number
	Host   string // Hostname or IP address
	IP     string // IP address
	RTT1   string // Round-trip time for first probe
	RTT2   string // Round-trip time for second probe
	RTT3   string // Round-trip time for third probe
	Status string // Status: "success", "timeout", or "error"
}

// TracerouteFormatter formats traceroute output for display.
type TracerouteFormatter struct{}

// NewTracerouteFormatter creates a new TracerouteFormatter instance.
func NewTracerouteFormatter() *TracerouteFormatter {
	return &TracerouteFormatter{}
}

// PrintTableHeader prints the table header for traceroute output.
func (f *TracerouteFormatter) PrintTableHeader() {
	headerStyle := style.HeaderStyle

	// 헤더 생성
	headers := make([]string, len(TracerouteTableColumns))
	for i, col := range TracerouteTableColumns {
		headers[i] = headerStyle.Width(col.Width).Render(col.Title)
	}

	border := style.TableBorderStyle.Render("─")

	// 헤더 출력
	fmt.Printf("%s  %s  %s  %s  %s\n", headers[0], headers[1], headers[2], headers[3], headers[4])

	// 구분선 출력
	fmt.Printf("%s  %s  %s  %s  %s\n",
		strings.Repeat(border, TracerouteTableHopWidth),
		strings.Repeat(border, TracerouteTableHostWidth),
		strings.Repeat(border, TracerouteTableRTT1Width),
		strings.Repeat(border, TracerouteTableRTT2Width),
		strings.Repeat(border, TracerouteTableRTT3Width))
}

// PrintHopLine prints a single hop line with formatted information.
func (f *TracerouteFormatter) PrintHopLine(hop TraceHop) {
	// Hop 번호
	hopStr := fmt.Sprintf("%-6d", hop.Hop)
	hopStr = style.HopStyle.Render(hopStr)

	// Host/IP
	hostStr := hop.Host
	if hostStr == "" {
		hostStr = "*"
	}
	if hop.Status == "timeout" {
		hostStr = "Request timed out"
		hostStr = style.HostTimeoutStyle.Render(fmt.Sprintf("%-40s", hostStr))
	} else {
		hostStr = style.HostStyle.Render(fmt.Sprintf("%-40s", hostStr))
	}

	// RTT 값들
	rtt1Str := f.formatRTT(hop.RTT1)
	rtt2Str := f.formatRTT(hop.RTT2)
	rtt3Str := f.formatRTT(hop.RTT3)

	fmt.Printf("%s  %s  %s  %s  %s\n", hopStr, hostStr, rtt1Str, rtt2Str, rtt3Str)
}

// formatRTT formats an RTT value with color coding based on latency.
// Values less than 30ms are green, 30-100ms are yellow, and above 100ms are red.
func (f *TracerouteFormatter) formatRTT(rtt string) string {
	if rtt == "*" {
		return style.RTTStyle.Render(fmt.Sprintf("%-12s", "*"))
	}

	// RTT 값에서 숫자만 추출
	rttValue := strings.TrimSuffix(rtt, " ms")
	val, err := strconv.ParseFloat(rttValue, 64)
	if err != nil {
		return style.RTTStyle.Render(fmt.Sprintf("%-12s", rtt))
	}

	// RTT 값에 따라 색상 변경
	var color lipgloss.Color
	if val < 30 {
		color = style.SuccessColor
	} else if val < 100 {
		color = style.WarningColor
	} else {
		color = style.DangerColor
	}

	return lipgloss.NewStyle().Foreground(color).Render(fmt.Sprintf("%-12s", rtt))
}
