package formatter

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"netmon/style"
)

// TableColumn defines a table column with its width and title.
type TableColumn struct {
	Width int    // Column width
	Title string // Column title
}

// CreateTable creates a formatted table with the given rows, columns, and total width.
// It applies consistent styling and returns the formatted table as a string.
func CreateTable(rows [][]string, columns []TableColumn, totalWidth int) string {
	// 헤더 생성
	headerStyle := style.HeaderStyle
	headerStyle = headerStyle.Align(lipgloss.Center)
	headers := make([]string, len(columns))
	for i, col := range columns {
		headers[i] = headerStyle.Width(col.Width).Render(col.Title)
	}

	// 테이블 생성
	t := table.New().
		Border(lipgloss.RoundedBorder()).
		BorderStyle(style.TableBorderStyle).
		StyleFunc(GetTableRowStyle).
		Headers(headers...).
		Rows(rows...).
		Width(totalWidth)

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
