package repository

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"io"
	"log/slog"
	"net/http"
	"os"
)

func (r *Repository) Put(eCtx echo.Context) error {
	req := eCtx.Request()
	if req.ContentLength <= 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "empty body not allowed")
	}

	if uint64(req.ContentLength) > r.MaxFileSize {
		return echo.NewHTTPError(http.StatusBadRequest,
			fmt.Sprintf("max allowed file size %d, got %d", r.MaxFileSize, req.ContentLength),
		)
	}

	requestID := req.Header.Get(echo.HeaderXRequestID)
	if requestID == "" {
		requestID = eCtx.Response().Header().Get(echo.HeaderXRequestID)
	}

	filename := eCtx.Param("filename")
	ctx := context.WithValue(req.Context(), RequestIdContextKey, requestID)

	if !r.AllowRewrite {
		r.logger.Debug("rewrite is not allowed, checking file existence", slog.String("filename", filename))
		err := r.Storage.Exists(ctx, filename)
		switch {
		case err == nil:
			return echo.NewHTTPError(
				http.StatusConflict,
				fmt.Sprintf("file %s exists, rewrite is disabled", filename),
			)
		case os.IsNotExist(err):
		default:
			r.logger.ErrorContext(ctx, "storage.Exists() call err",
				slog.String("err", err.Error()),
				slog.String("filename", filename),
			)
		}
	}

	body := io.LimitReader(eCtx.Request().Body, req.ContentLength)
	defer eCtx.Request().Body.Close()

	if err := r.Storage.Put(ctx, filename, body); err != nil {
		r.logger.ErrorContext(ctx, "storage.Put() err",
			slog.String("err", err.Error()),
			slog.String("filename", filename),
		)

		return err
	}

	return eCtx.NoContent(http.StatusNoContent)
}
