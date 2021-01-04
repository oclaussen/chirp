package chirp

import (
	"context"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"

	api "github.com/oclaussen/chirp/api/v1"
	"google.golang.org/grpc"
)

type ClipboardClient struct {
	conn   *grpc.ClientConn
	client api.ClipboardServiceClient
}

func NewClient(socketType string, addr string) (*ClipboardClient, error) {
	if socketType == "unix" {
		addr, err := filepath.Abs(addr)
		if err != nil {
			return nil, fmt.Errorf("could not get socket path: %w", err)
		}

		if _, err := os.Stat(addr); err != nil {
			return nil, fmt.Errorf("could not get unix socket: %w", err)
		}
	}

	conn, err := grpc.Dial(
		addr,
		grpc.WithInsecure(),
		grpc.WithContextDialer(func(ctx context.Context, addr string) (net.Conn, error) {
			return net.Dial(socketType, addr)
		}),
	)
	if err != nil {
		return nil, fmt.Errorf("could not connect to server: %w", err)
	}

	return &ClipboardClient{
		conn:   conn,
		client: api.NewClipboardServiceClient(conn),
	}, nil
}

func (c *ClipboardClient) Close() {
	c.conn.Close()
}

func (c *ClipboardClient) Copy() error {
	bytes, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		return fmt.Errorf("could not read stdin data: %w", err)
	}

	if _, err := c.client.Copy(context.Background(), &api.CopyRequest{Contents: string(bytes)}); err != nil {
		return fmt.Errorf("could not send copy command: %w", err)
	}

	return nil
}

func (c *ClipboardClient) Paste() error {
	response, err := c.client.Paste(context.Background(), &api.PasteRequest{})
	if err != nil {
		return fmt.Errorf("could not send paste command: %w", err)
	}

	fmt.Print(response.Contents)

	return nil
}
