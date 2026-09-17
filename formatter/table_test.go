package formatter

import (
	"strings"
	"testing"
)

var sampleColumns = []TableColumn{
	{Width: 10, Title: "PROTOCOL"},
	{Width: 25, Title: "LOCAL ADDRESS"},
	{Width: 8, Title: "PID"},
}

func TestScaleColumnWidths_UnchangedWhenItFits(t *testing.T) {
	// 43 + chrome(10) = 53 <= 80, and 0 means no terminal at all
	for _, termWidth := range []int{0, 53, 80} {
		got := scaleColumnWidths(sampleColumns, termWidth)
		if got[0] != 10 || got[1] != 25 || got[2] != 8 {
			t.Errorf("termWidth=%d: expected preferred widths, got %v", termWidth, got)
		}
	}
}

func TestScaleColumnWidths_KeepsTitlesThenSharesSurplus(t *testing.T) {
	// available = 45 - 10 = 35; title minimums are 10, 15, 5 (=30); surplus 5 is
	// shared by flex (0, 10, 3) -> +0, +3, +1
	got := scaleColumnWidths(sampleColumns, 45)
	want := []int{10, 18, 6}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("expected %v, got %v", want, got)
		}
	}
	if sum(got)+tableChromeWidth(len(got)) > 45 {
		t.Fatalf("table wider than terminal: %v", got)
	}
}

func TestScaleColumnWidths_TooNarrowShrinksProportionally(t *testing.T) {
	// available = 20 - 10 = 10 < title minimums, so scale 10/25/8 by 10/43 with floor 4
	got := scaleColumnWidths(sampleColumns, 20)
	want := []int{4, 5, 4}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("expected %v, got %v", want, got)
		}
	}
}

func TestCreateTable_TruncatesInsteadOfWrapping(t *testing.T) {
	rows := [][]string{{"TCP", strings.Repeat("x", 60), "1"}}
	out := CreateTable(rows, sampleColumns)

	lines := strings.Split(strings.TrimRight(out, "\n"), "\n")
	if len(lines) != 5 {
		t.Fatalf("expected 5 lines (border, header, separator, row, border), got %d:\n%s", len(lines), out)
	}
	if !strings.Contains(lines[3], "…") {
		t.Fatalf("expected long cell to be truncated with an ellipsis, got:\n%s", lines[3])
	}
	if !strings.Contains(lines[1], "LOCAL ADDRESS") {
		t.Fatalf("expected full header title, got:\n%s", lines[1])
	}
}

func sum(xs []int) int {
	total := 0
	for _, x := range xs {
		total += x
	}
	return total
}
