package chirp

import (
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/oclaussen/chirp/pkg/clipboard"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

var (
	copyCmd  = string('<')
	pasteCmd = string('>')
)

func Server(socketType string, addr string) error {
	if err := clipboard.Check(); err != nil {
		return err
	}

	if socketType == "unix" {
		addr, err := filepath.Abs(addr)
		if err != nil {
			return err
		}
		if _, err := os.Stat(addr); !os.IsNotExist(err) {
			if _, err = net.Dial("unix", addr); err == nil {
				return errors.New("socket already exists at: " + addr)
			}

			if err = os.Remove(addr); err != nil {
				return errors.Wrap(err, "could not remove stale socket: "+addr)
			}
		}
	}

	log.WithFields(log.Fields{"type": socketType, "address": addr}).Info("listening...")
	listener, err := net.Listen(socketType, addr)
	if err != nil {
		return err
	}
	defer listener.Close()

	go handleRequests(listener)

	killChan := make(chan os.Signal, 1)
	signal.Notify(killChan, os.Interrupt, syscall.SIGTERM)
	sig := <-killChan
	log.WithFields(log.Fields{"signal": sig}).Info("received signal")
	return nil

}

func handleRequests(listener net.Listener) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Error(err)
			return
		}
		log.Info("connection opened")

		go func() {
			defer func() {
				conn.Close()
				log.Info("connection closed")
			}()

			cmd := make([]byte, 1)
			if _, err := conn.Read(cmd); err != nil {
				log.Error(err)
				return
			}

			switch string(cmd) {
			case copyCmd:
				if err := clipboard.Copy(conn); err != nil {
					log.Error(err)
				}
			case pasteCmd:
				if err := clipboard.Paste(conn); err != nil {
					log.Error(err)
				}
			default:
				log.WithFields(log.Fields{"cmd": cmd}).Error("unexpected command")
			}
		}()
	}
}
