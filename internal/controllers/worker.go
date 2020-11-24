package controllers

import (
	"github.com/MontFerret/worker/pkg/worker"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"net/http"
)

type (
	ScriptDto struct {
		worker.Query
	}

	Worker struct {
		worker *worker.Worker
	}
)

func NewWorker(settings worker.CDPSettings) (*Worker, error) {
	w, err := worker.New(worker.WithCustomCDP(settings))

	if err != nil {
		return nil, errors.Wrap(err, "create a worker instance")
	}

	return &Worker{w}, nil
}

func (c *Worker) Use(e *echo.Echo) {
	e.POST("/", c.runScript)
}

func (c *Worker) runScript(ctx echo.Context) error {
	body := ScriptDto{}
	err := ctx.Bind(&body)

	if err != nil {
		ctx.Logger().Error("Failed to parse body", err)

		return ctx.JSON(
			http.StatusBadRequest,
			HttpError{err.Error()},
		)
	}

	out, err := c.worker.DoQuery(ctx.Request().Context(), body.Query)

	if err != nil {
		return ctx.JSON(
			http.StatusBadRequest,
			HttpError{err.Error()},
		)
	}

	return ctx.JSONBlob(http.StatusOK, out.Raw)
}
