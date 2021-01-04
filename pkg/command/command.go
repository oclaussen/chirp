package command

import (
	"fmt"

	"github.com/oclaussen/chirp/pkg/chirp"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const (
	name          = "chirp"
	qualifiedName = "com.wabenet.chirp"
	description   = "Access system clipboard over network"
)

var opts options

type options struct {
	socketType string
	address    string
}

func Execute() int {
	if err := NewCommand().Execute(); err != nil {
		log.Error(err)

		return 1
	}

	return 0
}

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{Use: name, Short: description}
	cmd.AddCommand(NewServerCommand())
	cmd.AddCommand(NewCopyCommand())
	cmd.AddCommand(NewPasteCommand())
	cmd.AddCommand(NewServiceCommand())

	flags := cmd.PersistentFlags()
	flags.StringVarP(&opts.socketType, "type", "t", "tcp", "socket type to listen on (tcp or unix)")
	flags.StringVarP(&opts.address, "address", "a", "127.0.0.1:24709", "address to bind to")

	return cmd
}

func NewServerCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "server",
		Short: "Start in server mode and wait for incoming clipboard requests",
		RunE: func(cmd *cobra.Command, args []string) error {
			serverLogging()

			return chirp.Server(opts.socketType, opts.address)
		},
	}
}

func NewCopyCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "copy",
		Short: "Send a copy request to the clipboard",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientLogging()

			client, err := chirp.NewClient(opts.socketType, opts.address)
			if err != nil {
				return fmt.Errorf("cannot create chirp client: %w", err)
			}

			return client.Copy()
		},
	}
}

func NewPasteCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "paste",
		Short: "Send a paste request to the clipboard",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientLogging()

			client, err := chirp.NewClient(opts.socketType, opts.address)
			if err != nil {
				return fmt.Errorf("cannot create chirp client: %w", err)
			}

			return client.Paste()
		},
	}
}

func serverLogging() {
	log.SetFormatter(&log.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
	})
}

func clientLogging() {
	log.SetFormatter(&log.TextFormatter{
		DisableTimestamp:       true,
		DisableLevelTruncation: true,
	})
}
