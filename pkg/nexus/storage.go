package nexus

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"time"

	"lua-mountain/pkg/option"
)

const (
	DefaultIndexUpdateInterval = time.Minute
	maxAssetSearchRetry = 5
)

type (
	Storage struct {
		Client     *HTTPClient
		Index      *AssetIndex
		logger     *slog.Logger
		repository *Repository
		cfg        StorageConfig
	}

	StorageConfig struct {
		IndexUpdateInterval time.Duration
		RepositoryName      string
	}
)

func WithStorageLogger(l *slog.Logger) option.ErrOption[*Storage] {
	return func(s *Storage) error {
		s.logger = l
		return nil
	}
}

func WithStorageConfig(cfg StorageConfig) option.ErrOption[*Storage] {
	return func(s *Storage) error {
		s.cfg = cfg
		return nil
	}
}

func NewStorage(ctx context.Context, client *HTTPClient, opts ...option.ErrOption[*Storage]) (*Storage, error) {
	s := &Storage{Client: client, Index: NewAssetIndex()}
	var err error
	for _, opt := range opts {
		if err = opt(s); err != nil {
			return nil, err
		}
	}

	if s.logger == nil {
		s.logger = slog.Default()
	}

	if s.cfg.RepositoryName == "" {
		return nil, errors.New("storage: repository name is required for nexus storage configuration")
	}

	if s.cfg.IndexUpdateInterval == 0 {
		s.cfg.IndexUpdateInterval = DefaultIndexUpdateInterval
	}

	// TODO: make own ctx for this request only
	s.repository, err = s.Client.GetRepository(ctx, s.cfg.RepositoryName)
	if err != nil {
		return nil, fmt.Errorf("storage: unable to load repository %s data: %w", s.cfg.RepositoryName, err)
	}

	if s.repository.Format != RawRepositoryFormat {
		return nil, fmt.Errorf(
			"storage: repository format err: only raw format supported, got %s", s.repository.Format,
		)
	}

	// TODO: add condvar, for blocking methods until index will be built
	s.logger.Info("storage: building first index")
	if err = s.BuildIndex(ctx); err != nil {
		return nil, fmt.Errorf("storage: repository index build err: %w", err)
	}

	go s.UpdateIndexOnInterval(ctx, s.cfg.IndexUpdateInterval)

	return s, nil
}

func (s *Storage) UpdateIndexOnInterval(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	for {
		select {
		case tick := <-ticker.C:
			s.logger.InfoContext(ctx, "start of nexus index update")

			// TODO: may be make a timeout smaller than update interval (about ms values)
			uCtx, done := context.WithTimeout(ctx, interval)
			if err := s.BuildIndex(uCtx); err != nil {
				s.logger.ErrorContext(uCtx, "nexus index update error", slog.String("err", err.Error()))
				done()
				continue
			}

			done()
			s.logger.InfoContext(ctx, "end of nexus index update",
				slog.Time("next", tick.Add(interval)),
			)
		case <-ctx.Done():
			s.logger.Info("nexus index update stopped")
			return
		}
	}
}

func (s *Storage) BuildIndex(ctx context.Context) error {
	var (
		token  string
		assets = make(map[string]Asset, s.Index.Count())
	)

	for {
		l, err := s.Client.GetAssetsList(ctx, s.repository.Name, token)
		if err != nil {
			return err
		}

		for _, asset := range l.Items {
			assets[asset.Path] = asset
		}

		if l.Next == "" {
			break
		}

		token = l.Next
	}

	s.Index.Replace(assets)
	return nil
}

func (s *Storage) Get(ctx context.Context, filename string) (io.ReadCloser, error) {
	asset := s.Index.Get(filename)
	if asset == nil {
		return nil, fmt.Errorf("storage.Get(): %s not found in index: %w", filename, os.ErrNotExist)
	}

	req, err := http.NewRequest(http.MethodGet, asset.DownloadUrl, nil)
	if err != nil {
		return nil, fmt.Errorf("storage.Get(): http request build err: %w", err)
	}

	resp, err := s.Client.doRequest(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("storage.Get(): http request err: %w", err)
	}

	return resp.Body, nil
}

// Exists - check key in index
func (s *Storage) Exists(_ context.Context, filename string) error {
	if s.Index.Has(filename) {
		return nil
	}

	return fmt.Errorf("storage.Exists(): %s not found in index", filename)
}

// Put - saves file and content in storage
func (s *Storage) Put(ctx context.Context, filename string, r io.Reader) (err error) {
	err = s.Client.SaveAsset(ctx, s.cfg.RepositoryName, filename, r)
	if err != nil {
		err = fmt.Errorf("storage.Put() - unable to save asset: %w", err)
		return
	}

	var (
		l *SearchList
		delay time.Duration
		retry int
	)

	for retry < maxAssetSearchRetry {
		delay = time.Duration(float64(time.Second) * 0.1 * float64(retry))
		time.Sleep(delay)
		l, err = s.Client.SearchAssetByName(ctx, s.cfg.RepositoryName, filename, "")
		if err != nil {
			err = fmt.Errorf("storage.Put() - unable to search saved asset: %w", err)
			return
		}

		if len(l.Items) == 0 || len(l.Items[0].Assets) == 0 {
			s.logger.DebugContext(ctx, "uploaded asset not found",
				slog.String("asset", filename),
				slog.Int("attempt", retry + 1),
				slog.Duration("delay", delay),
			)
			retry++
			continue
		} else {
			s.logger.DebugContext(ctx, "uploaded asset found",
				slog.String("asset", filename),
				slog.Int("attempt", retry + 1),
				slog.Duration("delay", delay),
			)
			break
		}
	}

	asset := l.Items[0].Assets[0]
	s.Index.Store(filename, asset)
	return nil
}

// Delete - deletes an asset in nexus
func (s *Storage) Delete(ctx context.Context, filename string) (err error) {
	asset := s.Index.Get(filename)
	if asset == nil {
		return
	}

	err = s.Client.DeleteAsset(ctx, asset.Id)
	if err != nil {
		err = fmt.Errorf("storage.Delete() - %w", err)
		return
	}

	s.Index.Delete(filename)

	return
}

// List - returns a slice of Asset.Path
func (s *Storage) List(_ context.Context) ([]string, error) {
	return s.Index.Keys(), nil
}
