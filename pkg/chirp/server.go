package chirp

import (
	"context"
	"fmt"
	"net"
	"os"
	"path/filepath"

	api "github.com/oclaussen/chirp/api/v1"
	"github.com/oclaussen/chirp/pkg/clipboard"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type ClipboardServer struct{}

func (server *ClipboardServer) Copy(ctx context.Context, request *api.CopyRequest) (*api.CopyResponse, error) {
	if err := clipboard.CopyString(request.Contents); err != nil {
		return nil, fmt.Errorf("could not copy from clipboard: %w", err)
	}

	return &api.CopyResponse{}, nil
}

func (server *ClipboardServer) Paste(ctx context.Context, request *api.PasteRequest) (*api.PasteResponse, error) {
	contents, err := clipboard.PasteString()
	if err != nil {
		return nil, fmt.Errorf("could not paste to clipboard: %w", err)
	}

	return &api.PasteResponse{Contents: contents}, nil
}

func Server(socketType string, addr string) error {
	if err := clipboard.Check(); err != nil {
		return fmt.Errorf("cannot use clipboard: %w", err)
	}

	if socketType == "unix" {
		addr, err := filepath.Abs(addr)
		if err != nil {
			return fmt.Errorf("could not get socket path: %w", err)
		}

		if _, err := os.Stat(addr); !os.IsNotExist(err) {
			if _, err = net.Dial("unix", addr); err == nil {
				return fmt.Errorf("socket already exists at %s: %w", addr, err)
			}

			if err = os.Remove(addr); err != nil {
				return fmt.Errorf("could not remove stale socket at %s: %w", addr, err)
			}
		}
	}

	log.WithFields(log.Fields{"type": socketType, "address": addr}).Info("listening...")

	listener, err := net.Listen(socketType, addr)
	if err != nil {
		return fmt.Errorf("could not start server socket: %w", err)
	}

	defer listener.Close()

	grpcServer := grpc.NewServer()
	api.RegisterClipboardServiceServer(grpcServer, &ClipboardServer{})

	return grpcServer.Serve(listener)
}
