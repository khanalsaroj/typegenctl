package browser

import (
	"os"
	"os/exec"
	"runtime"
)

func Open(url string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", url)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	default:
		cmd = exec.Command("xdg-open", url)
	}

	return cmd.Start()
}

func CanAutoOpen() bool {
	return os.Getenv("CI") == ""
}
