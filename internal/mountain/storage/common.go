package storage

import (
	"context"
	"io"
	"log/slog"

	"lua-mountain/pkg/attr"
)

type (
	Storage interface {
		Get(ctx context.Context, filename string) (io.ReadCloser, error)
		Exists(ctx context.Context, filename string) error
		Put(ctx context.Context, filename string, r io.Reader) error
		Delete(ctx context.Context, filename string) error
		List(ctx context.Context) ([]string, error)
	}

	Storages map[string]Storage
)

func InitStorages(ctx context.Context, cfg map[string]any, logger *slog.Logger) Storages {
	var (
		storages = make(Storages, len(cfg))
		err      error
	)

	for name, storageCfg := range cfg {
		switch storageCfg := storageCfg.(type) {
		case map[string]any:
			var (
				t  string
				ok bool
			)

			if t, ok = attr.GetTyped[string](storageCfg, "type"); !ok {
				logger.Warn("skip storage init: bad type", slog.String("name", name))
				continue
			}

			switch t {
			case "fs":
				storages[name], err = InitFsStorage(name, storageCfg, logger)
				if err != nil {
					delete(storages, name)
					logger.Error("storage initialization err", slog.String("err", err.Error()))
					continue
				}
			case "nexus":
				storages[name], err = InitNexusStorage(ctx, name, storageCfg, logger)
				if err != nil {
					delete(storages, name)
					logger.Error("storage initialization err", slog.String("err", err.Error()))
					continue
				}
			default:
				logger.Warn("skip storage init: unexpected type", slog.String("type", t))
			}

		default:
			logger.Warn("skip storage init: bad configuration", slog.String("name", name))
		}
	}

	return storages
}
