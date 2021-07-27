package worker

import (
	"fmt"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/worker/pkg/caching"
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
		cache     caching.Cache
	}

	Option func(opts *Options)
)

func (s CDPSettings) BaseURL() string {
	return fmt.Sprintf("http://%s:%d", s.Host, s.Port)
}

func (s CDPSettings) VersionURL() string {
	return fmt.Sprintf("%s/json/version", s.BaseURL())
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

func WithCache(cache caching.Cache) Option {
	return func(opts *Options) {
		opts.cache = cache
	}
}
