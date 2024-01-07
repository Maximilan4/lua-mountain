package mw

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"strings"
)

func AllowedExtensions(extensions []string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			filename := c.Param("filename")
			if !IsAllowedExtension(filename, extensions) {
				return echo.NewHTTPError(
					http.StatusBadRequest,
					fmt.Sprintf("filename %s has not allowed extension, allowed are: %v", filename, extensions),
				)
			}

			return next(c)
		}
	}
}

func IsAllowedExtension(filename string, allowed []string) bool {
	// TODO: may be use regular expression
	for _, ext := range allowed {
		if strings.HasSuffix(filename, ext) {
			return true
		}
	}

	return false
}