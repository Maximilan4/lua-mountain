package nexus

import (
	"context"
	"fmt"
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
		Url    *url.URL
		logger *slog.Logger
	}
)

func (hc *HTTPClient) GetAssetsList(ctx context.Context, repository string, nextToken string) (*AssetList, error) {
	// TODO: need to add limit param for a single page when nexus will support it
	var (
		err  error
		req  *http.Request
		resp *http.Response
	)

	u := hc.Url.JoinPath("assets")
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
	resp, err = hc.doRequest(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("GetAssetsList http req err: %w", err)
	}

	defer resp.Body.Close()
	return convertResponse[AssetList](resp.Body)
}

func (hc *HTTPClient) GetRepository(ctx context.Context, name string) (*Repository, error) {
	var (
		err  error
		req  *http.Request
		resp *http.Response
		addr = hc.Url.JoinPath("repositories", name).String()
	)

	req, err = http.NewRequest(http.MethodGet, addr, nil)
	if err != nil {
		return nil, fmt.Errorf("GetRepository build request err: %w", err)
	}

	resp, err = hc.doRequest(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("GetRepository http req err: %w", err)
	}

	defer resp.Body.Close()
	return convertResponse[Repository](resp.Body)
}

func (hc *HTTPClient) doRequest(ctx context.Context, req *http.Request) (resp *http.Response, err error) {
	tCtx, done := context.WithTimeout(ctx, hc.Timeout)
	defer done()
	hc.logger.DebugContext(tCtx, "http request start",
		slog.String("method", http.MethodGet),
		slog.String("addr", req.RequestURI),
	)

	resp, err = hc.Do(req)
	if err != nil {
		hc.logger.DebugContext(ctx, "request error",
			slog.String("method", http.MethodGet),
			slog.String("addr", req.RequestURI),
			slog.String("err", err.Error()),
		)
		return nil, err
	}

	hc.logger.DebugContext(ctx, "http request end",
		slog.String("method", http.MethodGet),
		slog.String("addr", req.RequestURI),
		slog.String("status", resp.Status),
	)

	if resp.StatusCode != http.StatusOK {
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
		Url:    u.JoinPath(RestAPIPrefix),
		Client: &http.Client{},
	}

	if cfg.Timeout == 0 {
		client.Timeout = DefaultClientTimeout
	} else {
		client.Timeout = cfg.Timeout
	}

	if cfg.Logger != nil {
		client.logger = cfg.Logger
	} else {
		client.logger = slog.Default()
	}

	return &client, nil
}
