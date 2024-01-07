package server

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	slogecho "github.com/samber/slog-echo"
	"lua-mountain/internal/mountain/logging"
)

type (
	Config struct {
		Address string `yaml:"address"`
		Port string `yaml:"port"`
	}
)

func Init() *echo.Echo {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	e.Use(
		slogecho.New(logging.DefaultLogger),
		middleware.RequestID(),
		middleware.Recover(),
	)

	return e
}
