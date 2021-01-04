package clipboard

import (
	"bytes"
	"fmt"
	"io"
	"strings"
)

const ErrNoClipboardFound Error = "no clipboard found"

type Error string

func (e Error) Error() string {
	return string(e)
}

func Check() error {
	if _, err := copyCmd(); err != nil {
		return err
	}

	if _, err := pasteCmd(); err != nil {
		return err
	}

	return nil
}

func Copy(source io.Reader) error {
	cmd, err := copyCmd()
	if err != nil {
		return err
	}

	in, err := cmd.StdinPipe()
	if err != nil {
		return fmt.Errorf("could not open copy stream: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("could not execute copy command: %w", err)
	}

	if _, err := io.Copy(in, source); err != nil {
		return fmt.Errorf("could not write copy data: %w", err)
	}

	in.Close()

	return cmd.Wait()
}

func CopyString(source string) error {
	return Copy(strings.NewReader(source))
}

func Paste(target io.Writer) error {
	cmd, err := pasteCmd()
	if err != nil {
		return err
	}

	out, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("could not open paste stream: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("could not execute paste command: %w", err)
	}

	if _, err := io.Copy(target, out); err != nil {
		return fmt.Errorf("could not read paste data: %w", err)
	}

	out.Close()

	return nil
}

func PasteString() (string, error) {
	var buf bytes.Buffer
	if err := Paste(&buf); err != nil {
		return "", err
	}

	return buf.String(), nil
}
