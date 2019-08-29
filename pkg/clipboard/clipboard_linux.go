// +build !darwin,!windows

package clipboard

import (
	"os"
	"os/exec"

	"github.com/pkg/errors"
)

func copyCmd() (*exec.Cmd, error) {
	if len(os.Getenv("WAYLAND_DISPLAY")) > 0 {
		if _, err := exec.LookPath("wl-copy"); err == nil {
			return exec.Command("wl-copy"), nil
		}
	}

	if _, err := exec.LookPath("xclip"); err == nil {
		return exec.Command("xclip", "-in", "-selection", "clipboard"), nil
	}

	if _, err := exec.LookPath("xsel"); err == nil {
		return exec.Command("xsel", "--input", "--clipboard"), nil
	}

	if _, err := exec.LookPath("termux-clipboard-set"); err == nil {
		return exec.Command("termux-clipboard-set"), nil
	}

	return nil, errors.New("no clipboard found")
}

func pasteCmd() (*exec.Cmd, error) {
	if len(os.Getenv("WAYLAND_DISPLAY")) > 0 {
		if _, err := exec.LookPath("wl-paste"); err == nil {
			return exec.Command("wl-paste", "--no-newline"), nil
		}
	}

	if _, err := exec.LookPath("xclip"); err == nil {
		return exec.Command("xclip", "-out", "-selection", "clipboard"), nil
	}

	if _, err := exec.LookPath("xsel"); err == nil {
		return exec.Command("xsel", "--output", "--clipboard"), nil
	}

	if _, err := exec.LookPath("termux-clipboard-get"); err == nil {
		return exec.Command("termux-clipboard-get"), nil
	}

	return nil, errors.New("no clipboard found")
}
