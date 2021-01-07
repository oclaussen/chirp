package command

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	name          = "chirp"
	qualifiedName = "com.wabenet.chirp"
	description   = "Access system clipboard over network"

	ConfKeyAddress  = "address"
	ConfKeyCAFile   = "tls-ca-file"
	ConfKeyCertFile = "tls-cert-file"
	ConfKeyKeyFile  = "tls-key-file"

	DefaultAddress = "tcp://127.0.0.1:24709"
)

func Execute() int {
	viper.SetDefault(ConfKeyAddress, DefaultAddress)

	viper.SetEnvPrefix(name)

	viper.SetConfigType("yaml")
	viper.AddConfigPath(fmt.Sprintf("/etc/%s", name))
	viper.AddConfigPath(fmt.Sprintf("$HOME/.%s", name))

	cmd := &cobra.Command{Use: name, Short: description}
	cmd.AddCommand(NewServerCommand())
	cmd.AddCommand(NewCopyCommand())
	cmd.AddCommand(NewPasteCommand())
	cmd.AddCommand(NewConfigureCommand())
	cmd.AddCommand(NewServiceCommand())

	flags := cmd.PersistentFlags()
	flags.StringP(ConfKeyAddress, "a", "", "address to bind to")
	flags.String(ConfKeyCertFile, "", "tls certificate file")
	flags.String(ConfKeyKeyFile, "", "tls key file")
	flags.String(ConfKeyCAFile, "", "tls certificate authority file")

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
