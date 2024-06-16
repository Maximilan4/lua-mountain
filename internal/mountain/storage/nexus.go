package storage

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"lua-mountain/pkg/attr"
	"lua-mountain/pkg/nexus"
)

func InitNexusStorage(ctx context.Context, name string, cfg map[string]any, logger *slog.Logger) (*nexus.Storage, error) {
	var (
		sCfg = nexus.StorageConfig{}
		cCfg = nexus.HTTPClientConfig{}
		ok   bool
		err  error
	)

	cCfg.Address, ok = attr.GetTyped[string](cfg, "address")
	if !ok || cCfg.Address == "" {
		return nil, errors.New("nexus storage init err: address is required")
	}

	cCfg.Timeout, err = attr.GetDuration(cfg, "request_timeout")
	if err != nil {
		logger.Warn("config key parse err", slog.String("err", err.Error()))
		cCfg.Timeout = nexus.DefaultClientTimeout
	}

	sCfg.IndexUpdateInterval, err = attr.GetDuration(cfg, "index_update_interval")
	if err != nil {
		logger.Warn("config key parse err", slog.String("err", err.Error()))
		sCfg.IndexUpdateInterval = nexus.DefaultIndexUpdateInterval
	}

	sCfg.RepositoryName, ok = attr.GetTyped[string](cfg, "repository")
	if !ok || sCfg.RepositoryName == "" {
		return nil, errors.New("nexus storage init err: repository is required")
	}

	nLogger := logger.With(slog.String("storage", name), slog.String("repo", sCfg.RepositoryName))
	cCfg.Logger = nLogger

	client, err := nexus.NewHTTPClient(&cCfg)
	if err != nil {
		return nil, fmt.Errorf("nexus storage init err: %w", err)
	}

	return nexus.NewStorage(ctx, client,
		nexus.WithStorageLogger(nLogger),
		nexus.WithStorageConfig(sCfg),
	)
}
