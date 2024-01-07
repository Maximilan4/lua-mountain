package filesystem

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"lua-mountain/pkg/option"
	"os"
	"path"
)

type (

	StorageConfig struct {
		Dir string
		Logger *slog.Logger
	}


	// TODO: lock system for preventing deleting/rewriting object
	Storage struct {
		Dir string
		Logger *slog.Logger
	}
)

func WithStorageLogger(logger *slog.Logger) option.ErrOption[*Storage] {
	return func(s *Storage) error {
		s.Logger = logger
		return nil
	}
}

func WithStorageConfig(cfg *StorageConfig) option.ErrOption[*Storage] {
	return func(s *Storage) error {
		s.Dir = cfg.Dir
		s.Logger = cfg.Logger
		return nil
	}
}

func NewStorage(opts ...option.ErrOption[*Storage]) (s *Storage, err error) {
	s = &Storage{}
	for _, opt := range opts {
		if err = opt(s); err != nil {
			return
		}
	}

	if s.Logger == nil {
		s.Logger = slog.Default()
	}

	if s.Dir == "" {
		return nil, errors.New("empty dir is not allowed")
	}

	if err = CreateDirIfNotExists(s.Dir); err != nil {
		return
	}

	return
}

func (s *Storage) Get(ctx context.Context, filename string) (io.ReadCloser, error) {
	fpath := path.Join(s.Dir, filename)
	s.Logger.DebugContext(ctx, "filesystem.Storage:Get() / os.Open()",
		slog.String("filepath", fpath),
	)

	return os.Open(fpath)
}

func (s *Storage) Exists(ctx context.Context, filename string) error {
	filepath := path.Join(s.Dir, filename)
	_, err := os.Stat(filepath)
	s.Logger.DebugContext(ctx, "filesystem.Storage:Exists()",
		slog.String("filepath", filepath),
	)

	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) Put(ctx context.Context, filename string, r io.Reader) error {
	fpath := path.Join(s.Dir, filename)
	s.Logger.DebugContext(ctx, "filesystem.Storage:Put() / os.Create()",
		slog.String("filepath", fpath),
	)

	file, err := os.Create(fpath)
	if err != nil {
		return err
	}

	defer file.Close()

	s.Logger.DebugContext(ctx, "filesystem.Storage:Put() / io.Copy()",
		slog.String("filepath", fpath),
	)

	if _, err = io.Copy(file, r); err != nil {
		return err
	}

	return nil
}

func (s *Storage) Delete(ctx context.Context, filename string) error {
	fpath := path.Join(s.Dir, filename)
	s.Logger.DebugContext(ctx, "filesystem.Storage:Delete() / os.Remove()",
		slog.String("filepath", fpath),
	)

	return os.Remove(fpath)
}

func (s *Storage) List(ctx context.Context) ([]string, error) {
	s.Logger.DebugContext(ctx, "reading a dir",
		slog.String("dir", s.Dir),
	)

	// TODO: read dir by batches
	entries, err := os.ReadDir(s.Dir)
	if err != nil {
		return nil, err
	}

	files := make([]string, 0, len(entries))
	var entry os.DirEntry
	for _, entry = range entries {
		if entry.IsDir() {
			s.Logger.DebugContext(ctx, "skip inner directory",
				slog.String("dir", entry.Name()),
			)
			continue
		}
		files = append(files, entry.Name())
	}

	return files, nil
}

