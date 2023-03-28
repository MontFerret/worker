package controllers

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/MontFerret/worker/pkg/worker"
)

type (
	Health struct {
		settings worker.CDPSettings
	}
)

func NewHealth(settings worker.CDPSettings) (*Health, error) {
	return &Health{settings}, nil
}

func (c *Health) Use(e *echo.Echo) {
	e.GET("/health", c.healthCheck)
}

func (c *Health) healthCheck(ctx echo.Context) error {
	out, err := http.Get(c.settings.VersionURL())

	if err != nil {
		ctx.Logger().Error("Failed to call Chrome", err)

		return ctx.NoContent(
			http.StatusFailedDependency,
		)
	}

	defer out.Body.Close()

	return ctx.NoContent(http.StatusOK)
}
