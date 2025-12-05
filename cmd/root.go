package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// newRootCmd creates and returns the root command for the netmon CLI.
// It sets up the base command with help and usage functions.
func newRootCmd(cfg Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "netmon",
		Short: "A modern, beautiful CLI tool for network monitoring and process management",
		Long: `A beautifully designed network monitoring tool built with Go that provides an intuitive interface for viewing active ports, network interfaces, routing tables, and managing processes on Linux, macOS, and Windows.`,
		Run: func(cmd *cobra.Command, args []string) {
			// No arguments provided, show help
			printCustomHelp(cmd, args, cfg)
		},
	}

	// Hide completion command from help (but keep it functional for Homebrew/install scripts)
	// This allows Homebrew to generate completions while hiding it from users
	cmd.CompletionOptions.HiddenDefaultCmd = true

	// Add all subcommands
	cmd.AddCommand(newRouteCmd())
	cmd.AddCommand(newIPCmd())
	cmd.AddCommand(newFindCmd())
	cmd.AddCommand(newLsCmd())
	cmd.AddCommand(newTracerouteCmd())
	cmd.AddCommand(newShutdownCmd())
	cmd.AddCommand(newVersionCmd(cfg))

	return cmd
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(cfg Config) {
	rootCmd := newRootCmd(cfg)

	// Set custom help and usage functions
	rootCmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		printCustomHelp(cmd, args, cfg)
	})
	rootCmd.SetUsageFunc(printCustomUsage)

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
