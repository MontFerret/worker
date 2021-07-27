package controllers

import (
	"io"
	"net/http"

	"github.com/labstack/echo/v4"

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

func NewWorker(worker *worker.Worker) (*Worker, error) {
	return &Worker{worker}, nil
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
