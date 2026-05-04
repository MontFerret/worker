package worker

import (
	"context"

	"github.com/pkg/errors"

	"github.com/MontFerret/worker/internal/storage"

	"github.com/MontFerret/ferret/v2"
	"github.com/MontFerret/ferret/v2/pkg/source"

	"github.com/MontFerret/worker/pkg/caching"
)

// Worker accepts FQL-script, run it and return result.
type Worker struct {
	engine *ferret.Engine
	cache  caching.Cache[*ferret.Plan]
}

// New returns Worker without file system access.
func New(setters ...Option) (*Worker, error) {
	opts, err := newOptions(setters)

	if err != nil {
		return nil, errors.Wrap(err, "create options")
	}

	engine, err := ferret.New(opts.Engine...)

	if err != nil {
		return nil, errors.Wrap(err, "create engine")
	}

	cache, err := storage.NewCache(opts.Cache...)

	if err != nil {
		return nil, errors.Wrap(err, "create cache storage")
	}

	return &Worker{
		engine: engine,
		cache:  cache,
	}, nil
}

func (w *Worker) DoQuery(ctx context.Context, query Query) (Result, error) {
	if query.Text == "" {
		return Result{}, errors.New("missed query text")
	}

	plan, err := w.compiledOrCached(ctx, query.Text)

	if err != nil {
		return Result{}, errors.Wrap(err, "compile query")
	}

	session, err := plan.NewSession(ctx, ferret.WithSessionParams(query.Params))

	if err != nil {
		return Result{}, errors.Wrap(err, "create session")
	}

	defer session.Close()

	out, err := session.Run(ctx)

	if err != nil {
		return Result{}, errors.Wrap(err, "run program")
	}

	return Result{
		ContentType: out.ContentType,
		Raw:         out.Content,
	}, nil
}

func (w *Worker) compiledOrCached(ctx context.Context, text string) (*ferret.Plan, error) {
	var plan *ferret.Plan

	if w.cache != nil {
		found, isFound := w.cache.Get(text)

		if isFound && found != nil {
			plan = found
		}
	}

	if plan == nil {
		compiled, err := w.engine.Compile(ctx, source.NewAnonymous(text))

		if err != nil {
			return nil, errors.Wrap(err, "compile")
		}

		plan = compiled
	}

	if w.cache != nil {
		w.cache.Set(text, plan)
	}

	return plan, nil
}
