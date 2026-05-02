package controllers

import (
	"errors"

	"github.com/MontFerret/ferret/v2/pkg/diagnostics"
)

type (
	HTTPError struct {
		Error   string `json:"error"`
		Details string `json:"details"`
	}
)

func newHTTPError(err error) HTTPError {
	if err == nil {
		return HTTPError{}
	}

	details := err.Error()

	var formattable diagnostics.Formattable
	if errors.As(err, &formattable) {
		details = formattable.Format()
	}

	return HTTPError{
		Error:   err.Error(),
		Details: details,
	}
}
