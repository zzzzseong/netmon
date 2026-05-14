//go:build !windows

package cmd

import (
	"os"
	"syscall"
)

// tracerouteSignals are the OS signals that cancel a running traceroute.
// On Unix, SIGTERM is included so that `kill <pid>` also terminates the child process.
var tracerouteSignals = []os.Signal{os.Interrupt, syscall.SIGTERM}
