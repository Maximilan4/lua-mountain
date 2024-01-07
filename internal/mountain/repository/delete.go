package repository

import (
	"context"
	"github.com/labstack/echo/v4"
	"log/slog"
	"net/http"
	"os"
)

func (r *Repository) Delete(eCtx echo.Context) error {
	req := eCtx.Request()
	if req.ContentLength != 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "body is forbidden for DELETE request")
	}

	requestID := req.Header.Get(echo.HeaderXRequestID)
	if requestID == "" {
		requestID = eCtx.Response().Header().Get(echo.HeaderXRequestID)
	}

	filename := eCtx.Param("filename")
	ctx := context.WithValue(req.Context(), RequestIdContextKey, requestID)

	if err := r.Storage.Delete(ctx, filename); err != nil && !os.IsNotExist(err) {
		r.logger.ErrorContext(ctx, "storage.Delete() err",
			slog.String("err", err.Error()),
			slog.String("filename", filename),
		)

		return err
	}

	return eCtx.NoContent(http.StatusNoContent)
}
