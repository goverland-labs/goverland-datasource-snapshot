package internal

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"github.com/s-larionov/process-manager"

	"github.com/goverland-labs/goverland-datasource-snapshot/internal/cli"
	"github.com/goverland-labs/goverland-datasource-snapshot/internal/config"
)

var (
	ErrCommandNotFound = errors.New("command not found")
)

type CliApplication struct {
	Application
	comLock  sync.RWMutex
	commands []cli.Command
}

func NewCliApplication(cfg config.App) (*CliApplication, error) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	a := &CliApplication{
		Application: Application{
			sigChan:   sigChan,
			cfg:       cfg,
			manager:   process.NewManager(),
			isCliMode: true,
		},
	}

	err := a.bootstrap()
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (a *CliApplication) RegisterCommand(cmd cli.Command) error {
	a.comLock.Lock()
	defer a.comLock.Unlock()

	a.commands = append(a.commands, cmd)

	return nil
}

func (a *CliApplication) ExecCommand(cmd string, args ...string) error {
	a.comLock.RLock()
	defer a.comLock.RUnlock()

	for _, c := range a.commands {
		if c.GetName() != cmd {
			continue
		}

		a, err := c.ParseArgs(args...)
		if err != nil {
			return err
		}

		return c.Execute(a)
	}

	return fmt.Errorf("%w: %q", ErrCommandNotFound, cmd)
}

func (a *CliApplication) PrintUsage() {
	a.comLock.RLock()
	defer a.comLock.RUnlock()

	usage := `Usage:
  cli command [param1] [param2] ...

Commands:`

	for _, cmd := range a.commands {
		usage += "\n" + a.prepareCommandUsage(cmd)
	}

	fmt.Println(usage)
}

func (a *CliApplication) prepareCommandUsage(cmd cli.Command) string {
	params := make([]string, 0)
	details := make([]string, 0)

	for name, desc := range cmd.GetArguments() {
		params = append(params, fmt.Sprintf("--%s=...", name))
		details = append(details, fmt.Sprintf("%s - %s", name, desc))
	}

	return fmt.Sprintf("%s %s\n  %s", cmd.GetName(), strings.Join(params, " "), strings.Join(details, "\n  "))
}

func (a *CliApplication) bootstrap() error {
	if err := a.Application.bootstrap(); err != nil {
		return err
	}

	if err := a.RegisterCommand(&cli.Import{
		Spaces:    a.spacesService,
		Proposals: a.proposalsService,
		Votes:     a.votesService,
	}); err != nil {
		return err
	}

	return nil
}
