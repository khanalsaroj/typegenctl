package usecase

import (
	"os"
	"path/filepath"
	"testing"
)

func TestArchiveAssetName(t *testing.T) {
	cases := []struct {
		goos, goarch, want string
	}{
		{"linux", "amd64", "typegenctl-linux-amd64.tar.gz"},
		{"darwin", "arm64", "typegenctl-darwin-arm64.tar.gz"},
		{"freebsd", "amd64", "typegenctl-freebsd-amd64.tar.gz"},
		{"windows", "amd64", "typegenctl-windows-amd64.zip"},
		{"windows", "386", "typegenctl-windows-386.zip"},
	}
	for _, tc := range cases {
		if got := archiveAssetName(tc.goos, tc.goarch); got != tc.want {
			t.Errorf("archiveAssetName(%q, %q) = %q, want %q", tc.goos, tc.goarch, got, tc.want)
		}
	}
}

func TestBinaryFileName(t *testing.T) {
	if got := binaryFileName("linux"); got != "typegenctl" {
		t.Errorf("binaryFileName(linux) = %q", got)
	}
	if got := binaryFileName("windows"); got != "typegenctl.exe" {
		t.Errorf("binaryFileName(windows) = %q", got)
	}
}

func TestParseChecksums(t *testing.T) {
	data := []byte(
		"abc123  typegenctl-linux-amd64.tar.gz\n" +
			"DEF456 *typegenctl-windows-amd64.zip\n" +
			"\n" +
			"shorty\n",
	)
	got := parseChecksums(data)

	if got["typegenctl-linux-amd64.tar.gz"] != "abc123" {
		t.Errorf("linux sum = %q, want abc123", got["typegenctl-linux-amd64.tar.gz"])
	}
	// '*' (binary-mode) prefix stripped and hex lowercased.
	if got["typegenctl-windows-amd64.zip"] != "def456" {
		t.Errorf("windows sum = %q, want def456", got["typegenctl-windows-amd64.zip"])
	}
	if _, ok := got["shorty"]; ok {
		t.Error("single-field line should be skipped")
	}
}

func TestSha256File(t *testing.T) {
	p := filepath.Join(t.TempDir(), "data.bin")
	if err := os.WriteFile(p, []byte("hello"), 0o644); err != nil {
		t.Fatal(err)
	}
	const want = "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824"
	got, err := sha256File(p)
	if err != nil {
		t.Fatal(err)
	}
	if got != want {
		t.Fatalf("sha256File = %q, want %q", got, want)
	}
}

func TestFindAsset(t *testing.T) {
	assets := []githubAsset{
		{Name: "checksums.txt", URL: "u1"},
		{Name: "typegenctl-linux-amd64.tar.gz", URL: "u2"},
	}
	if a, ok := findAsset(assets, "checksums.txt"); !ok || a.URL != "u1" {
		t.Errorf("findAsset(checksums.txt) = %+v, ok = %v", a, ok)
	}
	if _, ok := findAsset(assets, "missing"); ok {
		t.Error("expected missing asset to report not found")
	}
}
