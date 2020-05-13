package server

import (
	"net/http"

	"github.com/MontFerret/worker/pkg/worker"
	"github.com/labstack/echo"
)

// Server is HTTP server that wraps Ferret worker.
type Server struct {
	worker *worker.Worker
}

func New() *Server {
	worker, _ := worker.New()

	return &Server{
		worker,
	}
}

// Run start server that serve at the given port.
//
// Port should not begin with ":".
func (s *Server) Run(port string) error {
	router := echo.New()
	router.HideBanner = true

	router.POST("/", s.runScript)

	return router.Start("0.0.0.0:" + port)
}

type httpError struct {
	Error string `json:"error"`
}

type runScriptBody struct {
	worker.Query
}

func (s *Server) runScript(ctx echo.Context) error {
	reqctx := ctx.Request().Context()

	body := runScriptBody{}
	err := ctx.Bind(&body)

	if err != nil {
		return ctx.JSONPretty(
			http.StatusBadRequest,
			httpError{err.Error()},
			"  ",
		)
	}

	out, err := s.worker.DoQuery(reqctx, body.Query)

	if err != nil {
		return ctx.JSONPretty(
			http.StatusBadRequest,
			httpError{err.Error()},
			"  ",
		)
	}

	return ctx.JSONPretty(http.StatusOK, string(out.Raw), "  ")
}
