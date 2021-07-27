package controllers

import (
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"

	"github.com/MontFerret/worker/pkg/worker"
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
		var rawBodyContent string
		rawBody, rawErr := io.ReadAll(ctx.Request().Body)

		if rawErr == nil {
			rawBodyContent = string(rawBody)
		}

		ctx.Logger().Errorf("Failed to parse body: %s: %s", rawBodyContent, err)

		return ctx.JSON(
			http.StatusBadRequest,
			HttpError{err.Error()},
		)
	}

	ctx.Logger().Debugf("Received query to execute: %s", body.Query.Text)

	out, err := c.worker.DoQuery(ctx.Request().Context(), body.Query)

	if err != nil {
		ctx.Logger().Errorf("Failed to execute query: %s: %s", body.Query.Text, err)

		return ctx.JSON(
			http.StatusBadRequest,
			HttpError{err.Error()},
		)
	}

	return ctx.JSONBlob(http.StatusOK, out.Raw)
}
