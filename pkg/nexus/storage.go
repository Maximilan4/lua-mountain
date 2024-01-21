package nexus

import (
	"context"
	"fmt"
	"io"
	"log/slog"
)

type (
	Storage struct {
		Client *HTTPClient
		Index  *AssetIndex
		logger *slog.Logger
	}
)

func (s *Storage) Get(ctx context.Context, filename string) (io.ReadCloser, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Storage) Exists(_ context.Context, filename string) error {
	if s.Index.Has(filename) {
		return nil
	}

	return fmt.Errorf("file %s not exists", filename)
}

func (s *Storage) Put(ctx context.Context, filename string, r io.Reader) error {
	//TODO implement me
	panic("implement me")
}

func (s *Storage) Delete(ctx context.Context, filename string) error {
	//TODO implement me
	panic("implement me")
}

// List - returns a slice of Asset.Path
func (s *Storage) List(_ context.Context) ([]string, error) {
	return s.Index.Keys(), nil
}
