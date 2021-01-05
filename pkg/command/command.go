package command

import (
	"fmt"

	"github.com/oclaussen/chirp/pkg/chirp"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	name          = "chirp"
	qualifiedName = "com.wabenet.chirp"
	description   = "Access system clipboard over network"
)

func Execute() int {
	viper.SetDefault("type", "tcp")
	viper.SetDefault("address", "127.0.0.1:24709")

	viper.SetEnvPrefix(name)

	viper.SetConfigName(name)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(fmt.Sprintf("/etc/%s", name))
	viper.AddConfigPath(fmt.Sprintf("$HOME/.%s", name))

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Warn(fmt.Errorf("no configuration file found: %w", err))
		} else {
			log.Error(fmt.Errorf("could not read config file: %w", err))

			return 1
		}
	}

	cmd := &cobra.Command{Use: name, Short: description}
	cmd.AddCommand(NewServerCommand())
	cmd.AddCommand(NewCopyCommand())
	cmd.AddCommand(NewPasteCommand())
	cmd.AddCommand(NewServiceCommand())

	flags := cmd.PersistentFlags()
	flags.StringP("type", "t", "", "socket type to listen on (tcp or unix)")
	flags.StringP("address", "a", "", "address to bind to")

	if err := viper.BindPFlags(flags); err != nil {
		log.Error(err)

		return 1
	}

	if err := cmd.Execute(); err != nil {
		log.Error(err)

		return 1
	}

	return 0
}

func NewServerCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "server",
		Short: "Start in server mode and wait for incoming clipboard requests",
		RunE: func(cmd *cobra.Command, args []string) error {
			serverLogging()

			return chirp.Server(viper.GetString("type"), viper.GetString("address"))
		},
	}
}

func NewCopyCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "copy",
		Short: "Send a copy request to the clipboard",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientLogging()

			client, err := chirp.NewClient(viper.GetString("type"), viper.GetString("address"))
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

			client, err := chirp.NewClient(viper.GetString("type"), viper.GetString("address"))
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
