// +build darwin

package clipboard

import (
	"os/exec"

	"github.com/pkg/errors"
)

func copyCmd() (*exec.Cmd, error) {
	if _, err := exec.LookPath("pbcopy"); err == nil {
		return exec.Command("pbcopy"), nil
	}

	return nil, errors.New("no clipboard found")
}

func pasteCmd() (*exec.Cmd, error) {
	if _, err := exec.LookPath("pbpaste"); err == nil {
		return exec.Command("pbpaste"), nil
	}

	return nil, errors.New("no clipboard found")
}
