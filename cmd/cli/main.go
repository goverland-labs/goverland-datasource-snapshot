package main

import (
	"os"

	"github.com/caarlos0/env/v8"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/s-larionov/process-manager"

	"github.com/goverland-labs/datasource-snapshot/internal"
	"github.com/goverland-labs/datasource-snapshot/internal/config"
	"github.com/goverland-labs/datasource-snapshot/internal/logger"
)

var (
	cfg config.App
)

func init() {
	err := env.Parse(&cfg)
	if err != nil {
		panic(err)
	}

	level, err := zerolog.ParseLevel(cfg.LogLevel)
	if err != nil {
		panic(err)
	}
	zerolog.SetGlobalLevel(level)
	process.SetLogger(&logger.ProcessManagerLogger{})
}

func main() {
	application, err := internal.NewCliApplication(cfg)
	if err != nil {
		panic(err)
	}

	args := os.Args
	if len(args) <= 1 {
		application.PrintUsage()
		os.Exit(0)
	}

	cmd := args[1]
	params := args[2:]

	if err := application.ExecCommand(cmd, params...); err != nil {
		log.Err(err).
			Fields(map[string]any{
				"command":   cmd,
				"arguments": os.Args,
			}).
			Msg("unable to execute command")
		os.Exit(1)
	}
}
