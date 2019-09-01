package command

import (
	"os"

	"github.com/oclaussen/chirp/pkg/chirp"
	log "github.com/sirupsen/logrus"
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
		Run: func(cmd *cobra.Command, args []string) {
			serverLogging()
			if err := chirp.Server(opts.socketType, opts.address); err != nil {
				log.Error(err)
				os.Exit(1)
			}
		},
	}
}

func NewCopyCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "copy",
		Short: "Send a copy request to the clipboard",
		Run: func(cmd *cobra.Command, args []string) {
			clientLogging()
			if err := chirp.Copy(opts.socketType, opts.address); err != nil {
				log.Error(err)
				os.Exit(1)
			}
		},
	}
}

func NewPasteCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "paste",
		Short: "Send a paste request to the clipboard",
		Run: func(cmd *cobra.Command, args []string) {
			clientLogging()
			if err := chirp.Paste(opts.socketType, opts.address); err != nil {
				log.Error(err)
				os.Exit(1)
			}
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
