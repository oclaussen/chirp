// +build darwin

package clipboard

import (
	"os/exec"
)

func copyCmd() (*exec.Cmd, error) {
	if _, err := exec.LookPath("pbcopy"); err == nil {
		return exec.Command("pbcopy"), nil
	}

	return nil, ErrNoClipboardFound
}

func pasteCmd() (*exec.Cmd, error) {
	if _, err := exec.LookPath("pbpaste"); err == nil {
		return exec.Command("pbpaste"), nil
	}

	return nil, ErrNoClipboardFound
}
