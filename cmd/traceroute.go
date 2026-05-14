package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"strconv"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
	"netmon/formatter"
	"netmon/parser"
	"netmon/style"
	"netmon/utils"
)

const (
	// tracerouteTimeoutSeconds is the timeout in seconds for Unix traceroute
	tracerouteTimeoutSeconds = 3
	// tracerouteProbeCount is the number of probes per hop for Unix traceroute
	tracerouteProbeCount = 3
	// tracertTimeoutMs is the timeout in milliseconds for Windows tracert
	tracertTimeoutMs = 3000
)

// newTracerouteCmd creates and returns the traceroute command.
// It traces the network path to a destination with animated loading.
func newTracerouteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "traceroute <host>",
		Short: "Trace route to network host",
		Long:  `Trace the network path to a destination with animated loading.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			target := args[0]

			// Validate input: hostname or IP address
			if err := utils.ValidateHostname(target); err != nil {
				return fmt.Errorf("invalid hostname or IP address: %w", err)
			}

			// Print header
			header := fmt.Sprintf("Tracing route to %s", target)
			headerStyle := lipgloss.NewStyle().
				Foreground(style.PrimaryColor).
				Bold(true)
			fmt.Println(headerStyle.Render(header))
			fmt.Println()

			// Create formatter and parser
			fmtter := formatter.NewTracerouteFormatter()
			parser := parser.NewTracerouteParser()

			// Print table header
			fmtter.PrintTableHeader()

			// Execute traceroute in real-time with cancellation support
			err := executeTracerouteStreaming(target, parser)
			if err != nil {
				return err
			}

			fmt.Println()
			return nil
		},
	}
}

// executeTracerouteStreaming executes traceroute in real-time and streams the output.
// It handles signal cancellation (Ctrl+C) to properly terminate the traceroute process.
func executeTracerouteStreaming(target string, parser *parser.TracerouteParser) error {
	var cmd *exec.Cmd
	var cmdName string

	if runtime.GOOS == "windows" {
		// Windows: tracert
		cmdName = "tracert"
		cmd = exec.Command("tracert", "-d", "-w", strconv.Itoa(tracertTimeoutMs), target)
	} else {
		// Unix-like: traceroute
		cmdName = "traceroute"
		timeoutStr := strconv.Itoa(tracerouteTimeoutSeconds)
		probeStr := strconv.Itoa(tracerouteProbeCount)
		cmd = exec.Command("traceroute", "-n", "-w", timeoutStr, "-q", probeStr, target)
	}

	// Check if command exists in PATH
	if _, err := exec.LookPath(cmdName); err != nil {
		// Print OS-specific installation instructions
		installMsg := getTracerouteInstallMessage()
		return fmt.Errorf("%s command not found in PATH\n%s", cmdName, installMsg)
	}

	// Create stdout pipe
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to create pipe: %w", err)
	}
	// Ensure pipe is closed for resource management
	defer stdout.Close()

	// Start command
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start traceroute: %w", err)
	}

	// Set up signal handling for graceful cancellation
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, tracerouteSignals...)
	defer signal.Stop(sigChan)

	// Goroutine to handle cancellation
	done := make(chan error, 1)
	go func() {
		// Read and parse output in real-time
		scanner := bufio.NewScanner(stdout)

		if runtime.GOOS == "windows" {
			parser.ParseWindowsTracert(scanner)
		} else {
			parser.ParseUnixTraceroute(scanner)
		}

		// Wait for process to finish
		done <- cmd.Wait()
	}()

	// Wait for either completion or cancellation
	select {
	case sig := <-sigChan:
		// Signal received (Ctrl+C)
		fmt.Fprintf(os.Stderr, "\nReceived signal: %v. Terminating traceroute...\n", sig)
		cmd.Process.Kill()
		<-done // Wait for cleanup
		return fmt.Errorf("traceroute interrupted")
	case err := <-done:
		// Process completed
		if err != nil {
			// Partial results may already be shown, but propagate failure for reliable exit codes.
			return fmt.Errorf("traceroute command exited with error: %w", err)
		}
		return nil
	}
}

// getTracerouteInstallMessage returns OS-specific installation instructions for traceroute.
func getTracerouteInstallMessage() string {
	var msg string

	switch runtime.GOOS {
	case "linux":
		msg = `Please install traceroute using one of the following commands:
  • Debian/Ubuntu: sudo apt-get install traceroute
  • RHEL/CentOS:   sudo yum install traceroute
  • Fedora:        sudo dnf install traceroute
  • Arch Linux:    sudo pacman -S traceroute`
	case "darwin":
		msg = `Please install traceroute using Homebrew:
  brew install traceroute
  
Note: macOS usually includes traceroute by default. If you see this message,
traceroute may have been removed or is not in your PATH.`
	case "windows":
		msg = `Windows includes tracert by default. If you see this message,
please check your system PATH configuration.`
	default:
		msg = fmt.Sprintf(`Please install traceroute for your operating system (%s).
Check your system's package manager for installation instructions.`, runtime.GOOS)
	}

	return style.DescStyle.Render(msg)
}
