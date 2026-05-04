package controllers

import (
	stderrors "errors"
	"strings"
	"testing"

	pkgerrors "github.com/pkg/errors"

	"github.com/MontFerret/ferret/v2/pkg/diagnostics"
	"github.com/MontFerret/ferret/v2/pkg/source"
)

func TestNewHTTPErrorFormatsWrappedDiagnostics(t *testing.T) {
	src := source.New("query.fql", "RETURN @missing")
	diag := &diagnostics.Diagnostic{
		Kind:    diagnostics.TypeError,
		Message: "missing query parameter",
		Source:  src,
		Spans: []diagnostics.ErrorSpan{
			diagnostics.NewMainErrorSpan(source.Span{Start: 7, End: 15}, "parameter is required"),
		},
	}
	diagSet := diagnostics.NewDiagnosticsOf[*diagnostics.Diagnostic]([]*diagnostics.Diagnostic{diag})

	got := newHTTPError(pkgerrors.Wrap(diagSet, "compile query"))

	if got.Error != "compile query: Found 1 errors" {
		t.Fatalf("unexpected error: %q", got.Error)
	}

	if got.Details != diagSet.Format() {
		t.Fatalf("expected formatted diagnostics, got:\n%s", got.Details)
	}

	for _, want := range []string{
		"TypeError: missing query parameter",
		"--> query.fql:1:8",
		"parameter is required",
	} {
		if !strings.Contains(got.Details, want) {
			t.Fatalf("expected details to contain %q, got:\n%s", want, got.Details)
		}
	}
}

func TestNewHTTPErrorFallsBackToPlainError(t *testing.T) {
	err := stderrors.New("invalid request body")

	got := newHTTPError(err)

	if got.Error != err.Error() {
		t.Fatalf("unexpected error: %q", got.Error)
	}

	if got.Details != err.Error() {
		t.Fatalf("unexpected details: %q", got.Details)
	}
}
