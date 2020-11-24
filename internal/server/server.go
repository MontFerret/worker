package server

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// Server is HTTP server that wraps Ferret worker.
type Server struct {
	router *echo.Echo
}

func New() (*Server, error) {
	router := echo.New()
	router.HideBanner = true

	router.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		Skipper:      middleware.DefaultSkipper,
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))
	router.Use(middleware.BodyLimit("1M"))
	router.Use(middleware.RequestID())
	router.Use(middleware.Logger())
	router.Use(middleware.Recover())
	router.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: 5,
	}))

	return &Server{router}, nil
}

func (s *Server) Router() *echo.Echo {
	return s.router
}

// Run start server that serve at the given port.
//
// Port should not begin with ":".
func (s *Server) Run(port uint64) error {
	return s.router.Start(fmt.Sprintf("0.0.0.0:%d", port))
}
