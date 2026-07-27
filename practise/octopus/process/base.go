package process

import (
	"os/exec"
	"runtime"
)

func OpenBrowser(url string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin": // macOS
		cmd = exec.Command("open", url)
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", url)
	default: // Linux 和其他 Unix
		cmd = exec.Command("xdg-open", url)
	}
	return cmd.Start() // Start 不会等待命令完成
}
