package worker

import (
	"context"

	"github.com/MontFerret/ferret/pkg/compiler"
	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/MontFerret/ferret/pkg/drivers/cdp"
	"github.com/MontFerret/ferret/pkg/drivers/http"
	"github.com/MontFerret/ferret/pkg/runtime"
	"github.com/pkg/errors"

	"github.com/MontFerret/worker/pkg/caching"
)

// Worker accepts FQL-script, run it and return result.
type Worker struct {
	comp    *compiler.Compiler
	drivers []drivers.Driver
	cache   caching.Cache
}

// New returns Worker without file system access.
func New(setters ...Option) (*Worker, error) {
	opts := newOptions()

	for _, setter := range setters {
		setter(opts)
	}

	var comp *compiler.Compiler

	if opts.noStdlib {
		comp = compiler.New(compiler.WithoutStdlib())
	} else {
		comp = compiler.New()
	}

	for _, functions := range opts.functions {
		if err := comp.RegisterFunctions(&functions); err != nil {
			return nil, err
		}
	}

	return &Worker{
		comp: comp,
		drivers: []drivers.Driver{
			cdp.NewDriver(
				cdp.WithAddress(opts.cdp.BaseURL()),
			),
			http.NewDriver(),
		},
		cache: opts.cache,
	}, nil
}

func (w *Worker) DoQuery(ctx context.Context, query Query) (Result, error) {
	if query.Text == "" {
		return Result{}, errors.New("missed query text")
	}

	program, err := w.compiledOrCached(query.Text)

	if err != nil {
		return Result{}, errors.Wrap(err, "compile query")
	}

	for _, d := range w.drivers {
		ctx = drivers.WithContext(ctx, d)
	}

	r, err := program.Run(ctx, runtime.WithParams(query.Params))

	if err != nil {
		return Result{}, errors.Wrap(err, "run program")
	}

	return Result{
		Raw: r,
	}, nil
}

func (w *Worker) compiledOrCached(text string) (*runtime.Program, error) {
	var program *runtime.Program

	if w.cache != nil {
		found, isFound := w.cache.Get(text)

		if isFound && found != nil {
			program = found.(*runtime.Program)
		}
	}

	if program == nil {
		compiled, err := w.comp.Compile(text)

		if err != nil {
			return nil, errors.Wrap(err, "compile")
		}

		program = compiled
	}

	if w.cache != nil {
		w.cache.Set(text, program)
	}

	return program, nil
}
