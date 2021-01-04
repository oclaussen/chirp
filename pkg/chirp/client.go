package chirp

import (
	"context"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
	"time"

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
			return nil, err
		}
		if _, err := os.Stat(addr); err != nil {
			return nil, err
		}
	}

	conn, err := grpc.Dial(
		addr,
		grpc.WithInsecure(),
		grpc.WithDialer(func(addr string, timeout time.Duration) (net.Conn, error) {
			return net.DialTimeout(socketType, addr, timeout)
		}),
	)
	if err != nil {
		return nil, err
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
		return err
	}

	if _, err := c.client.Copy(context.Background(), &api.CopyRequest{Contents: string(bytes)}); err != nil {
		return err
	}

	return nil
}

func (c *ClipboardClient) Paste() error {
	response, err := c.client.Paste(context.Background(), &api.PasteRequest{})
	if err != nil {
		return err
	}

	fmt.Print(response.Contents)
	return nil
}
