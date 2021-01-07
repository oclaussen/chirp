package command

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"

	"github.com/oclaussen/go-gimme/ssl"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewConfigureCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "configure",
		Short: "generates configuration files based on parameters",
		RunE: func(_ *cobra.Command, _ []string) error {
			return Configure()
		},
	}
}

func Configure() error {
	appdir := filepath.Join(os.Getenv("HOME"), fmt.Sprintf(".%s", name))

	u, err := url.Parse(viper.GetString(ConfKeyAddress))
	if err != nil {
		return fmt.Errorf("invalid address: %w", err)
	}

	_, files, err := ssl.GimmeCertificates(&ssl.Options{
		Org:          name,
		Hosts:        []string{u.Hostname()},
		WriteToFiles: &ssl.Files{Directory: appdir},
	})
	if err != nil {
		return fmt.Errorf("could not generate tls certificates: %w", err)
	}

	viper.Set(ConfKeyCAFile, files.CAFile)
	viper.Set(ConfKeyCertFile, files.ClientCertFile)
	viper.Set(ConfKeyKeyFile, files.ClientKeyFile)

	if err := viper.WriteConfigAs(filepath.Join(appdir, "client.yaml")); err != nil {
		return fmt.Errorf("could not write client config: %w", err)
	}

	viper.Set(ConfKeyCAFile, files.CAFile)
	viper.Set(ConfKeyCertFile, files.ServerCertFile)
	viper.Set(ConfKeyKeyFile, files.ServerKeyFile)

	if err := viper.WriteConfigAs(filepath.Join(appdir, "server.yaml")); err != nil {
		return fmt.Errorf("could not write server config: %w", err)
	}

	return nil
}
