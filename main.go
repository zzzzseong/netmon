package main

import (
	"netmon/cmd"
)

// Version 정보
const Version = "1.6.4"

// BuildDate is set at build time via -ldflags "-X main.BuildDate=<date>"
var BuildDate = "dev"

func main() {
	// Create configuration with version information
	cfg := cmd.Config{
		Version:   Version,
		BuildDate: BuildDate,
	}

	// Execute Cobra commands with configuration
	cmd.Execute(cfg)
}
