package chirp

import (
	"context"
	"fmt"
	"net"

	api "github.com/oclaussen/chirp/api/v1"
	"github.com/oclaussen/chirp/pkg/clipboard"
	"google.golang.org/grpc"
)

type ClipboardServer struct {
	listener net.Listener
	server   *grpc.Server
}

func (s *ClipboardServer) Copy(ctx context.Context, request *api.CopyRequest) (*api.CopyResponse, error) {
	if err := clipboard.CopyString(request.Contents); err != nil {
		return nil, fmt.Errorf("could not copy from clipboard: %w", err)
	}

	return &api.CopyResponse{}, nil
}

func (s *ClipboardServer) Paste(ctx context.Context, request *api.PasteRequest) (*api.PasteResponse, error) {
	contents, err := clipboard.PasteString()
	if err != nil {
		return nil, fmt.Errorf("could not paste to clipboard: %w", err)
	}

	return &api.PasteResponse{Contents: contents}, nil
}

func NewServer(config *Config) (*ClipboardServer, error) {
	if err := clipboard.Check(); err != nil {
		return nil, fmt.Errorf("cannot use clipboard: %w", err)
	}

	protocol, addr, err := config.DialOptions()
	if err != nil {
		return nil, fmt.Errorf("invalid connection config: %w", err)
	}

	if _, err = net.Dial(protocol, addr); err == nil {
		return nil, fmt.Errorf("server already exists at %s: %w", addr, err)
	}

	creds, err := config.TLSServerOptions()
	if err != nil {
		return nil, err
	}

	listener, err := net.Listen(protocol, addr)
	if err != nil {
		return nil, fmt.Errorf("could not start server socket: %w", err)
	}

	return &ClipboardServer{
		listener: listener,
		server:   grpc.NewServer(grpc.Creds(creds)),
	}, nil
}

func (s *ClipboardServer) Listen() error {
	defer s.listener.Close()

	api.RegisterClipboardServiceServer(s.server, s)

	return s.server.Serve(s.listener)
}
