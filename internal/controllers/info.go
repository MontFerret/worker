package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"

	"github.com/MontFerret/worker/pkg/worker"
)

type (
	chromeVersionInternal struct {
		Browser  string `json:"Browser"`
		Protocol string `json:"Protocol-Version"`
		V8       string `json:"V8-Version"`
		WebKit   string `json:"WebKit-Version"`
	}

	ChromeVersionDto struct {
		Browser  string `json:"browser"`
		Protocol string `json:"protocol"`
		V8       string `json:"v8"`
		WebKit   string `json:"webkit"`
	}

	VersionDto struct {
		Worker string           `json:"worker"`
		Chrome ChromeVersionDto `json:"chrome"`
		Ferret string           `json:"ferret"`
	}

	InfoDto struct {
		IP      string     `json:"ip"`
		Version VersionDto `json:"version"`
	}

	InfoSettings struct {
		Version       string
		FerretVersion string
		CDP           worker.CDPSettings
	}

	Info struct {
		settings InfoSettings
	}
)

func NewInfo(settings InfoSettings) (*Info, error) {
	return &Info{settings}, nil
}

func (c *Info) Use(e *echo.Echo) {
	e.GET("/info", func(ctx echo.Context) error {
		version, err := c.version(ctx.Request().Context())

		if err != nil {
			ctx.Logger().Error("failed to retrieve version", err)

			return ctx.NoContent(
				http.StatusFailedDependency,
			)
		}

		ip, err := c.ip(ctx.Request().Context())

		if err != nil {
			ctx.Logger().Error("failed to retrieve ip address", err)

			return ctx.NoContent(
				http.StatusFailedDependency,
			)
		}

		return ctx.JSON(
			http.StatusOK,
			InfoDto{
				IP:      ip,
				Version: version,
			},
		)
	})
}

func (c *Info) version(_ context.Context) (VersionDto, error) {
	var chromeVersion ChromeVersionDto

	if !c.settings.CDP.Disabled {
		chromeVersionResp, err := http.Get(c.settings.CDP.VersionURL())

		if err != nil {
			return VersionDto{}, errors.Wrap(err, "call Chrome")
		}

		defer chromeVersionResp.Body.Close()

		chromeVersionBlob, err := io.ReadAll(chromeVersionResp.Body)

		if err != nil {
			return VersionDto{}, errors.Wrap(err, "read response from Chrome")
		}

		chromeVersionInternal := chromeVersionInternal{}

		err = json.Unmarshal(chromeVersionBlob, &chromeVersionInternal)

		if err != nil {
			return VersionDto{}, errors.Wrap(err, "parse response from Chrome")
		}

		chromeVersion = ChromeVersionDto(chromeVersionInternal)
	}

	return VersionDto{
		Worker: c.settings.Version,
		Ferret: c.settings.FerretVersion,
		Chrome: chromeVersion,
	}, nil
}

func (c *Info) ip(_ context.Context) (string, error) {
	rsp, err := http.Get("http://checkip.amazonaws.com")

	if err != nil {
		return "", errors.Wrap(err, "call service")
	}

	defer rsp.Body.Close()

	buf, err := io.ReadAll(rsp.Body)

	if err != nil {
		return "", errors.Wrap(err, "parse response")
	}

	return string(bytes.TrimSpace(buf)), nil
}
