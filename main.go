package main

import (
	"netmon/cmd"
)

// Version 정보
const (
	Version   = "1.2.2"
	BuildDate = "2025-12-05"
)

func main() {
	// Create configuration with version information
	cfg := cmd.Config{
		Version:   Version,
		BuildDate: BuildDate,
	}

	// Execute Cobra commands with configuration
	cmd.Execute(cfg)
}
