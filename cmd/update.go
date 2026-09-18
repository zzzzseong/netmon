package cmd

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

const githubRepo = "zzzzseong/netmon"

// httpClient bounds every network call so a stalled GitHub response cannot hang the command forever.
var httpClient = &http.Client{Timeout: 2 * time.Minute}

type releaseInfo struct {
	TagName string `json:"tag_name"`
}

func newUpdateCmd(cfg Config) *cobra.Command {
	return &cobra.Command{
		Use:   "update",
		Short: "Update netmon to the latest version",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runUpdate(cfg)
		},
	}
}

func runUpdate(cfg Config) error {
	fmt.Println("Checking for updates...")

	latest, err := fetchLatestTag()
	if err != nil {
		return fmt.Errorf("failed to fetch latest version: %w", err)
	}

	current := "dev build"
	if cfg.Version != "dev" {
		current = "v" + cfg.Version
		if !isNewerVersion(latest, current) {
			fmt.Printf("Already up to date (%s)\n", current)
			return nil
		}
	}

	selfPath, err := resolveExecutable()
	if err != nil {
		return err
	}
	if isHomebrewPath(selfPath) {
		return upgradeWithHomebrew(current, latest)
	}

	fmt.Printf("Updating %s → %s\n", current, latest)

	platform, err := currentPlatform()
	if err != nil {
		return err
	}

	tmpDir, err := os.MkdirTemp("", "netmon-update-*")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmpDir)

	archive := archiveName(platform)
	archivePath := filepath.Join(tmpDir, archive)
	downloadURL := fmt.Sprintf("https://github.com/%s/releases/download/%s/%s", githubRepo, latest, archive)

	fmt.Printf("Downloading %s...\n", archive)
	if err := downloadFile(archivePath, downloadURL); err != nil {
		return fmt.Errorf("download failed: %w", err)
	}

	fmt.Println("Verifying checksum...")
	checksumURL := fmt.Sprintf("https://github.com/%s/releases/download/%s/SHA256SUMS", githubRepo, latest)
	if err := verifyChecksum(archivePath, archive, checksumURL); err != nil {
		return err
	}

	if err := extractBinary(archivePath, tmpDir); err != nil {
		return fmt.Errorf("extraction failed: %w", err)
	}

	newBinary := filepath.Join(tmpDir, execName())
	fmt.Printf("Installing to %s...\n", selfPath)
	if err := installBinary(newBinary, selfPath); err != nil {
		return fmt.Errorf("installation failed: %w", err)
	}

	fmt.Printf("Updated to %s\n", latest)
	return nil
}

func fetchLatestTag() (string, error) {
	resp, err := httpClient.Get(fmt.Sprintf("https://api.github.com/repos/%s/releases/latest", githubRepo))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	return parseLatestTag(resp)
}

// parseLatestTag extracts tag_name from a releases/latest response, surfacing
// the API's own message on a non-200 status (e.g. rate limiting) instead of a
// misleading "empty tag_name" error.
func parseLatestTag(resp *http.Response) (string, error) {
	if resp.StatusCode != http.StatusOK {
		var apiErr struct {
			Message string `json:"message"`
		}
		_ = json.NewDecoder(resp.Body).Decode(&apiErr)
		if apiErr.Message != "" {
			return "", fmt.Errorf("GitHub API returned %s: %s", resp.Status, apiErr.Message)
		}
		return "", fmt.Errorf("GitHub API returned %s", resp.Status)
	}

	var rel releaseInfo
	if err := json.NewDecoder(resp.Body).Decode(&rel); err != nil {
		return "", err
	}
	if rel.TagName == "" {
		return "", fmt.Errorf("empty tag_name in GitHub API response")
	}
	return rel.TagName, nil
}

func currentPlatform() (string, error) {
	goos := runtime.GOOS
	goarch := runtime.GOARCH
	switch goos {
	case "linux", "darwin", "windows":
	default:
		return "", fmt.Errorf("unsupported OS: %s", goos)
	}
	switch goarch {
	case "amd64", "arm64":
	default:
		return "", fmt.Errorf("unsupported architecture: %s", goarch)
	}
	return fmt.Sprintf("%s-%s", goos, goarch), nil
}

func execName() string {
	if runtime.GOOS == "windows" {
		return "netmon.exe"
	}
	return "netmon"
}

// archiveName returns the release asset name for the platform.
// The release workflow publishes .zip for Windows and .tar.gz for everything else.
func archiveName(platform string) string {
	if runtime.GOOS == "windows" {
		return fmt.Sprintf("netmon-%s.zip", platform)
	}
	return fmt.Sprintf("netmon-%s.tar.gz", platform)
}

func downloadFile(dest, url string) error {
	resp, err := httpClient.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP %d from %s", resp.StatusCode, url)
	}

	f, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.Copy(f, resp.Body)
	return err
}

