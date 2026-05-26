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

	for {
		// Clear alternate screen and move to top-left on every iteration.
		// Safe to do here (no blank-screen flash) because we're in the alt buffer.
		fmt.Print("\033[H\033[2J")
		fmt.Printf("netmon %s  —  every %.0fs  —  %s  —  Ctrl+C to stop\n\n",
			name, interval.Seconds(), time.Now().Format("2006-01-02 15:04:05"))

		if err := fn(); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
		}

		select {
		case <-ctx.Done():
			return nil
		case <-time.After(interval):
		}
	}
}
