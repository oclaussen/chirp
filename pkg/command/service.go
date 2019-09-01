package command

import (
	"os"
	"os/user"
	"runtime"

	"github.com/kardianos/service"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
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
		Run: func(cmd *cobra.Command, args []string) {
			clientLogging()
			svc, err := newService()
			if err != nil {
				log.Error(err)
				os.Exit(1)
			}
			if err := svc.Install(); err != nil {
				log.Error(err)
				os.Exit(1)
			}
		},
	}
	return cmd
}

func NewServiceUninstallCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "uninstall",
		Short: "Uninstall chirp system service",
		Run: func(cmd *cobra.Command, args []string) {
			clientLogging()
			svc, err := newService()
			if err != nil {
				log.Error(err)
				os.Exit(1)
			}
			if err := svc.Uninstall(); err != nil {
				log.Error(err)
				os.Exit(1)
			}
		},
	}
	return cmd
}

func NewServiceStartCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start",
		Short: "Start chirp system service",
		Run: func(cmd *cobra.Command, args []string) {
			clientLogging()
			svc, err := newService()
			if err != nil {
				log.Error(err)
				os.Exit(1)
			}
			if err := svc.Start(); err != nil {
				log.Error(err)
				os.Exit(1)
			}
		},
	}
	return cmd
}

func NewServiceStopCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stop",
		Short: "Stop chirp system service",
		Run: func(cmd *cobra.Command, args []string) {
			clientLogging()
			svc, err := newService()
			if err != nil {
				log.Error(err)
				os.Exit(1)
			}
			if err := svc.Stop(); err != nil {
				log.Error(err)
				os.Exit(1)
			}
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
		Name:        qualified_name,
		DisplayName: name,
		Description: description,
		Option:      map[string]interface{}{},
		Arguments:   []string{"server", "--type", opts.socketType, "--address", opts.address},
	}

	if u, err := user.Current(); err == nil && u.Uid != "0" {
		config.UserName = u.Username
	}

	if runtime.GOOS == "darwin" {
		config.Option["UserService"] = true
	}

	return service.New(&program{}, config)
}