func verifyChecksum(filePath, filename, checksumURL string) error {
	resp, err := httpClient.Get(checksumURL)
	if err != nil {
		return fmt.Errorf("failed to download SHA256SUMS: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var expected string
	for _, line := range strings.Split(string(body), "\n") {
		fields := strings.Fields(line)
		if len(fields) == 2 && fields[1] == filename {
			expected = fields[0]
			break
		}
	}
	if expected == "" {
		return fmt.Errorf("no checksum entry found for %s", filename)
	}

	f, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return err
	}
	actual := fmt.Sprintf("%x", h.Sum(nil))

	if expected != actual {
		return fmt.Errorf("checksum mismatch:\n  expected: %s\n  actual:   %s", expected, actual)
	}
	return nil
}

func extractBinary(archivePath, destDir string) error {
	if strings.HasSuffix(archivePath, ".zip") {
		return extractZip(archivePath, destDir)
	}
	return extractTarGz(archivePath, destDir)
}

func extractZip(zipPath, destDir string) error {
	zr, err := zip.OpenReader(zipPath)
	if err != nil {
		return err
	}
	defer zr.Close()

	target := execName()
	for _, zf := range zr.File {
		if zf.FileInfo().IsDir() || filepath.Base(zf.Name) != target {
			continue
		}
		in, err := zf.Open()
		if err != nil {
			return err
		}
		out, err := os.OpenFile(filepath.Join(destDir, target), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
		if err != nil {
			in.Close()
			return err
		}
		_, copyErr := io.Copy(out, in)
		in.Close()
		out.Close()
		return copyErr
	}
	return fmt.Errorf("binary %s not found in archive", target)
}

func extractTarGz(tarballPath, destDir string) error {
	f, err := os.Open(tarballPath)
	if err != nil {
		return err
	}
	defer f.Close()

	gz, err := gzip.NewReader(f)
	if err != nil {
		return err
	}
	defer gz.Close()

	target := execName()
	tr := tar.NewReader(gz)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if hdr.Typeflag != tar.TypeReg || filepath.Base(hdr.Name) != target {
			continue
		}
		out, err := os.OpenFile(filepath.Join(destDir, target), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
		if err != nil {
			return err
		}
		_, copyErr := io.Copy(out, tr)
		out.Close()
		return copyErr
	}
	return fmt.Errorf("binary %s not found in archive", target)
}

// isNewerVersion reports whether candidate is a strictly newer semantic version
// than current. Both are "vMAJOR.MINOR.PATCH" tags. Unparseable input is treated
// as newer so an unexpected tag format never blocks an update.
func isNewerVersion(candidate, current string) bool {
	parse := func(v string) ([3]int, bool) {
		var out [3]int
		parts := strings.Split(strings.TrimPrefix(v, "v"), ".")
		if len(parts) != 3 {
			return out, false
		}
		for i, p := range parts {
			n, err := strconv.Atoi(p)
			if err != nil {
				return out, false
			}
			out[i] = n
		}
		return out, true
	}
	c, okC := parse(candidate)
	cur, okCur := parse(current)
	if !okC || !okCur {
		return candidate != current
	}
	for i := 0; i < 3; i++ {
		if c[i] != cur[i] {
			return c[i] > cur[i]
		}
	}
	return false
}

// homebrewFormula is the fully qualified tap formula, so brew never resolves a
// different "netmon" from another tap.
const homebrewFormula = "zzzzseong/netmon/netmon"

// upgradeWithHomebrew delegates the upgrade to brew so the Cellar and brew's
// own records stay in sync. brew update runs first because brew upgrade only
// auto-updates taps once a day and would otherwise miss a fresh release.
func upgradeWithHomebrew(current, latest string) error {
	brew, err := exec.LookPath("brew")
	if err != nil {
		return fmt.Errorf("netmon was installed with Homebrew but `brew` is not in PATH; run `brew upgrade netmon` manually")
	}

	fmt.Printf("Installed with Homebrew; updating %s → %s via brew\n", current, latest)
	for _, args := range [][]string{{"update", "--quiet"}, {"upgrade", homebrewFormula}} {
		c := exec.Command(brew, args...)
		c.Stdout, c.Stderr = os.Stdout, os.Stderr
		if err := c.Run(); err != nil {
			return fmt.Errorf("brew %s failed: %w", strings.Join(args, " "), err)
		}
	}
	return nil
}

// isHomebrewPath reports whether the binary lives inside a Homebrew (or Linuxbrew) prefix.
func isHomebrewPath(path string) bool {
	p := filepath.ToSlash(path)
	return strings.Contains(p, "/Cellar/") ||
		strings.Contains(p, "/homebrew/") ||
		strings.Contains(p, "/linuxbrew/")
}

func resolveExecutable() (string, error) {
	self, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("cannot determine current executable: %w", err)
	}
	return filepath.EvalSymlinks(self)
}

func installBinary(src, dest string) error {
	// Write to a temp file next to dest, then rename atomically.
	tmp := dest + ".new"
	err := copyExecutable(src, tmp)
	if err == nil {
		if runtime.GOOS == "windows" {
			// A running .exe cannot be overwritten on Windows, but it can be renamed.
			// Move the current binary aside first; the leftover .old file is removed on the next update.
			old := dest + ".old"
			_ = os.Remove(old)
			if err := os.Rename(dest, old); err != nil {
				return err
			}
		}
		return os.Rename(tmp, dest)
	}
	// Permission denied — retry with sudo (not available on Windows).
	if os.IsPermission(err) && runtime.GOOS != "windows" {
		fmt.Println("Elevated privileges required, retrying with sudo...")
		return sudoInstall(src, dest)
	}
	return err
}

func copyExecutable(src, dest string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.OpenFile(dest, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	return err
}

func sudoInstall(src, dest string) error {
	cmd := exec.Command("sudo", "install", "-m", "0755", src, dest)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
