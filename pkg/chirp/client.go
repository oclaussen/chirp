package chirp

import (
	"io"
	"net"
	"os"
	"path/filepath"
)

func Copy(socketType string, addr string) error {
	if socketType == "unix" {
		addr, err := filepath.Abs(addr)
		if err != nil {
			return err
		}
		if _, err := os.Stat(addr); err != nil {
			return err
		}
	}

	conn, err := net.Dial(socketType, addr)
	if err != nil {
		return err
	}
	defer conn.Close()

	if _, err := conn.Write([]byte(copyCmd)); err != nil {
		return err
	}

	if _, err := io.Copy(conn, os.Stdin); err != nil {
		return err
	}

	return nil
}

func Paste(socketType string, addr string) error {
	if socketType == "unix" {
		addr, err := filepath.Abs(addr)
		if err != nil {
			return err
		}
		if _, err := os.Stat(addr); err != nil {
			return err
		}
	}

	conn, err := net.Dial(socketType, addr)
	if err != nil {
		return err
	}
	defer conn.Close()

	if _, err := conn.Write([]byte(pasteCmd)); err != nil {
		return err
	}

	if _, err = io.Copy(os.Stdout, conn); err != nil {
		return err
	}

	return nil
}
