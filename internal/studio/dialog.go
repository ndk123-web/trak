package studio

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)

// OpenNativeFolderDialog invokes the native operating system folder browser dialog and returns the selected path.
func OpenNativeFolderDialog() (string, error) {
	if runtime.GOOS == "windows" {
		psScript := `
Add-Type -AssemblyName System.Windows.Forms
$f = New-Object System.Windows.Forms.FolderBrowserDialog
$f.Description = "Select Trak Workspace Directory"
$f.ShowNewFolderButton = $true
if ($f.ShowDialog() -eq [System.Windows.Forms.DialogResult]::OK) {
    Write-Output $f.SelectedPath
}
`
		cmd := exec.Command("powershell", "-NoProfile", "-NonInteractive", "-Command", psScript)
		out, err := cmd.Output()
		if err != nil {
			return "", err
		}
		return strings.TrimSpace(string(out)), nil
	} else if runtime.GOOS == "darwin" {
		cmd := exec.Command("osascript", "-e", `POSIX path of (choose folder with prompt "Select Trak Workspace Directory")`)
		out, err := cmd.Output()
		if err != nil {
			return "", err
		}
		return strings.TrimSpace(string(out)), nil
	} else {
		cmd := exec.Command("zenity", "--file-selection", "--directory", "--title=Select Trak Workspace Directory")
		out, err := cmd.Output()
		if err == nil {
			return strings.TrimSpace(string(out)), nil
		}
		return "", fmt.Errorf("native folder dialog not available")
	}
}

// OpenBrowser opens the given URL in the user's default system browser.
func OpenBrowser(url string) {
	switch runtime.GOOS {
	case "windows":
		exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		exec.Command("open", url).Start()
	default:
		exec.Command("xdg-open", url).Start()
	}
}
