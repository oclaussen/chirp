package chirp

import (
	"bytes"
	"context"
	"net"
	"os"
	"path/filepath"
	"strings"

	api "github.com/oclaussen/chirp/api/v1"
	"github.com/oclaussen/chirp/pkg/clipboard"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type ClipboardServer struct{}

func (server *ClipboardServer) Copy(ctx context.Context, request *api.CopyRequest) (*api.CopyResponse, error) {
	if err := clipboard.Copy(strings.NewReader(request.Contents)); err != nil {
		return nil, err
	}

	return &api.CopyResponse{}, nil
}

func (server *ClipboardServer) Paste(ctx context.Context, request *api.PasteRequest) (*api.PasteResponse, error) {
	var buf bytes.Buffer
	if err := clipboard.Paste(&buf); err != nil {
		return nil, err
	}

	return &api.PasteResponse{Contents: buf.String()}, nil
}

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

	grpcServer := grpc.NewServer()
	api.RegisterClipboardServiceServer(grpcServer, &ClipboardServer{})
	return grpcServer.Serve(listener)
}
