package commands

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"log/slog"
	"lua-mountain/internal/mountain/config"
	"lua-mountain/internal/mountain/logging"
	"lua-mountain/internal/mountain/repository"
	"lua-mountain/internal/mountain/server"
	"lua-mountain/internal/mountain/server/mw"
	"lua-mountain/internal/mountain/storage"
)

func StartCommand() *cli.Command {
	return &cli.Command{
		Name:                   "serve",
		Usage:                  "mountain serve",
		Description:            "starts a rocks server",
		Category:               "",
		Action:                 startRocksServer,
	}
}

func startRocksServer(c *cli.Context) error {
	cfg := config.Get()
	srv := server.Init()
	storages := storage.InitStorages(cfg.Storages, logging.DefaultLogger)
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

	address := fmt.Sprintf("%s:%s", cfg.Listen.Address, cfg.Listen.Port)
	logging.DefaultLogger.Info("starting mountain on",
		slog.String("address", address),
	)

	if err := srv.Start(address); err != nil {
		return err
	}

	return nil
}