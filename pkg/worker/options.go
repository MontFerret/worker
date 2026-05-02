package worker

import (
	"fmt"

	"github.com/MontFerret/ferret/v2"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/worker/pkg/caching"
)

type (
	CDPSettings struct {
		Host     string `json:"host"`
		Port     uint64 `json:"port"`
		Disabled bool   `json:"disabled"`
	}

	Options struct {
		engine []ferret.Option
		cdp    CDPSettings
		cache  caching.Cache[*ferret.Plan]
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
		engine: []ferret.Option{},
		cdp: CDPSettings{
			Host: "127.0.0.1",
			Port: 9222,
		},
	}
}

func WithFunctions(functions *runtime.Functions) Option {
	return func(opts *Options) {
		opts.engine = append(opts.engine, ferret.WithFunctions(functions))
	}
}

func WithoutStdlib() Option {
	return func(opts *Options) {
		opts.engine = append(opts.engine, ferret.WithoutStdlib())
	}
}

func WithCustomCDP(cdp CDPSettings) Option {
	return func(opts *Options) {
		opts.cdp = cdp
	}
}

func WithCache(cache caching.Cache[*ferret.Plan]) Option {
	return func(opts *Options) {
		opts.cache = cache
	}
}
