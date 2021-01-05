package command

import (
	"os/user"
	"runtime"

	"github.com/kardianos/service"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewServiceCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "service",
		Short: "Manage the chirp server as a system service",
	}

	cmd.AddCommand(NewServiceInstallCommand())
	cmd.AddCommand(NewServiceUninstallCommand())
	cmd.AddCommand(NewServiceStartCommand())
	cmd.AddCommand(NewServiceStopCommand())

	return cmd
}

func NewServiceInstallCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "install",
		Short: "Install chirp server as a system service",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientLogging()

			svc, err := newService()
			if err != nil {
				return err
			}

			return svc.Install()
		},
	}

	return cmd
}

func NewServiceUninstallCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "uninstall",
		Short: "Uninstall chirp system service",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientLogging()

			svc, err := newService()
			if err != nil {
				return err
			}

			return svc.Uninstall()
		},
	}

	return cmd
}

func NewServiceStartCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start",
		Short: "Start chirp system service",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientLogging()

			svc, err := newService()
			if err != nil {
				return err
			}

			return svc.Start()
		},
	}

	return cmd
}

func NewServiceStopCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stop",
		Short: "Stop chirp system service",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientLogging()

			svc, err := newService()
			if err != nil {
				return err
			}

			return svc.Stop()
		},
	}

	return cmd
}

type program struct{}

func (p *program) Start(_ service.Service) error {
	return nil
}

func (p *program) Stop(_ service.Service) error {
	return nil
}

func newService() (service.Service, error) {
	config := &service.Config{
		Name:        qualifiedName,
		DisplayName: name,
		Description: description,
		Option:      map[string]interface{}{},
		Arguments:   []string{"server", "--type", viper.GetString("type"), "--address", viper.GetString("address")},
	}

	if u, err := user.Current(); err == nil && u.Uid != "0" {
		config.UserName = u.Username
	}

	if runtime.GOOS == "darwin" {
		config.Option["UserService"] = true
	}

	return service.New(&program{}, config)
}
