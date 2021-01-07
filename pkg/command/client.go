package command

import (
	"fmt"

	"github.com/oclaussen/chirp/pkg/chirp"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewCopyCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "copy",
		Short: "Send a copy request to the clipboard",
		RunE: func(_ *cobra.Command, _ []string) error {
			return withClient(func(c *chirp.ClipboardClient) error {
				return c.Copy()
			})
		},
	}
}

func NewPasteCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "paste",
		Short: "Send a paste request to the clipboard",
		RunE: func(_ *cobra.Command, _ []string) error {
			return withClient(func(c *chirp.ClipboardClient) error {
				return c.Paste()
			})
		},
	}
}

func withClient(f func(*chirp.ClipboardClient) error) error {
	log.SetFormatter(&log.TextFormatter{
		DisableTimestamp:       true,
		DisableLevelTruncation: true,
	})

	viper.SetConfigName("client")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Warn("no configuration file found: %w")
		} else {
			return fmt.Errorf("could not read config file: %w", err)
		}
	}

	client, err := chirp.NewClient(&chirp.Config{
		Address:         viper.GetString(ConfKeyAddress),
		CertificateFile: viper.GetString(ConfKeyCertFile),
		KeyFile:         viper.GetString(ConfKeyKeyFile),
		CAFile:          viper.GetString(ConfKeyCAFile),
	})
	if err != nil {
		return fmt.Errorf("cannot create chirp client: %w", err)
	}

	defer client.Close()

	return f(client)
}
