package repository

import (
	"log/slog"
	"lua-mountain/internal/mountain/storage"
)

const (
	defaultMaxFileSize = 2 << 30
)

type (
	Config struct {
		Prefix                string   `yaml:"prefix"`
		Storage               string   `yaml:"storage"`
		AllowedFileExtensions []string `yaml:"allowed_file_extensions"`
		AllowRewrite          bool     `yaml:"allow_rewrite"`
		MaxFileSize           uint64   `yaml:"max_file_size"`
	}

	Repository struct {
		Prefix                string
		Storage               storage.Storage
		logger                *slog.Logger
		AllowedFileExtensions []string
		AllowRewrite          bool
		MaxFileSize           uint64
	}
)

func New(cfg *Config, storage storage.Storage, logger *slog.Logger) *Repository {

	repo := &Repository{
		Storage:               storage,
		AllowRewrite:          cfg.AllowRewrite,
		MaxFileSize:           cfg.MaxFileSize,
		AllowedFileExtensions: cfg.AllowedFileExtensions,
		logger: logger.With(
			slog.String("prefix", cfg.Prefix),
		),
	}

	if repo.MaxFileSize == 0 {
		repo.MaxFileSize = defaultMaxFileSize
	}

	if repo.AllowedFileExtensions == nil {
		repo.AllowedFileExtensions = []string{".rockspec", ".src.rock", ".all.rock"}
	}

	repo.logger.Info("repo created",
		slog.Bool("rewrite", repo.AllowRewrite),
		slog.Uint64("max_file_size", repo.MaxFileSize),
		slog.Any("allowed_file_extensions", repo.AllowedFileExtensions),
	)

	return repo
}
