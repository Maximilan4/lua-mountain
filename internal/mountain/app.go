package mountain

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"runtime/debug"

	"github.com/urfave/cli/v2"

	"lua-mountain/internal/mountain/commands"
	"lua-mountain/internal/mountain/config"
	"lua-mountain/internal/mountain/logging"
)

var Engine *cli.App

func init() {
	cli.VersionPrinter = versionPrinter
	Engine = &cli.App{
		Name:  "mountain",
		Usage: "mountain",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "config",
				Usage:    "--config=config.yaml",
				Aliases:  []string{"c"},
				Required: false,
			},
			&cli.BoolFlag{
				Name:     "debug",
				Aliases:  []string{"d"},
				Usage:    "--debug",
				Required: false,
			},
		},
		Commands: []*cli.Command{
			commands.StartCommand(),
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

	if c.Bool("debug") {
		slog.SetLogLoggerLevel(slog.LevelDebug)
	}

	return nil
}

func versionPrinter(ctx *cli.Context) {
	base := fmt.Sprintf("%s; version %s; ", ctx.App.Name, ctx.App.Version)
	d, ok := debug.ReadBuildInfo()
	if !ok {
		fmt.Println(base)
		return
	}

	base += fmt.Sprintf("%s; ", d.GoVersion)
	settings := make(map[string]string, len(d.Settings))
	for _, setting := range d.Settings {
		switch setting.Key {
		case "vcs.revision":
			base += fmt.Sprintf("revision %s; ", setting.Value[:8])
		case "vcs.time":
			base += fmt.Sprintf("time %s;", setting.Value)
		default:
			settings[setting.Key] = setting.Value
		}
	}

	fmt.Println(base)
	if ctx.Bool("debug") {
		fmt.Println("build info:")
		for k, v := range settings {
			fmt.Printf("- %s: %s\n", k, v)
		}

	}

	return
}
