package command

import (
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
		RunE: func(_ *cobra.Command, _ []string) error {
			return withService(func(s service.Service) error {
				return s.Install()
			})
		},
	}

	return cmd
}

func NewServiceUninstallCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "uninstall",
		Short: "Uninstall chirp system service",
		RunE: func(_ *cobra.Command, _ []string) error {
			return withService(func(s service.Service) error {
				return s.Uninstall()
			})
		},
	}

	return cmd
}

func NewServiceStartCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start",
		Short: "Start chirp system service",
		RunE: func(_ *cobra.Command, _ []string) error {
			return withService(func(s service.Service) error {
				return s.Start()
			})
		},
	}

	return cmd
}

func NewServiceStopCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stop",
		Short: "Stop chirp system service",
		RunE: func(_ *cobra.Command, _ []string) error {
			return withService(func(s service.Service) error {
				return s.Stop()
			})
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
		Arguments:   []string{"server"},
	}

	if u, err := user.Current(); err == nil && u.Uid != "0" {
		config.UserName = u.Username
	}

	if runtime.GOOS == "darwin" {
		config.Option["UserService"] = true
	}

	return service.New(&program{}, config)
}

func withService(f func(service.Service) error) error {
	log.SetFormatter(&log.TextFormatter{
		DisableTimestamp:       true,
		DisableLevelTruncation: true,
	})

	svc, err := newService()
	if err != nil {
		return err
	}

	return f(svc)
}
