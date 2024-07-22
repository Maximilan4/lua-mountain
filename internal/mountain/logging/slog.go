package logging

import (
	"log/slog"
	"os"
	"strings"

	"lua-mountain/internal/mountain/repository"
	"lua-mountain/pkg/slogan"
)

type (
	Config struct {
		Target string     `yaml:"target"`
		Format string     `yaml:"format"`
		Level  slog.Level `yaml:"level"`
	}
)

var (
	DefaultLogger = slog.Default()
)

func Init(cfg *Config) error {
	var (
		err     error
		handler slog.Handler
		target  *os.File
	)

	if cfg.Target == "" {
		target = os.Stdout
	} else if target, err = os.Open(cfg.Target); err != nil {
		return err
	}

	switch strings.ToLower(cfg.Format) {
	case "json":
		handler = slogan.NewJSONHandler(target, &slog.HandlerOptions{
			Level: cfg.Level,
		}, repository.RequestIdContextKey)
	case "text":
		fallthrough
	default:
		handler = slogan.NewTextHandler(target, &slog.HandlerOptions{
			Level: cfg.Level,
		}, repository.RequestIdContextKey)
	}

	DefaultLogger = slog.New(handler)
	slog.SetDefault(DefaultLogger)
	return nil
}
