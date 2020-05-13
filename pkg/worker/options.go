package worker

import (
	"fmt"
	"github.com/MontFerret/ferret/pkg/runtime/core"
)

type (
	CDPSettings struct {
		Host string `json:"host"`
		Port uint64 `json:"port"`
	}

	Options struct {
		functions []core.Functions
		noStdlib  bool
		cdp       CDPSettings
	}

	Option func(opts *Options)
)

func (s CDPSettings) URL() string {
	return fmt.Sprintf("http://%s:%d", s.Host, s.Port)
}

func newOptions() *Options {
	return &Options{
		functions: make([]core.Functions, 0, 5),
		cdp: CDPSettings{
			Host: "127.0.0.1",
			Port: 9222,
		},
	}
}

func WithFunctions(functions core.Functions) Option {
	return func(opts *Options) {
		opts.functions = append(opts.functions, functions)
	}
}

func WithoutStdlib() Option {
	return func(opts *Options) {
		opts.noStdlib = true
	}
}

func WithCustomCDP(cdp CDPSettings) Option {
	return func(opts *Options) {
		opts.cdp = cdp
	}
}
