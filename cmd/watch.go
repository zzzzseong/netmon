package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"time"

	"golang.org/x/term"
)

const (
	watchHeaderLines = 2 // command info line + blank line
	watchFooterLines = 1 // pagination indicator
	tableChrome      = 4 // top border + column headers + separator + bottom border
)

// runWithWatch runs fn on every interval, paginating output that exceeds the
// terminal height. [ / ] and left/right arrow keys navigate pages.
func runWithWatch(name string, interval time.Duration, fn func() (string, error)) error {
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

	// Initial fetch
	if out, err := fn(); err == nil {
		content = out
	}
	fmt.Print("\033[H\033[2J")
	page = renderWatch(name, interval, content, page)

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
			page = renderWatch(name, interval, content, page)
		case <-ticker.C:
			if out, err := fn(); err == nil {
				content = out
			}
			page = renderWatch(name, interval, content, page)
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
func renderWatch(name string, interval time.Duration, content string, page int) int {
	_, termHeight, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil || termHeight <= 0 {
		termHeight = 24
	}

	fmt.Print("\033[H")
	fmt.Printf("netmon %s  —  every %.0fs  —  %s  —  [ ] navigate  Ctrl+C to stop\n\n",
		name, interval.Seconds(), time.Now().Format("2006-01-02 15:04:05"))

	lines := trimLines(content)
	available := termHeight - watchHeaderLines

	if len(lines) <= available {
		// Content fits — no pagination needed
		fmt.Print(strings.Join(lines, "\n"))
		fmt.Print("\033[J")
		return 0
	}

	// Reserve a line for the pagination footer
	available -= watchFooterLines

	var pages [][]string
	if isTableOutput(lines) {
		// Preserve top border + header + separator on every page; bottom border at end
		head := lines[:3]
		data := lines[3 : len(lines)-1]
		bottom := lines[len(lines)-1]
		pageSize := available - tableChrome
		if pageSize < 1 {
			pageSize = 1
		}
		for i := 0; i < len(data); i += pageSize {
			end := i + pageSize
			if end > len(data) {
				end = len(data)
			}
			p := make([]string, 0, 3+end-i+1)
			p = append(p, head...)
			p = append(p, data[i:end]...)
			p = append(p, bottom)
			pages = append(pages, p)
		}
	} else {
		// Non-table (e.g. stats box): simple line-based pagination
		pageSize := available
		if pageSize < 1 {
			pageSize = 1
		}
		for i := 0; i < len(lines); i += pageSize {
			end := i + pageSize
			if end > len(lines) {
				end = len(lines)
			}
			pages = append(pages, lines[i:end])
		}
	}

	totalPages := len(pages)
	if totalPages == 0 {
		totalPages = 1
	}
	if page >= totalPages {
		page = totalPages - 1
	}
	if page < 0 {
		page = 0
	}

	fmt.Print(strings.Join(pages[page], "\n"))
	fmt.Printf("\n\n  Page %d/%d    [ prev   ] next\n", page+1, totalPages)
	fmt.Print("\033[J")

	return page
}

// isTableOutput returns true when lines look like a lipgloss table
// (╭ top, ├ separator on line 2, ╰ bottom).
func isTableOutput(lines []string) bool {
	if len(lines) < 4 {
		return false
	}
	return strings.HasPrefix(lines[0], "╭") &&
		strings.HasPrefix(lines[2], "├") &&
		strings.HasPrefix(lines[len(lines)-1], "╰")
}

func trimLines(s string) []string {
	lines := strings.Split(strings.TrimRight(s, "\n"), "\n")
	return lines
}
