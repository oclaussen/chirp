package chirp

import (
	"context"
	"fmt"
        "io"
	"os"

	api "github.com/oclaussen/chirp/api/v1"
	"google.golang.org/grpc"
)

type ClipboardClient struct {
	conn   *grpc.ClientConn
	client api.ClipboardServiceClient
}

func NewClient(config *Config) (*ClipboardClient, error) {
	creds, err := config.TLSClientOptions()
	if err != nil {
		return nil, err
	}

	conn, err := grpc.NewClient(
		config.Address,
		grpc.WithTransportCredentials(creds),
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
	bytes, err := io.ReadAll(os.Stdin)
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

	fmt.Fprint(os.Stdout, response.Contents)

	return nil
}
