package server

import (
	"fmt"
	"net/http"
	"time"

	"golang.org/x/time/rate"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/ziflex/lecho/v3"
)

type (
	// Options is a set of options for server.
	Options struct {
		// RequestLimit is a number of requests per second for each IP.
		// If value is 0, rate limit is disabled.
		RequestLimit uint64

		// RequestLimitTimeWindow is a period of requests limit in seconds.
		// If value is 0, rate limit is set default value.
		RequestLimitTimeWindow uint64

		// BodyLimit is a maximum size of request body.
		// If value is 0, body limit is disabled.
		BodyLimit uint64
	}

	// Server is HTTP server that wraps Ferret worker.
	Server struct {
		router *echo.Echo
	}
)

func New(logger *lecho.Logger, opts Options) (*Server, error) {
	router := echo.New()
	router.Logger = logger
	router.HideBanner = true

	if opts.RequestLimit > 0 {
		var dur time.Duration

		if opts.RequestLimitTimeWindow > 0 {
			dur = time.Second * time.Duration(opts.RequestLimitTimeWindow)
		}

		router.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStoreWithConfig(middleware.RateLimiterMemoryStoreConfig{
			Rate:      rate.Limit(opts.RequestLimit),
			ExpiresIn: dur,
		})))
	}

	router.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		Skipper:      middleware.DefaultSkipper,
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	if opts.BodyLimit > 0 {
		router.Use(middleware.BodyLimit(fmt.Sprintf("%d", opts.BodyLimit)))
	}

	router.Use(middleware.RequestID())
	router.Use(lecho.Middleware(lecho.Config{
		Logger: logger,
	}))
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
