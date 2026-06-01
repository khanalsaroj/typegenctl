package usecase

import (
	"archive/tar"
	"archive/zip"
	"bufio"
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/khanalsaroj/typegenctl/internal/app"
	"github.com/khanalsaroj/typegenctl/internal/version"
)

const (
	binaryName  = "typegenctl"
	releaseAPI  = "https://api.github.com/repos/khanalsaroj/typegenctl/releases/latest"
	httpTimeout = 60 * time.Second
)

type githubAsset struct {
	Name string `json:"name"`
	URL  string `json:"browser_download_url"`
}

type githubRelease struct {
	TagName string        `json:"tag_name"`
	Assets  []githubAsset `json:"assets"`
}

func httpClient() *http.Client { return &http.Client{Timeout: httpTimeout} }

func SelfUpdate(options *app.Options) error {
	fmt.Println("Checking for updates...")

	current := strings.TrimPrefix(version.Version, "v")

	rel, err := fetchLatestRelease()
	if err != nil {
		return err
	}
	latest := strings.TrimPrefix(rel.TagName, "v")

	fmt.Printf("Current version: %s\n", current)
	fmt.Printf("Latest version:  %s\n", latest)

	upToDate := current == latest
	if upToDate && !options.Force {
		fmt.Println("Already up to date.")
		return nil
	}

	if options.CheckUpdate {
		if upToDate {
			fmt.Println("No update available (use --force to reinstall).")
		} else {
			fmt.Println("Update available.")
		}
		return nil
	}

	fmt.Println("Downloading and verifying update...")
	return installRelease(rel)
}

func fetchLatestRelease() (*githubRelease, error) {
	req, err := http.NewRequest(http.MethodGet, releaseAPI, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "typegenctl-self-update")
	req.Header.Set("Accept", "application/vnd.github+json")

	resp, err := httpClient().Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch latest release: %s", resp.Status)
	}

	var rel githubRelease
	if err := json.NewDecoder(resp.Body).Decode(&rel); err != nil {
		return nil, err
	}
	if rel.TagName == "" {
		return nil, errors.New("latest release has no tag")
	}
	return &rel, nil
}

// installRelease downloads the platform archive, verifies it against the
// release's checksums.txt, extracts the binary, and replaces the running
// executable. It refuses to proceed if the checksum cannot be verified.
func installRelease(rel *githubRelease) error {
	goos, goarch := runtime.GOOS, runtime.GOARCH

	archiveName := archiveAssetName(goos, goarch)
	asset, ok := findAsset(rel.Assets, archiveName)
	if !ok {
		return fmt.Errorf("no release asset %q for %s/%s", archiveName, goos, goarch)
	}
	checksumAsset, ok := findAsset(rel.Assets, "checksums.txt")
	if !ok {
		return errors.New("release is missing checksums.txt; refusing to update without verification")
	}

	tmpDir, err := os.MkdirTemp("", "typegenctl-update-")
	if err != nil {
		return err
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	archivePath := filepath.Join(tmpDir, archiveName)
	fmt.Printf("Downloading %s...\n", archiveName)
	if err := downloadToFile(asset.URL, archivePath); err != nil {
		return err
	}

	// Verify the checksum before touching anything executable.
	sumsData, err := downloadBytes(checksumAsset.URL)
	if err != nil {
		return fmt.Errorf("failed to download checksums: %w", err)
	}
	want, ok := parseChecksums(sumsData)[archiveName]
	if !ok {
		return fmt.Errorf("no checksum listed for %s", archiveName)
	}
	got, err := sha256File(archivePath)
	if err != nil {
		return err
	}
	if !strings.EqualFold(got, want) {
		return fmt.Errorf("checksum mismatch for %s:\n  expected %s\n  got      %s", archiveName, want, got)
	}
	fmt.Println("Checksum verified.")

	binInArchive := binaryFileName(goos)
	extractedBin := filepath.Join(tmpDir, binInArchive)
	if err := extractBinary(archivePath, binInArchive, extractedBin); err != nil {
		return err
	}

	if err := replaceExecutable(extractedBin); err != nil {
		return err
	}

	fmt.Println("Update complete.")
	return nil
}

// archiveAssetName returns the release archive filename for the given platform,
// matching the naming used by the release workflow.
func archiveAssetName(goos, goarch string) string {
	if goos == "windows" {
		return fmt.Sprintf("%s-%s-%s.zip", binaryName, goos, goarch)
	}
	return fmt.Sprintf("%s-%s-%s.tar.gz", binaryName, goos, goarch)
}

// binaryFileName returns the executable name as packaged inside the archive.
func binaryFileName(goos string) string {
	if goos == "windows" {
		return binaryName + ".exe"
	}
	return binaryName
}

func findAsset(assets []githubAsset, name string) (githubAsset, bool) {
	for _, a := range assets {
		if a.Name == name {
			return a, true
		}
	}
	return githubAsset{}, false
}

// parseChecksums parses `sha256sum`-style output ("<hex>  <filename>") into a
// map of base filename -> lowercase hex digest.
func parseChecksums(data []byte) map[string]string {
	out := make(map[string]string)
	sc := bufio.NewScanner(bytes.NewReader(data))
	for sc.Scan() {
		fields := strings.Fields(strings.TrimSpace(sc.Text()))
		if len(fields) < 2 {
			continue
		}
		// sha256sum prefixes binary-mode entries with '*'.
		name := strings.TrimPrefix(fields[len(fields)-1], "*")
		out[filepath.Base(name)] = strings.ToLower(fields[0])
	}
	return out
}

func sha256File(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer func() { _ = f.Close() }()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

func downloadToFile(url, dest string) error {
	resp, err := httpGet(url)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download failed (%s): %s", url, resp.Status)
	}

	f, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer func() { _ = f.Close() }()

	_, err = io.Copy(f, resp.Body)
	return err
}

func downloadBytes(url string) ([]byte, error) {
	resp, err := httpGet(url)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("download failed (%s): %s", url, resp.Status)
	}
	return io.ReadAll(resp.Body)
}

