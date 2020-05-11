package worker

import (
	"context"

	"github.com/MontFerret/ferret/pkg/compiler"
	"github.com/MontFerret/ferret/pkg/runtime"
	"github.com/pkg/errors"
)

// Worker accepts FQL-script, run it and return result.
type Worker struct {
	comp *compiler.Compiler
}

// DefaultWorker is the default Worker.
var DefaultWorker = Worker{
	comp: compiler.New(),
}

type (
	// Query is the FQL-script.
	Query struct {
		Text   string                 `json:"text"`
		Params map[string]interface{} `json:"params"`
	}

	// Result is the result of Query.
	Result struct {
		Raw []byte
	}
)

func (w *Worker) DoQuery(ctx context.Context, query Query) (Result, error) {
	program, err := w.comp.Compile(query.Text)
	if err != nil {
		return Result{}, errors.Wrap(err, "compile query")
	}

	r, err := program.Run(ctx, runtime.WithParams(query.Params))
	if err != nil {
		return Result{}, errors.Wrap(err, "run program")
	}

	return Result{
		Raw: r,
	}, nil
}
