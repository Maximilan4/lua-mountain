package nexus

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"time"
)

const (
	RestAPIPrefix        = "service/rest/v1"
	DefaultClientTimeout = time.Second * 30
)

type (
	HTTPClientConfig struct {
		Address string
		Timeout time.Duration
		Logger  *slog.Logger
	}

	// HTTPClient - minimal needed nexus http client
	HTTPClient struct {
		*http.Client
		Url            *url.URL
		logger         *slog.Logger
		RequestTimeout time.Duration
	}
)

func (hc *HTTPClient) DeleteAsset(ctx context.Context, id string) error {
	var (
		err  error
		req  *http.Request
		resp *http.Response
	)

	u := hc.Url.JoinPath(RestAPIPrefix, "assets", id)

	req, err = http.NewRequest(http.MethodDelete, u.String(), nil)
	if err != nil {
		return fmt.Errorf("DeleteAsset build request err: %w", err)
	}

	tCtx, done := context.WithTimeout(ctx, hc.RequestTimeout)
	defer done()

	resp, err = hc.doRequest(tCtx, req)
	if err != nil {
		return fmt.Errorf("DeleteAsset http req err: %w", err)
	}

	defer resp.Body.Close()
	return nil
}

func (hc *HTTPClient) SaveAsset(ctx context.Context, repository, filepath string, f io.Reader) error {
	// TODO: need to add limit param for a single page when nexus will support it
	var (
		err  error
		req  *http.Request
		resp *http.Response
	)

	u := hc.Url.JoinPath("repository", repository, filepath)

	req, err = http.NewRequest(http.MethodPut, u.String(), f)
	if err != nil {
		return fmt.Errorf("SaveAsset err: %w", err)
	}

	resp, err = hc.doRequest(ctx, req)
	if err != nil {
		return fmt.Errorf("SaveAsset http req err: %w", err)
	}

	defer resp.Body.Close()
	return nil
}

func (hc *HTTPClient) SearchAssetByName(ctx context.Context, repository, filename string, nextToken string) (*SearchList, error) {
	var (
		err  error
		req  *http.Request
		resp *http.Response
	)

	u := hc.Url.JoinPath(RestAPIPrefix, "search")
	q := u.Query()
	if nextToken != "" {
		q.Add("continuationToken", nextToken)
	}

	q.Add("repository", repository)
	q.Add("name", filename)
	u.RawQuery = q.Encode()

	req, err = http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("SearchAssetByName build request err: %w", err)
	}

	tCtx, done := context.WithTimeout(ctx, hc.RequestTimeout)
	defer done()

	resp, err = hc.doRequest(tCtx, req)
	if err != nil {
		return nil, fmt.Errorf("SearchAssetByName http req err: %w", err)
	}

	defer resp.Body.Close()
	return jsonTo[SearchList](resp.Body)
}

func (hc *HTTPClient) GetAssetsList(ctx context.Context, repository string, nextToken string) (*AssetList, error) {
	// TODO: need to add limit param for a single page when nexus will support it
	var (
		err  error
		req  *http.Request
		resp *http.Response
	)

	u := hc.Url.JoinPath(RestAPIPrefix, "assets")
	q := u.Query()
	if nextToken != "" {
		q.Add("continuationToken", nextToken)
	}

	q.Add("repository", repository)
	u.RawQuery = q.Encode()

	req, err = http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("GetAssetsList build request err: %w", err)
	}

	tCtx, done := context.WithTimeout(ctx, hc.RequestTimeout)
	defer done()

	resp, err = hc.doRequest(tCtx, req)
	if err != nil {
		return nil, fmt.Errorf("GetAssetsList http req err: %w", err)
	}

	defer resp.Body.Close()
	return jsonTo[AssetList](resp.Body)
}

func (hc *HTTPClient) GetRepository(ctx context.Context, name string) (*Repository, error) {
	var (
		err  error
		req  *http.Request
		resp *http.Response
		addr = hc.Url.JoinPath(RestAPIPrefix, "repositories", name).String()
	)

	req, err = http.NewRequest(http.MethodGet, addr, nil)
	if err != nil {
		return nil, fmt.Errorf("GetRepository build request err: %w", err)
	}

	tCtx, done := context.WithTimeout(ctx, hc.RequestTimeout)
	defer done()

	resp, err = hc.doRequest(tCtx, req)
	if err != nil {
		return nil, fmt.Errorf("GetRepository http req err: %w", err)
	}

	defer resp.Body.Close()
	return jsonTo[Repository](resp.Body)
}

func (hc *HTTPClient) doRequest(ctx context.Context, req *http.Request) (*http.Response, error) {
	start := time.Now()
	resp, err := hc.Do(req.WithContext(ctx))
	if err != nil {
		hc.logger.DebugContext(ctx, "http request error",
			slog.String("method", req.Method),
			slog.String("addr", req.URL.String()),
			slog.String("err", err.Error()),
			slog.Duration("dur", time.Since(start)),
		)
		return nil, err
	}

	hc.logger.DebugContext(ctx, "http request end",
		slog.String("method", req.Method),
		slog.String("addr", req.URL.String()),
		slog.String("status", resp.Status),
		slog.Duration("dur", time.Since(start)),
	)

	if resp.StatusCode < http.StatusOK || resp.StatusCode > 299 {
		return nil, fmt.Errorf("%s %s END=%d", req.Method, req.RequestURI, resp.StatusCode)
	}

	return resp, nil
}

func NewHTTPClient(cfg *HTTPClientConfig) (*HTTPClient, error) {
	u, err := url.Parse(cfg.Address)
	if err != nil {
		return nil, err
	}

	client := HTTPClient{
		Url:    u,
		Client: &http.Client{},
	}

	if cfg.Timeout == 0 {
		client.RequestTimeout = DefaultClientTimeout
	} else {
		client.RequestTimeout = cfg.Timeout
	}

	if cfg.Logger != nil {
		client.logger = cfg.Logger
	} else {
		client.logger = slog.Default()
	}

	return &client, nil
}
