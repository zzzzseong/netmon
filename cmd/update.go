package cmd

import (
	"archive/tar"
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
	"strings"

	"github.com/spf13/cobra"
)

const githubRepo = "zzzzseong/netmon"

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

	current := "v" + cfg.Version
	if latest == current {
		fmt.Printf("Already up to date (%s)\n", current)
		return nil
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

	tarball := fmt.Sprintf("netmon-%s.tar.gz", platform)
	tarballPath := filepath.Join(tmpDir, tarball)
	downloadURL := fmt.Sprintf("https://github.com/%s/releases/download/%s/%s", githubRepo, latest, tarball)

	fmt.Printf("Downloading %s...\n", tarball)
	if err := downloadFile(tarballPath, downloadURL); err != nil {
		return fmt.Errorf("download failed: %w", err)
	}

	fmt.Println("Verifying checksum...")
	checksumURL := fmt.Sprintf("https://github.com/%s/releases/download/%s/SHA256SUMS", githubRepo, latest)
	if err := verifyChecksum(tarballPath, tarball, checksumURL); err != nil {
		return err
	}

	if err := extractBinary(tarballPath, tmpDir); err != nil {
		return fmt.Errorf("extraction failed: %w", err)
	}

	selfPath, err := resolveExecutable()
	if err != nil {
		return err
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
	resp, err := http.Get(fmt.Sprintf("https://api.github.com/repos/%s/releases/latest", githubRepo))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

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

func downloadFile(dest, url string) error {
	resp, err := http.Get(url)
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
	resp, err := http.Get(checksumURL)
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

func extractBinary(tarballPath, destDir string) error {
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
		return os.Rename(tmp, dest)
	}
	// Permission denied — retry with sudo.
	if os.IsPermission(err) {
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
