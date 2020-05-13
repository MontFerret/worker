package server

import "github.com/MontFerret/worker/pkg/worker"

type (
	httpError struct {
		Error string `json:"error"`
	}

	runScriptBody struct {
		worker.Query
	}
)
