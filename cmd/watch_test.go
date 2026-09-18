package cmd

import (
	"fmt"
	"testing"
)

func fakeTable(dataRows int) []string {
	lines := []string{"╭───╮", "│ H │", "├───┤"}
	for i := 0; i < dataRows; i++ {
		lines = append(lines, fmt.Sprintf("│ %d │", i))
	}
	return append(lines, "╰───╯")
}

func TestPaginateWatch_TableKeepsChromeOnEveryPage(t *testing.T) {
	lines := fakeTable(10)
	pages := paginateWatch(lines, 8) // 8 - tableChrome(4) = 4 data rows per page

	if len(pages) != 3 {
		t.Fatalf("expected 3 pages, got %d", len(pages))
	}
	for i, p := range pages {
		if p[0] != "╭───╮" || p[1] != "│ H │" || p[2] != "├───┤" || p[len(p)-1] != "╰───╯" {
			t.Errorf("page %d is missing table chrome: %v", i, p)
		}
		if len(p) > 8 {
			t.Errorf("page %d has %d lines, exceeds available 8", i, len(p))
		}
	}
	if pages[2][3] != "│ 8 │" || pages[2][4] != "│ 9 │" {
		t.Errorf("last page should hold rows 8 and 9, got %v", pages[2])
	}
}

func TestPaginateWatch_PlainTextSplitsByLine(t *testing.T) {
	lines := make([]string, 10)
	for i := range lines {
		lines[i] = fmt.Sprintf("line %d", i)
	}
	pages := paginateWatch(lines, 4)
	if len(pages) != 3 || len(pages[0]) != 4 || len(pages[2]) != 2 {
		t.Fatalf("expected pages of 4,4,2 lines, got %d pages: %v", len(pages), pages)
	}
}

func TestPaginateWatch_TinyTerminalStillProgresses(t *testing.T) {
	lines := fakeTable(3)
	pages := paginateWatch(lines, 1) // available smaller than chrome: one data row per page
	if len(pages) != 3 {
		t.Fatalf("expected 3 pages of one row each, got %d", len(pages))
	}
}

func TestIsTableOutput(t *testing.T) {
	if !isTableOutput(fakeTable(1)) {
		t.Error("expected table output to be detected")
	}
	if isTableOutput([]string{"Network Summary", "TCP: 3"}) {
		t.Error("plain text must not be detected as a table")
	}
	colored := fakeTable(1)
	for i := range colored {
		colored[i] = "\x1b[38;5;63m" + colored[i] + "\x1b[0m" // lipgloss border color on a TTY
	}
	if !isTableOutput(colored) {
		t.Error("table output with ANSI-colored borders must be detected")
	}
}
