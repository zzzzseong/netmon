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
