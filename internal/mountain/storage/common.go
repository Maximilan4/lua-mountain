package storage

import (
	"context"
	"io"
	"log/slog"
	"lua-mountain/pkg/attr"
	"lua-mountain/pkg/filesystem"
)

const (
	DefaultStorageDir = "/var/mountain"
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


func InitStorages(cfg map[string]any, logger *slog.Logger) Storages {
	var (
		storages = make(Storages, len(cfg))
		err error
	)

	for name, storageCfg := range cfg {
		switch storageCfg := storageCfg.(type) {
		case map[string]any:
			var (
				t string
				ok bool
			)

			if t, ok = attr.GetTyped[string](storageCfg, "type"); !ok {
				slog.Warn("skip storage init: bad type", slog.String("name", name))
				continue
			}

			switch t {
			case "fs":
				storages[name], err = InitFsStorage(name, storageCfg, logger)
				if err != nil {
					delete(storages, name)
					slog.Error("storage initialization err", slog.String("err", err.Error()))
					continue
				}

			}


		default:
			slog.Warn("skip storage init: bad configuration", slog.String("name", name))
		}
	}

	return storages
}

func InitFsStorage(name string, cfg map[string]any, logger *slog.Logger) (*filesystem.Storage, error ){
	var (
		sCfg = filesystem.StorageConfig{}
		ok bool
	)

	sCfg.Dir, ok = attr.GetTyped[string](cfg, "dir")
	if !ok {
		sCfg.Dir = DefaultStorageDir
	}
	sCfg.Logger = logger.With(slog.String("storage", name), slog.String("dir", sCfg.Dir))
	sCfg.Logger.Info("loading new fs storage")

	return filesystem.NewStorage(filesystem.WithStorageConfig(&sCfg))
}