package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "netmon",
	Short: "A modern, beautiful CLI tool for network monitoring and process management",
	Long: `A beautifully designed network monitoring tool built with Go that provides an intuitive interface for viewing active ports, network interfaces, routing tables, and managing processes on Linux, macOS, and Windows.`,
	Run: func(cmd *cobra.Command, args []string) {
		// No arguments provided, show help
		printCustomHelp(cmd, args)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	// Set custom help and usage functions
	rootCmd.SetHelpFunc(printCustomHelp)
	rootCmd.SetUsageFunc(printCustomUsage)
	
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.netmon.yaml)")
}