func httpGet(url string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "typegenctl-self-update")
	return httpClient().Do(req)
}

func extractBinary(archivePath, wantName, dest string) error {
	if strings.HasSuffix(archivePath, ".zip") {
		return extractFromZip(archivePath, wantName, dest)
	}
	return extractFromTarGz(archivePath, wantName, dest)
}

func extractFromTarGz(archivePath, wantName, dest string) error {
	f, err := os.Open(archivePath)
	if err != nil {
		return err
	}
	defer func() { _ = f.Close() }()

	gz, err := gzip.NewReader(f)
	if err != nil {
		return err
	}
	defer func() { _ = gz.Close() }()

	tr := tar.NewReader(gz)
	for {
		hdr, err := tr.Next()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return err
		}
		if filepath.Base(hdr.Name) == wantName {
			return writeBinary(dest, tr)
		}
	}
	return fmt.Errorf("binary %q not found in archive", wantName)
}

func extractFromZip(archivePath, wantName, dest string) error {
	r, err := zip.OpenReader(archivePath)
	if err != nil {
		return err
	}
	defer func() { _ = r.Close() }()

	for _, zf := range r.File {
		if filepath.Base(zf.Name) != wantName {
			continue
		}
		rc, err := zf.Open()
		if err != nil {
			return err
		}
		err = writeBinary(dest, rc)
		_ = rc.Close()
		return err
	}
	return fmt.Errorf("binary %q not found in archive", wantName)
}

func writeBinary(dest string, r io.Reader) error {
	out, err := os.OpenFile(dest, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o755)
	if err != nil {
		return err
	}
	if _, err := io.Copy(out, r); err != nil {
		_ = out.Close()
		return err
	}
	return out.Close()
}

// replaceExecutable swaps the currently running binary for newBin. It stages
// the replacement in the destination directory so the final move is atomic, and
// moves the running binary aside first so the swap also works on Windows (where
// a running executable can be renamed but not overwritten).
func replaceExecutable(newBin string) error {
	exe, err := os.Executable()
	if err != nil {
		return err
	}
	if resolved, err := filepath.EvalSymlinks(exe); err == nil {
		exe = resolved
	}

	dir := filepath.Dir(exe)
	staged := filepath.Join(dir, "."+filepath.Base(exe)+".new")
	if err := copyFile(newBin, staged, 0o755); err != nil {
		return fmt.Errorf("cannot stage update in %s: %w%s", dir, err, permHint(err))
	}

	old := exe + ".old"
	_ = os.Remove(old) // clean a leftover from a previous update, if any

	if err := os.Rename(exe, old); err != nil {
		_ = os.Remove(staged)
		return fmt.Errorf("cannot replace %s: %w%s", exe, err, permHint(err))
	}
	if err := os.Rename(staged, exe); err != nil {
		_ = os.Rename(old, exe) // roll back so a binary always exists
		_ = os.Remove(staged)
		return fmt.Errorf("cannot move new binary into place: %w", err)
	}

	// Best-effort: on Windows the old binary stays locked until the process exits.
	_ = os.Remove(old)
	return nil
}

func copyFile(src, dst string, mode os.FileMode) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer func() { _ = in.Close() }()

	out, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, mode)
	if err != nil {
		return err
	}
	if _, err := io.Copy(out, in); err != nil {
		_ = out.Close()
		return err
	}
	return out.Close()
}

func permHint(err error) string {
	if errors.Is(err, os.ErrPermission) {
		return "\nhint: the binary is in a protected directory; re-run with elevated privileges (e.g. sudo typegenctl self-update)"
	}
	return ""
}
