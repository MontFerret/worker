package server

import "github.com/labstack/echo/v4"

type Route interface {
	Use(e *echo.Echo)
}
