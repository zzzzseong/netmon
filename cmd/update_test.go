package cmd

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"crypto/sha256"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestArchiveName_MatchesReleaseAssets(t *testing.T) {
	got := archiveName("linux-amd64")
	if runtime.GOOS == "windows" {
		if got != "netmon-linux-amd64.zip" {
			t.Fatalf("expected .zip on Windows, got %s", got)
		}
		return
	}
	if got != "netmon-linux-amd64.tar.gz" {
		t.Fatalf("expected .tar.gz, got %s", got)
	}
}

func TestIsNewerVersion(t *testing.T) {
	cases := []struct {
		candidate, current string
		want               bool
	}{
		{"v1.7.0", "v1.6.6", true},
		{"v1.6.6", "v1.7.0", false}, // never downgrade
		{"v1.7.0", "v1.7.0", false},
		{"v2.0.0", "v1.99.99", true},
		{"v1.7.10", "v1.7.9", true},
		{"nightly", "v1.7.0", true}, // unparseable tag still triggers an update
	}
	for _, c := range cases {
		if got := isNewerVersion(c.candidate, c.current); got != c.want {
			t.Errorf("isNewerVersion(%q, %q) = %v, want %v", c.candidate, c.current, got, c.want)
		}
	}
}

func TestIsHomebrewPath(t *testing.T) {
	cases := map[string]bool{
		"/opt/homebrew/Cellar/netmon/1.7.0/bin/netmon":  true,
		"/usr/local/Cellar/netmon/1.7.0/bin/netmon":     true,
		"/home/linuxbrew/.linuxbrew/bin/netmon":         true,
		"/usr/local/bin/netmon":                         false,
		`C:\Users\me\AppData\Local\Programs\netmon.exe`: false,
	}
	for path, want := range cases {
		if got := isHomebrewPath(path); got != want {
			t.Errorf("isHomebrewPath(%q) = %v, want %v", path, got, want)
		}
	}
}

func TestVerifyChecksum(t *testing.T) {
	dir := t.TempDir()
	artifact := filepath.Join(dir, "netmon-linux-amd64.tar.gz")
	content := []byte("release bytes")
	if err := os.WriteFile(artifact, content, 0644); err != nil {
		t.Fatal(err)
	}
	sum := fmt.Sprintf("%x", sha256.Sum256(content))

	serve := func(body string) *httptest.Server {
		return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, body)
		}))
	}

	t.Run("match", func(t *testing.T) {
		srv := serve("deadbeef  netmon-darwin-arm64.tar.gz\n" + sum + "  netmon-linux-amd64.tar.gz\n")
		defer srv.Close()
		if err := verifyChecksum(artifact, "netmon-linux-amd64.tar.gz", srv.URL); err != nil {
			t.Fatalf("expected checksum to match, got %v", err)
		}
	})

	t.Run("mismatch", func(t *testing.T) {
		srv := serve(strings.Repeat("0", 64) + "  netmon-linux-amd64.tar.gz\n")
		defer srv.Close()
		if err := verifyChecksum(artifact, "netmon-linux-amd64.tar.gz", srv.URL); err == nil {
			t.Fatal("expected checksum mismatch error")
		}
	})

	t.Run("missing entry", func(t *testing.T) {
		srv := serve("deadbeef  netmon-darwin-arm64.tar.gz\n")
		defer srv.Close()
		if err := verifyChecksum(artifact, "netmon-linux-amd64.tar.gz", srv.URL); err == nil {
			t.Fatal("expected error when SHA256SUMS has no entry for the artifact")
		}
	})
}

func TestExtractTarGz(t *testing.T) {
	dir := t.TempDir()
	archive := filepath.Join(dir, "netmon.tar.gz")
	f, err := os.Create(archive)
	if err != nil {
		t.Fatal(err)
	}
	gz := gzip.NewWriter(f)
	tw := tar.NewWriter(gz)
	for name, body := range map[string]string{"README": "ignore me", execName(): "binary"} {
		if err := tw.WriteHeader(&tar.Header{Name: name, Mode: 0755, Size: int64(len(body)), Typeflag: tar.TypeReg}); err != nil {
			t.Fatal(err)
		}
		if _, err := tw.Write([]byte(body)); err != nil {
			t.Fatal(err)
		}
	}
	tw.Close()
	gz.Close()
	f.Close()

	if err := extractBinary(archive, dir); err != nil {
		t.Fatalf("extract failed: %v", err)
	}
	got, err := os.ReadFile(filepath.Join(dir, execName()))
	if err != nil || string(got) != "binary" {
		t.Fatalf("expected extracted binary content, got %q (err %v)", got, err)
	}
}

func TestExtractZip(t *testing.T) {
	dir := t.TempDir()
	archive := filepath.Join(dir, "netmon.zip")
	f, err := os.Create(archive)
	if err != nil {
		t.Fatal(err)
	}
	zw := zip.NewWriter(f)
	for name, body := range map[string]string{"LICENSE": "ignore me", execName(): "binary"} {
		w, err := zw.Create(name)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := w.Write([]byte(body)); err != nil {
			t.Fatal(err)
		}
	}
	zw.Close()
	f.Close()

	if err := extractBinary(archive, dir); err != nil {
		t.Fatalf("extract failed: %v", err)
	}
	got, err := os.ReadFile(filepath.Join(dir, execName()))
	if err != nil || string(got) != "binary" {
		t.Fatalf("expected extracted binary content, got %q (err %v)", got, err)
	}
}

func TestExtractBinary_MissingBinary(t *testing.T) {
	dir := t.TempDir()
	archive := filepath.Join(dir, "empty.tar.gz")
	f, _ := os.Create(archive)
	gz := gzip.NewWriter(f)
	tar.NewWriter(gz).Close()
	gz.Close()
	f.Close()

	if err := extractBinary(archive, dir); err == nil {
		t.Fatal("expected error when archive has no binary")
	}
}
