package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/MontFerret/worker/pkg/worker"
	"github.com/labstack/echo/v4"
)

// Server is HTTP server that wraps Ferret worker.
type Server struct {
	settings Settings
	worker   *worker.Worker
}

func New(settings Settings) (*Server, error) {
	worker, err := worker.New(worker.WithCustomCDP(settings.CDP))

	if err != nil {
		return nil, err
	}

	return &Server{
		settings,
		worker,
	}, nil
}

// Run start server that serve at the given port.
//
// Port should not begin with ":".
func (s *Server) Run(port uint64) error {
	router := echo.New()
	router.HideBanner = true

	router.POST("/", s.runScript)
	router.GET("/version", s.version)
	router.GET("/health", s.healthCheck)

	return router.Start(fmt.Sprintf("0.0.0.0:%d", port))
}

func (s *Server) version(ctx echo.Context) error {
	chromeVersionResp, err := http.Get(s.settings.CDP.VersionURL())

	if err != nil {
		ctx.Logger().Error("Failed to get response from Chrome", err)

		return ctx.NoContent(
			http.StatusFailedDependency,
		)
	}

	defer chromeVersionResp.Body.Close()

	chromeVersionBlob, err := ioutil.ReadAll(chromeVersionResp.Body)

	if err != nil {
		ctx.Logger().Error("Failed to read response from Chrome", err)

		return ctx.NoContent(
			http.StatusInternalServerError,
		)
	}

	chromeVersion := chromeVersionInternal{}

	err = json.Unmarshal(chromeVersionBlob, &chromeVersion)

	if err != nil {
		ctx.Logger().Error("Failed to parse response from Chrome", err)

		return ctx.NoContent(
			http.StatusInternalServerError,
		)
	}

	return ctx.JSON(
		http.StatusOK,
		Version{
			Worker: s.settings.Version,
			Ferret: s.settings.FerretVersion,
			Chrome: ChromeVersion{
				Browser:  chromeVersion.Browser,
				Protocol: chromeVersion.Protocol,
				V8:       chromeVersion.V8,
				WebKit:   chromeVersion.WebKit,
			},
		},
	)
}

func (s *Server) healthCheck(ctx echo.Context) error {
	out, err := http.Get(s.settings.CDP.VersionURL())

	if err != nil {
		ctx.Logger().Error("Failed to call Chrome", err)

		return ctx.NoContent(
			http.StatusFailedDependency,
		)
	}

	defer out.Body.Close()

	return ctx.NoContent(http.StatusOK)
}

func (s *Server) runScript(ctx echo.Context) error {
	reqctx := ctx.Request().Context()

	body := Script{}
	err := ctx.Bind(&body)

	if err != nil {
		ctx.Logger().Error("Failed to parse body", err)

		return ctx.JSON(
			http.StatusBadRequest,
			HttpError{err.Error()},
		)
	}

	out, err := s.worker.DoQuery(reqctx, body.Query)

	if err != nil {
		return ctx.JSON(
			http.StatusBadRequest,
			HttpError{err.Error()},
		)
	}

	return ctx.JSONBlob(http.StatusOK, out.Raw)
}
