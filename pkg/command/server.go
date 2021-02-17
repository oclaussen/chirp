package command

import (
	"errors"
	"fmt"

	"github.com/oclaussen/chirp/pkg/chirp"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewServerCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "server",
		Short: "Start in server mode and wait for incoming clipboard requests",
		RunE: func(_ *cobra.Command, _ []string) error {
			return withServer(func(s *chirp.ClipboardServer) error {
				return s.Listen()
			})
		},
	}
}

func withServer(f func(*chirp.ClipboardServer) error) error {
	log.SetFormatter(&log.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
	})

	viper.SetConfigName("server")

	if err := viper.ReadInConfig(); err != nil {
		var e *viper.ConfigFileNotFoundError
		if errors.As(err, &e) {
			log.Warn("no configuration file found: %w")
		} else {
			return fmt.Errorf("could not read config file: %w", err)
		}
	}

	server, err := chirp.NewServer(&chirp.Config{
		Address:         viper.GetString(ConfKeyAddress),
		CertificateFile: viper.GetString(ConfKeyCertFile),
		KeyFile:         viper.GetString(ConfKeyKeyFile),
		CAFile:          viper.GetString(ConfKeyCAFile),
	})
	if err != nil {
		return fmt.Errorf("cannot create chirp server: %w", err)
	}

	return f(server)
}
