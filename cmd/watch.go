package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/charmbracelet/x/ansi"
	"golang.org/x/term"
)

const (
	watchHeaderLines = 2 // command info line + blank line
	watchFooterLines = 2 // blank line + pagination indicator
	tableChrome      = 4 // top border + column headers + separator + bottom border
)

// runWithWatch runs fn on every interval, paginating output that exceeds the
// terminal height. [ / ] and left/right arrow keys navigate pages; q or Ctrl+C exits.
// If a refresh fails, the previous output stays on screen and the error is shown in the header.
func runWithWatch(name string, interval time.Duration, fn func() (string, error)) error {
	if interval <= 0 {
		return fmt.Errorf("interval must be at least 1 second")
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	fmt.Print("\033[?1049h\033[?25l")
	defer fmt.Print("\033[?1049l\033[?25h")

	// Raw mode for single-keypress navigation (skipped if stdin is not a TTY)
	stdinFd := int(os.Stdin.Fd())
	if term.IsTerminal(stdinFd) {
		old, err := term.MakeRaw(stdinFd)
		if err == nil {
			defer term.Restore(stdinFd, old)
		}
	}

	navCh := make(chan int, 1)
	go readWatchKeys(navCh, cancel)

	page := 0
	content := ""
	var lastErr error

	refresh := func() {
		out, err := fn()
		lastErr = err
		if err == nil {
			content = out
		}
	}

	// Initial fetch
	refresh()
	fmt.Print("\033[H\033[2J")
	page = renderWatch(name, interval, content, lastErr, page)

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil
		case delta := <-navCh:
			page += delta
			if page < 0 {
				page = 0
			}
			page = renderWatch(name, interval, content, lastErr, page)
		case <-ticker.C:
			refresh()
			page = renderWatch(name, interval, content, lastErr, page)
		}
	}
}

func readWatchKeys(navCh chan<- int, cancel context.CancelFunc) {
	buf := make([]byte, 4)
	for {
		n, err := os.Stdin.Read(buf)
		if err != nil || n == 0 {
			return
		}
		b := buf[:n]
		switch {
		case n == 1 && b[0] == '[':
			select {
			case navCh <- -1:
			default:
			}
		case n == 1 && b[0] == ']':
			select {
			case navCh <- 1:
			default:
			}
		case n == 1 && (b[0] == 'q' || b[0] == 3): // q or Ctrl+C
			cancel()
			return
		case n >= 3 && b[0] == 27 && b[1] == '[': // ESC sequence
			switch b[2] {
			case 'D': // left arrow
				select {
				case navCh <- -1:
				default:
				}
			case 'C': // right arrow
				select {
				case navCh <- 1:
				default:
				}
			}
		}
	}
}

// renderWatch prints the current page and returns the (possibly clamped) page index.
func renderWatch(name string, interval time.Duration, content string, refreshErr error, page int) int {
	termWidth, termHeight, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil || termHeight <= 0 {
		termWidth, termHeight = 80, 24
	}

	header := fmt.Sprintf("netmon %s  —  every %.0fs  —  %s  —  [ ] navigate  q or Ctrl+C to stop",
		name, interval.Seconds(), time.Now().Format("2006-01-02 15:04:05"))
	if refreshErr != nil {
		header += "  —  ⚠ refresh failed: " + refreshErr.Error()
	}
	// Keep the header on one line so the pagination math below stays valid.
	header = ansi.Truncate(header, termWidth, "…")

	// The terminal is in raw mode, which disables output post-processing, so a
	// bare "\n" no longer returns the cursor to column 0. Build the frame and
	// emit "\r\n" explicitly, erasing the rest of each line first so a shorter
	// frame leaves no residue from the previous one.
	var frame strings.Builder
	frame.WriteString("\033[H")
	fmt.Fprintf(&frame, "%s\n\n", header)

	lines := trimLines(content)
	available := termHeight - watchHeaderLines

	if len(lines) <= available {
		// Content fits — no pagination needed
		frame.WriteString(strings.Join(lines, "\n"))
		frame.WriteString("\033[J")
		fmt.Print(strings.ReplaceAll(frame.String(), "\n", "\033[K\r\n"))
		return 0
	}

	// Reserve a line for the pagination footer
	pages := paginateWatch(lines, available-watchFooterLines)

	totalPages := len(pages)
	if page >= totalPages {
		page = totalPages - 1
	}
	if page < 0 {
		page = 0
	}

	frame.WriteString(strings.Join(pages[page], "\n"))
	// No trailing newline: the footer sits on the last row and a newline
	// there would scroll the screen, shifting every following frame.
	fmt.Fprintf(&frame, "\n\n  Page %d/%d    [ prev   ] next", page+1, totalPages)
	frame.WriteString("\033[J")
	fmt.Print(strings.ReplaceAll(frame.String(), "\n", "\033[K\r\n"))

	return page
}

// paginateWatch splits lines into pages that each fit in available lines.
// Table output keeps its top border, column headers, separator and bottom
// border on every page; anything else is split line by line. Always returns
// at least one page.
func paginateWatch(lines []string, available int) [][]string {
	var pages [][]string
	if isTableOutput(lines) {
		head := lines[:3]
		data := lines[3 : len(lines)-1]
		bottom := lines[len(lines)-1]
		pageSize := max(available-tableChrome, 1)
		for i := 0; i < len(data); i += pageSize {
			end := min(i+pageSize, len(data))
			p := make([]string, 0, 3+end-i+1)
			p = append(p, head...)
			p = append(p, data[i:end]...)
			p = append(p, bottom)
			pages = append(pages, p)
		}
	} else {
		pageSize := max(available, 1)
		for i := 0; i < len(lines); i += pageSize {
			end := min(i+pageSize, len(lines))
			pages = append(pages, lines[i:end])
		}
	}
	if len(pages) == 0 {
		pages = [][]string{lines}
	}
	return pages
}

// isTableOutput returns true when lines look like a lipgloss table
// (╭ top, ├ separator on line 2, ╰ bottom).
func isTableOutput(lines []string) bool {
	if len(lines) < 4 {
		return false
	}
	// Borders carry ANSI color codes on a real terminal; strip them before matching.
	return strings.HasPrefix(ansi.Strip(lines[0]), "╭") &&
		strings.HasPrefix(ansi.Strip(lines[2]), "├") &&
		strings.HasPrefix(ansi.Strip(lines[len(lines)-1]), "╰")
}

func trimLines(s string) []string {
	return strings.Split(strings.TrimRight(s, "\n"), "\n")
}
