package command

import (
	"github.com/oclaussen/chirp/pkg/chirp"
	"github.com/spf13/cobra"
)

type options struct {
	socketType string
	address    string
}

var opts options

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "chirp",
		Short: "Access system clipboard over network",
	}

	cmd.AddCommand(NewServerCommand())
	cmd.AddCommand(NewCopyCommand())
	cmd.AddCommand(NewPasteCommand())

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
			return chirp.Server(opts.socketType, opts.address)
		},
	}
}

func NewCopyCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "copy",
		Short: "Send a copy request to the clipboard",
		RunE: func(cmd *cobra.Command, args []string) error {
			return chirp.Copy(opts.socketType, opts.address)
		},
	}
}

func NewPasteCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "paste",
		Short: "Send a paste request to the clipboard",
		RunE: func(cmd *cobra.Command, args []string) error {
			return chirp.Paste(opts.socketType, opts.address)
		},
	}
}
