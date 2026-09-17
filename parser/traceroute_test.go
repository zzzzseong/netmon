package parser

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

func captureStdout(t *testing.T, fn func()) string {
	t.Helper()

	oldStdout := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("failed to create pipe: %v", err)
	}
	os.Stdout = w

	fn()

	_ = w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, r); err != nil {
		t.Fatalf("failed to read captured stdout: %v", err)
	}
	_ = r.Close()

	return buf.String()
}

func TestParseWindowsTracert_LocaleAgnosticTimeoutAndSuccess(t *testing.T) {
	input := strings.Join([]string{
		"Tracing route to 8.8.8.8 over a maximum of 30 hops:",
		"",
		"  1     *        *        *     Request timed out.",
		"  2     *        *        *     요청 시간이 만료되었습니다.",
		"  3    10 ms    11 ms    12 ms  1.1.1.1",
	}, "\n")

	p := NewTracerouteParser()
	out := captureStdout(t, func() {
		scanner := bufio.NewScanner(strings.NewReader(input))
		p.ParseWindowsTracert(scanner)
	})

	if !strings.Contains(out, "Request timed out") {
		t.Fatalf("expected timeout row in output, got: %q", out)
	}
	if !strings.Contains(out, "1.1.1.1") {
		t.Fatalf("expected success row with IP, got: %q", out)
	}
}

func TestParseUnixTraceroute_KeepsFirstHopWhenHeaderOnStderr(t *testing.T) {
	// BSD/macOS traceroute writes the "traceroute to ..." header to stderr,
	// so the first stdout line is already hop 1 and must not be skipped.
	input := strings.Join([]string{
		" 1  10.0.0.1  1.234 ms  1.345 ms  1.456 ms",
		" 2  * * *",
	}, "\n")

	p := NewTracerouteParser()
	out := captureStdout(t, func() {
		p.ParseUnixTraceroute(bufio.NewScanner(strings.NewReader(input)))
	})

	if !strings.Contains(out, "10.0.0.1") {
		t.Fatalf("expected hop 1 in output, got: %q", out)
	}
	if !strings.Contains(out, "Request timed out") {
		t.Fatalf("expected timeout row for hop 2, got: %q", out)
	}
}

func TestParseUnixTraceroute_IgnoresHeaderOnStdout(t *testing.T) {
	// Linux traceroute writes the header to stdout; it must be ignored, not treated as a hop.
	input := strings.Join([]string{
		"traceroute to 1.1.1.1 (1.1.1.1), 30 hops max, 60 byte packets",
		" 1  10.0.0.1  1.234 ms  1.345 ms  1.456 ms",
	}, "\n")

	p := NewTracerouteParser()
	out := captureStdout(t, func() {
		p.ParseUnixTraceroute(bufio.NewScanner(strings.NewReader(input)))
	})

	if strings.Count(out, "\n") != 1 {
		t.Fatalf("expected exactly one hop line, got: %q", out)
	}
	if !strings.Contains(out, "10.0.0.1") {
		t.Fatalf("expected hop 1 in output, got: %q", out)
	}
}
