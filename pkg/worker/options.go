package worker

import (
	"fmt"

	"github.com/MontFerret/contrib/modules/csv"
	"github.com/MontFerret/contrib/modules/toml"
	"github.com/MontFerret/contrib/modules/web/article"
	"github.com/MontFerret/contrib/modules/web/html"
	"github.com/MontFerret/contrib/modules/web/html/drivers/cdp"
	"github.com/MontFerret/contrib/modules/web/html/drivers/memory"
	"github.com/MontFerret/contrib/modules/web/robots"
	"github.com/MontFerret/contrib/modules/web/sitemap"
	"github.com/MontFerret/contrib/modules/xml"
	"github.com/MontFerret/contrib/modules/yaml"
	"github.com/MontFerret/ferret/v2"
	"github.com/MontFerret/worker/pkg/caching"
)

type (
	CDPSettings struct {
		Host     string `json:"host"`
		Port     uint64 `json:"port"`
		Disabled bool   `json:"disabled"`
	}

	options struct {
		engine []ferret.Option
		cdp    CDPSettings
		cache  []caching.Option
	}

	Options struct {
		Engine []ferret.Option
		Cache  []caching.Option
	}

	Option func(opts *options)
)

func (s CDPSettings) BaseURL() string {
	return fmt.Sprintf("http://%s:%d", s.Host, s.Port)
}

func (s CDPSettings) VersionURL() string {
	return fmt.Sprintf("%s/json/version", s.BaseURL())
}

func newOptions(setters []Option) (Options, error) {
	opts := &options{
		cdp: CDPSettings{
			Host: "127.0.0.1",
			Port: 9222,
		},
		cache: []caching.Option{
			caching.WithSize(50),
		},
	}

	for _, setter := range setters {
		setter(opts)
	}

	drivers, err := html.New(
		html.WithDefaultDriver(memory.New()),
		html.WithDrivers(cdp.New(
			cdp.WithAddress(opts.cdp.BaseURL()),
		)),
	)

	if err != nil {
		return Options{}, fmt.Errorf("create HTML module: %w", err)
	}

	def := []ferret.Option{
		ferret.WithModules(
			csv.New(),
			toml.New(),
			xml.New(),
			yaml.New(),
			article.New(),
			robots.New(),
			sitemap.New(),
			drivers,
		),
	}

	if len(opts.engine) > 0 {
		def = append(def, opts.engine...)
	}

	return Options{
		Engine: def,
		Cache:  opts.cache,
	}, nil
}

func WithEngineOptions(engineOpts ...ferret.Option) Option {
	return func(opts *options) {
		opts.engine = append(opts.engine, engineOpts...)
	}
}

func WithCustomCDP(cdp CDPSettings) Option {
	return func(opts *options) {
		opts.cdp = cdp
	}
}

func WithCacheSize(cache uint) Option {
	return func(opts *options) {
		opts.cache = append(opts.cache, caching.WithSize(cache))
	}
}

func WithFSRoot(rootDir string) Option {
	return func(opts *options) {
		opts.engine = append(opts.engine, ferret.WithFSRoot(rootDir))
	}
}
