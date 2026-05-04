package controllers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"

	workerpkg "github.com/MontFerret/worker/pkg/worker"
)

func TestWorkerRunScriptReturnsDiagnosticDetails(t *testing.T) {
	wkr, err := workerpkg.New()
	if err != nil {
		t.Fatalf("new worker: %v", err)
	}

	controller, err := NewWorker(wkr)
	if err != nil {
		t.Fatalf("new worker controller: %v", err)
	}

	e := echo.New()
	controller.Use(e)

	req := httptest.NewRequest(
		http.MethodPost,
		"/",
		strings.NewReader(`{"text":"RETURN { a: @foo, b: @bar }"}`),
	)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d: %s", http.StatusBadRequest, rec.Code, rec.Body.String())
	}

	var payload HTTPError
	if err := json.Unmarshal(rec.Body.Bytes(), &payload); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if !strings.Contains(payload.Error, "run program: Found") {
		t.Fatalf("expected wrapped execution error, got %q", payload.Error)
	}

	if payload.Details == "" {
		t.Fatal("expected details to be populated")
	}

	if payload.Details == payload.Error {
		t.Fatalf("expected formatted diagnostics, got only plain error: %q", payload.Details)
	}

	for _, want := range []string{
		"anonymous",
		"RETURN { a: @foo, b: @bar }",
		"@foo",
		"@bar",
	} {
		if !strings.Contains(payload.Details, want) {
			t.Fatalf("expected details to contain %q, got:\n%s", want, payload.Details)
		}
	}
}
