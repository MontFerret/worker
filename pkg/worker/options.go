package worker

import (
	"fmt"

	"github.com/MontFerret/contrib/modules/csv"
	"github.com/MontFerret/contrib/modules/db/sqlite"
	"github.com/MontFerret/contrib/modules/document/pdf"
	"github.com/MontFerret/contrib/modules/document/xlsx"
	"github.com/MontFerret/contrib/modules/net/rest"
	"github.com/MontFerret/contrib/modules/security/jwt"
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
	"github.com/MontFerret/ferret/v2/pkg/module"
	"github.com/MontFerret/worker/pkg/caching"
)

type (
	CDPSettings struct {
		Host     string `json:"host"`
		Port     uint64 `json:"port"`
		Disabled bool   `json:"disabled"`
	}

	options struct {
		engine     []ferret.Option
		cdp        CDPSettings
		cache      []caching.Option
		httpPolicy HTTPPolicy
		rest       bool
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
		httpPolicy: DefaultHTTPPolicy(),
	}

	for _, setter := range setters {
		setter(opts)
	}

	if err := validateHTTPPolicy(opts.httpPolicy); err != nil {
		return Options{}, err
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

	mods := []module.Module{
		csv.New(),
		toml.New(),
		xml.New(),
		yaml.New(),
		article.New(),
		robots.New(),
		sitemap.New(),
		drivers,
		pdf.New(),
		xlsx.New(),
		sqlite.New(sqlite.WithMemoryOnly()),
		jwt.New(),
	}

	if opts.rest {
		mods = append(mods, rest.New())
	}

	def := []ferret.Option{
		ferret.WithModules(mods...),
		ferret.WithNetwork(newNetwork(opts.httpPolicy)),
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

// WithHTTPPolicy configures Ferret policy-backed HTTP egress for Worker.
func WithHTTPPolicy(policy HTTPPolicy) Option {
	return func(opts *options) {
		opts.httpPolicy = policy
	}
}

// WithRESTModule enables the NET::REST module. Queries can use it to make
// outbound HTTP requests through the configured Ferret HTTP policy.
func WithRESTModule() Option {
	return func(opts *options) {
		opts.rest = true
	}
}
