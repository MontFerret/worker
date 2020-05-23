package server

import "github.com/MontFerret/worker/pkg/worker"

type Settings struct {
	Version       string
	FerretVersion string
	CDP           worker.CDPSettings
}
