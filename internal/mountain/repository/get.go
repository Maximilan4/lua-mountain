package repository

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"log/slog"
	"net/http"
)

func (r *Repository) Get(eCtx echo.Context) error {
	requestID := eCtx.Request().Header.Get(echo.HeaderXRequestID)
	if requestID == "" {
		requestID = eCtx.Response().Header().Get(echo.HeaderXRequestID)
	}

	ctx := context.WithValue(eCtx.Request().Context(), RequestIdContextKey, requestID)

	filename := eCtx.Param("filename")
	if err := r.Storage.Exists(ctx, filename); err != nil {
		r.logger.ErrorContext(ctx, "storage.Exists() call err",
			slog.String("err", err.Error()),
			slog.String("filename", filename),
		)
		return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("file %s not found", filename))
	}

	f, err := r.Storage.Get(ctx, filename)
	if err != nil {
		r.logger.ErrorContext(ctx, "storage.Get() call err",
			slog.String("err", err.Error()),
			slog.String("filename", filename),
		)
		return err
	}

	defer f.Close()

	return eCtx.Stream(http.StatusOK, echo.MIMETextPlain, f)
}