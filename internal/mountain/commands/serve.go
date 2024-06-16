package commands

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/urfave/cli/v2"

	"lua-mountain/internal/mountain/config"
	"lua-mountain/internal/mountain/logging"
	"lua-mountain/internal/mountain/repository"
	"lua-mountain/internal/mountain/server"
	"lua-mountain/internal/mountain/server/mw"
	"lua-mountain/internal/mountain/storage"
)

func StartCommand() *cli.Command {
	return &cli.Command{
		Name:        "serve",
		Usage:       "mountain serve",
		Description: "starts a rocks server",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "address",
				Category: "http",
				Usage:    "--address 0.0.0.0",
			},
			&cli.StringFlag{
				Name:     "port",
				Category: "http",
				Usage:    "--port 2023",
			},
		},
		Category: "",
		Action:   startRocksServer,
	}
}

func startRocksServer(c *cli.Context) error {
	cfg := config.Get()
	srv := server.Init()
	storages := storage.InitStorages(context.Background(), cfg.Storages, logging.DefaultLogger)
	for _, repoCfg := range cfg.Repositories {
		st, ok := storages[repoCfg.Storage]
		if !ok {
			slog.Warn("unable to find repository storage",
				slog.String("repository", repoCfg.Prefix),
				slog.String("storage", repoCfg.Storage),
			)
			continue
		}

		repo := repository.New(&repoCfg, st, logging.DefaultLogger)
		extMw := mw.AllowedExtensions(repo.AllowedFileExtensions)
		rGroup := srv.Group(repoCfg.Prefix)
		for _, man := range []string{"manifest", "manifest-5.1", "manifest-5.2", "manifest-5.3", "manifest-5.4"} {
			rGroup.GET("/"+man, repo.GetManifest)
			rGroup.GET("/"+man+".json", repo.GetManifestJson)
			rGroup.GET("/"+man+".zip", repo.GetManifestZip)
		}

		rGroup.GET("/:filename", repo.Get, extMw)
		rGroup.PUT("/:filename", repo.Put, extMw)
		rGroup.DELETE("/:filename", repo.Delete, extMw)
	}

	for _, r := range srv.Routes() {
		if r.Method == "echo_route_not_found" {
			continue
		}

		logging.DefaultLogger.Info("mountain`s route inited",
			slog.String("method", r.Method),
			slog.String("path", r.Path),
		)
	}

	var (
		address = c.String("address")
		port    = c.String("port")
	)

	if address == "" {
		address = cfg.Listen.Address
	}

	if port == "" {
		port = cfg.Listen.Port
	}

	bindAddress := fmt.Sprintf("%s:%s", address, port)
	logging.DefaultLogger.Info("starting mountain on",
		slog.String("address", bindAddress),
	)

	if err := srv.Start(bindAddress); err != nil {
		return err
	}

	return nil
}
