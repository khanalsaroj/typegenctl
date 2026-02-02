package usecase

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/sarojkhanal/typegenctl/internal/app"
	"github.com/sarojkhanal/typegenctl/internal/version"
)

const (
	installerURL = "https://raw.githubusercontent.com/khanalsaroj/typegenctl/refs/heads/main/main/install.sh"
	releaseAPI   = "https://api.github.com/repos/khanalsaroj/typegenctl/releases/latest"
)

type githubRelease struct {
	TagName string `json:"tag_name"`
}

func SelfUpdate(options *app.Options) error {
	fmt.Println("Checking for updates...")

	current := strings.TrimPrefix(version.Version, "v")

	latest, err := fetchLatestVersion()
	if err != nil {
		return err
	}

	fmt.Printf("Current version: %s\n", current)
	fmt.Printf("Latest version:  %s\n", latest)

	if current == latest && !options.Force {
		fmt.Println("Already up to date.")
		return nil
	}

	if options.CheckUpdate {
		fmt.Println("Update available.")
		return nil
	}

	fmt.Println("Downloading and installing update...")
	return runInstaller()
}

func fetchLatestVersion() (string, error) {
	req, err := http.NewRequest("GET", releaseAPI, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("User-Agent", "typegenctl-self-update")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer func(Body io.ReadCloser) {
		if Body.Close() != nil {
		}
	}(resp.Body)

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("failed to fetch latest version: %s", resp.Status)
	}

	var rel githubRelease
	if err := json.NewDecoder(resp.Body).Decode(&rel); err != nil {
		return "", err
	}

	return strings.TrimPrefix(rel.TagName, "v"), nil
}

func runInstaller() error {
	cmd := exec.Command("bash", "-c",
		fmt.Sprintf("curl -fsSL %s | sudo bash", installerURL),
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	return cmd.Run()
}
