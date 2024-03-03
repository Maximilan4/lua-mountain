package storage

import (
	"log/slog"

	"lua-mountain/pkg/attr"
	"lua-mountain/pkg/filesystem"
)

const (
	DefaultStorageDir = "/var/mountain"
)

func InitFsStorage(name string, cfg map[string]any, logger *slog.Logger) (*filesystem.Storage, error) {
	var (
		sCfg = filesystem.StorageConfig{}
		ok   bool
	)

	sCfg.Dir, ok = attr.GetTyped[string](cfg, "dir")
	if !ok {
		sCfg.Dir = DefaultStorageDir
	}
	sCfg.Logger = logger.With(slog.String("storage", name), slog.String("dir", sCfg.Dir))
	sCfg.Logger.Info("loading new fs storage")

	return filesystem.NewStorage(filesystem.WithStorageConfig(&sCfg))
}
