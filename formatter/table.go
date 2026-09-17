package formatter

import (
	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/charmbracelet/x/ansi"
	"golang.org/x/term"
	"netmon/style"
)

// TableColumn defines a table column with its preferred width and title.
type TableColumn struct {
	Width int    // Preferred column width; shrinks proportionally on narrow terminals
	Title string // Column title
}

// minColumnWidth is the smallest width a column can shrink to on a narrow terminal.
const minColumnWidth = 4

// terminalWidth returns the current terminal width and true when stdout is a terminal.
func terminalWidth() (int, bool) {
	w, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil || w <= 0 {
		return 0, false
	}
	return w, true
}

// fitWidth caps preferred to the terminal width when stdout is a terminal.
func fitWidth(preferred int) int {
	if w, ok := terminalWidth(); ok && w < preferred {
		return w
	}
	return preferred
}

// tableChromeWidth is the horizontal space a table needs beyond its column widths:
// one cell padding on each side of every column plus one border per column boundary.
func tableChromeWidth(columnCount int) int {
	return 3*columnCount + 1
}

// fitColumnWidths returns the column widths to use for the current terminal.
// Without a terminal (pipe, tests) the preferred widths are used unchanged.
func fitColumnWidths(columns []TableColumn) []int {
	termWidth, ok := terminalWidth()
	if !ok {
		termWidth = 0
	}
	return scaleColumnWidths(columns, termWidth)
}

// scaleColumnWidths returns column widths that fit a terminal of termWidth
// columns (0 means unlimited). When the table would be too wide, every column
// first keeps enough room for its title and the remaining space is shared in
// proportion to the preferred widths. If even the titles do not fit, all
// columns shrink proportionally instead.
func scaleColumnWidths(columns []TableColumn, termWidth int) []int {
	widths := make([]int, len(columns))
	sum := 0
	for i, col := range columns {
		widths[i] = col.Width
		sum += col.Width
	}

	if termWidth <= 0 || sum+tableChromeWidth(len(columns)) <= termWidth {
		return widths
	}
	available := termWidth - tableChromeWidth(len(columns))

	// Minimum per column: the title plus HeaderStyle's one-cell padding on each side.
	mins := make([]int, len(columns))
	minSum, flexSum := 0, 0
	for i, col := range columns {
		mins[i] = min(col.Width, lipgloss.Width(col.Title)+2)
		minSum += mins[i]
		flexSum += col.Width - mins[i]
	}

	if minSum <= available {
		surplus := available - minSum
		for i, col := range columns {
			widths[i] = mins[i]
			if flexSum > 0 {
				widths[i] += surplus * (col.Width - mins[i]) / flexSum
			}
		}
		return widths
	}

	for i := range widths {
		widths[i] = max(widths[i]*available/sum, minColumnWidth)
	}
	return widths
}

// CreateTable creates a formatted table with the given rows and columns.
// Column widths are honoured exactly: headers and cells are truncated with an
// ellipsis to their column so nothing wraps, every row stays on one line, and
// the table shrinks proportionally on narrow terminals.
func CreateTable(rows [][]string, columns []TableColumn) string {
	widths := fitColumnWidths(columns)

	headers := make([]string, len(columns))
	total := tableChromeWidth(len(columns))
	for i, col := range columns {
		// HeaderStyle pads one cell on each side, so the title itself gets width-2.
		title := ansi.Truncate(col.Title, widths[i]-2, "…")
		headers[i] = style.HeaderStyle.Width(widths[i]).Render(title)
		total += widths[i]
	}

	fitted := make([][]string, len(rows))
	for r, row := range rows {
		fitted[r] = make([]string, len(row))
		for c, cell := range row {
			if c < len(widths) {
				cell = ansi.Truncate(cell, widths[c], "…")
			}
			fitted[r][c] = cell
		}
	}

	t := table.New().
		Border(lipgloss.RoundedBorder()).
		BorderStyle(style.TableBorderStyle).
		StyleFunc(GetTableRowStyle).
		Headers(headers...).
		Rows(fitted...).
		Width(total).
		Wrap(false)

	return t.String()
}

// GetTableRowStyle returns the style for a table row based on its index.
// Even rows use TableRowEvenColor, odd rows use TableRowOddColor.
func GetTableRowStyle(row, col int) lipgloss.Style {
	color := style.TableRowEvenColor
	if row%2 != 0 {
		color = style.TableRowOddColor
	}
	return lipgloss.NewStyle().
		Foreground(color).
		Align(lipgloss.Left).
		Padding(0, 1)
}
