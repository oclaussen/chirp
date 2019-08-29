package clipboard

import (
	"io"
)

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
		return err
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	if _, err := io.Copy(in, source); err != nil {
		return err
	}
	in.Close()

	return cmd.Wait()
}

func Paste(target io.Writer) error {
	cmd, err := pasteCmd()
	if err != nil {
		return err
	}

	out, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	if _, err := io.Copy(target, out); err != nil {
		return err
	}
	out.Close()

	return nil
}
