package repository

import (
	"archive/zip"
	"context"
	"github.com/labstack/echo/v4"
	"io"
	"log/slog"
	"lua-mountain/internal/mountain/luarocks"
	"net/http"
	"path/filepath"
	"regexp"
	"strings"
)


var (
	semverRegexp = regexp.MustCompile(
	"(0|[1-9]\\d*)\\.(0|[1-9]\\d*)\\.(0|[1-9]\\d*)(?:-((?:0|[1-9]\\d*|\\d*[a-zA-Z-][0-9a-zA-Z-]*)" +
		"(?:\\.(?:0|[1-9]\\d*|\\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\\+([0-9a-zA-Z-]+(?:\\.[0-9a-zA-Z-]+)*))?$",
	 )
)


func (r *Repository) GetManifest(eCtx echo.Context) error {
	req := eCtx.Request()
	resp := eCtx.Response()
	requestID := req.Header.Get(echo.HeaderXRequestID)
	if requestID == "" {
		requestID = resp.Header().Get(echo.HeaderXRequestID)
	}

	ctx := context.WithValue(req.Context(), RequestIdContextKey, requestID)

	list, err := r.Storage.List(ctx)
	if err != nil {
		r.logger.ErrorContext(ctx, "storage.List() err", slog.String("err", err.Error()))
		return err
	}

	resp.Header().Add("Content-Type", "text/x-lua")
	writer := luarocks.NewWriter(resp)
	if err = writer.WriteRepositoryPackages(r.getRocksList(ctx, list)); err != nil {
		return err
	}

	return nil
}

func (r *Repository) GetManifestJson(eCtx echo.Context) error {
	req := eCtx.Request()
	resp := eCtx.Response()
	requestID := req.Header.Get(echo.HeaderXRequestID)
	if requestID == "" {
		requestID = resp.Header().Get(echo.HeaderXRequestID)
	}

	ctx := context.WithValue(req.Context(), RequestIdContextKey, requestID)

	list, err := r.Storage.List(ctx)
	if err != nil {
		r.logger.ErrorContext(ctx, "storage.List() err", slog.String("err", err.Error()))
		return err
	}

	return eCtx.JSON(http.StatusOK, r.getRocksList(ctx, list))
}


func (r *Repository) GetManifestZip(eCtx echo.Context) error {
	req := eCtx.Request()
	resp := eCtx.Response()
	requestID := req.Header.Get(echo.HeaderXRequestID)
	if requestID == "" {
		requestID = resp.Header().Get(echo.HeaderXRequestID)
	}

	ctx := context.WithValue(req.Context(), RequestIdContextKey, requestID)
	list, err := r.Storage.List(ctx)
	if err != nil {
		r.logger.ErrorContext(ctx, "storage.List() err", slog.String("err", err.Error()))
		return err
	}

	var (
		filename, _ = strings.CutSuffix(filepath.Base(req.URL.Path), filepath.Ext(req.URL.Path))
		archive = zip.NewWriter(resp)
		f io.Writer
	)

	defer archive.Close()
	f, err = archive.Create(filename)
	if err != nil {
		r.logger.ErrorContext(ctx, "unable to create archive",
			slog.String("err", err.Error()),
			slog.String("filename", filename),
		)

		return err
	}
	writer := luarocks.NewWriter(f)
	if err = writer.WriteRepositoryPackages(r.getRocksList(ctx, list)); err != nil {
		return err
	}

	return nil
}


func (r *Repository) getRocksList(ctx context.Context, list []string) luarocks.RocksList {
	var (
		rocks = make(luarocks.RocksList, 0, len(list))
		fileName, rockName, version, arch string
		found bool
		index int
	)

	for _, fileName = range list {
		if fileName, found = strings.CutSuffix(fileName, ".rockspec"); found {
			arch = "rockspec"
		} else if fileName, found = strings.CutSuffix(fileName, ".rock"); found {
			arch = fileName[strings.LastIndexByte(fileName, '.') + 1:]
			fileName, _ = strings.CutSuffix(fileName, fileName[strings.LastIndexByte(fileName, '.'):])
		} else {
			r.logger.DebugContext(ctx, "unable to define file arch", slog.String("filename", fileName))
			continue
		}

		if index = strings.Index(fileName, "scm"); index > 0 {
			version = fileName[index:]
			rockName = fileName[:index-1]
		} else {
			match := semverRegexp.FindAllString(fileName, 1)
			if len(match) == 0 {
				r.logger.DebugContext(ctx, "unable to parse version", slog.String("filename", fileName))
				continue
			}

			version = match[0]
			rockName, _ = strings.CutSuffix(fileName, version)
		}
		rocks.Add(strings.TrimRight(rockName, "-."), version, arch)
	}

	return rocks
}