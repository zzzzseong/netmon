package main

import (
	"netmon/cmd"
)

// Version and BuildDate are injected by the release workflow via
// -ldflags "-X main.Version=<semver> -X main.BuildDate=<date>".
// Source builds without ldflags report "dev".
var (
	Version   = "dev"
	BuildDate = "unknown"
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
