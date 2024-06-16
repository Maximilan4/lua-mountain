package mountain

import (
	"context"
	"log/slog"
	"os"

	"github.com/urfave/cli/v2"

	"lua-mountain/internal/mountain/commands"
	"lua-mountain/internal/mountain/config"
	"lua-mountain/internal/mountain/logging"
)

var Engine *cli.App

func init() {
	Engine = &cli.App{
		Name:  "mountain",
		Usage: "mountain",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "config",
				Usage:    "--config=config.yaml",
				Required: false,
			},
		},
		Commands: []*cli.Command{
			commands.StartCommand(),
			commands.VersionCommand(),
		},
		Before:       onBefore,
		Action:       cli.ShowAppHelp,
		BashComplete: cli.DefaultAppComplete,
		Reader:       os.Stdin,
		Writer:       os.Stdout,
		ErrWriter:    os.Stderr,
	}

}

func Start(ctx context.Context, args []string) error {
	return Engine.RunContext(ctx, args)
}

func onBefore(c *cli.Context) error {
	configPath := c.String("config")
	var err error

	if configPath != "" {
		slog.Debug("loading config file", slog.String("path", configPath))
		if err = config.Load(configPath); err != nil {
			return err
		}

		return nil
	}

	slog.Debug("searching config file at", slog.Any("paths", config.DefaultSearchDirs))
	var p string
	if p, err = config.Search(config.DefaultSearchDirs...); err != nil {
		return err
	}

	slog.Debug("config founded, loading", slog.String("path", p))

	if err = config.Load(p); err != nil {
		return err
	}

	if err = logging.Init(&config.Get().Logs); err != nil {
		return err
	}

	return nil
}
