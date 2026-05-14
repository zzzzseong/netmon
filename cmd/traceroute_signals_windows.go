//go:build windows

package cmd

import "os"

// tracerouteSignals are the OS signals that cancel a running traceroute.
// On Windows, only os.Interrupt (Ctrl+C) is reliably deliverable.
var tracerouteSignals = []os.Signal{os.Interrupt}
