package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"
)

func runWithWatch(name string, interval time.Duration, fn func() error) error {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	// Enter alternate screen buffer: isolates watch output from shell history
	// and prevents scroll corruption when content exceeds terminal height.
	// Restore main screen + cursor on exit.
	fmt.Print("\033[?1049h\033[?25l")
	defer fmt.Print("\033[?1049l\033[?25h")

	// Clear alternate screen once on entry
	fmt.Print("\033[H\033[2J")

	for {
		// Move to top-left and overwrite in place — no blank-screen flash
		fmt.Print("\033[H")
		fmt.Printf("netmon %s  —  every %.0fs  —  %s  —  Ctrl+C to stop\n\n",
			name, interval.Seconds(), time.Now().Format("2006-01-02 15:04:05"))

		if err := fn(); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
		}

		// Clear from cursor to end (handles content shrinkage between refreshes)
		fmt.Print("\033[J")

		select {
		case <-ctx.Done():
			return nil
		case <-time.After(interval):
		}
	}
}
