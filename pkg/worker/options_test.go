package worker

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync/atomic"
	"testing"
)

func TestRESTModuleDisabledByDefault(t *testing.T) {
	wkr, err := New()
	if err != nil {
		t.Fatalf("new worker: %v", err)
	}

	_, err = wkr.DoQuery(context.Background(), Query{
		Text: `LET api = NET::REST::CLIENT("http://example.test") RETURN api`,
	})
	if err == nil {
		t.Fatal("expected NET::REST to be unavailable by default")
	}

	if !strings.Contains(err.Error(), "unresolved function") {
		t.Fatalf("expected unresolved function error, got %v", err)
	}
}

func TestWithRESTModuleEnablesRESTClient(t *testing.T) {
	var called atomic.Bool
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called.Store(true)

		if r.URL.Path != "/health" {
			t.Fatalf("expected /health path, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{"ok": true})
	}))
	defer server.Close()

	wkr, err := New(
		WithRESTModule(),
		WithHTTPPolicy(testHTTPPolicy(t, server.URL)),
	)
	if err != nil {
		t.Fatalf("new worker: %v", err)
	}

	out, err := wkr.DoQuery(context.Background(), Query{
		Text: `
			LET api = NET::REST::CLIENT({
				baseUrl: @baseUrl,
				encoding: "json"
			})
			LET res = QUERY ONE "/health" IN api USING http OPTIONS {
				timeout: 1000
			}
			RETURN res.ok
		`,
		Params: map[string]interface{}{
			"baseUrl": server.URL,
		},
	})
	if err != nil {
		t.Fatalf("do query: %v", err)
	}
	if !called.Load() {
		t.Fatal("expected HTTP server to be called")
	}

	if strings.TrimSpace(string(out.Raw)) != "true" {
		t.Fatalf("expected true output, got %s", out.Raw)
	}
}

func TestDefaultFerretHTTPEgressDisabled(t *testing.T) {
	var called atomic.Bool
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		called.Store(true)
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	wkr, err := New()
	if err != nil {
		t.Fatalf("new worker: %v", err)
	}

	_, err = wkr.DoQuery(context.Background(), Query{
		Text: `RETURN IO::NET::HTTP::GET(@url)`,
		Params: map[string]interface{}{
			"url": server.URL,
		},
	})
	if err == nil {
		t.Fatal("expected HTTP request to be blocked")
	}
	if !strings.Contains(err.Error(), "outbound requests are disabled") {
		t.Fatalf("expected disabled HTTP error, got %v", err)
	}
	if called.Load() {
		t.Fatal("expected server not to be called")
	}
}

func TestRESTDisallowedHostBlockedBeforeRequest(t *testing.T) {
	var called atomic.Bool
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		called.Store(true)
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	policy := DefaultHTTPPolicy()
	policy.AllowedHosts = []string{"allowed.example"}
	policy.AllowLocalhost = true

	wkr, err := New(
		WithRESTModule(),
		WithHTTPPolicy(policy),
	)
	if err != nil {
		t.Fatalf("new worker: %v", err)
	}

	_, err = wkr.DoQuery(context.Background(), Query{
		Text: `
			LET api = NET::REST::CLIENT(@baseUrl)
			RETURN QUERY ONE "/" IN api USING http
		`,
		Params: map[string]interface{}{
			"baseUrl": server.URL,
		},
	})
	if err == nil {
		t.Fatal("expected disallowed host error")
	}
	if !strings.Contains(err.Error(), "not allowed") {
		t.Fatalf("expected not allowed error, got %v", err)
	}
	if called.Load() {
		t.Fatal("expected server not to be called")
	}
}

func TestRESTMaxRequestSizeBlocksOversizedBody(t *testing.T) {
	var called atomic.Bool
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		called.Store(true)
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	policy := testHTTPPolicy(t, server.URL)
	policy.MaxRequestSize = 4

	wkr, err := New(
		WithRESTModule(),
		WithHTTPPolicy(policy),
	)
	if err != nil {
		t.Fatalf("new worker: %v", err)
	}

	_, err = wkr.DoQuery(context.Background(), Query{
		Text: `
			LET api = NET::REST::CLIENT(@baseUrl)
			RETURN QUERY ONE "/submit" IN api USING http WITH {
				body: @body
			}
		`,
		Params: map[string]interface{}{
			"baseUrl": server.URL,
			"body":    "12345",
		},
	})
	if err == nil {
		t.Fatal("expected request body limit error")
	}
	if !strings.Contains(err.Error(), "request body exceeds limit") {
		t.Fatalf("expected request body limit error, got %v", err)
	}
	if called.Load() {
		t.Fatal("expected server not to be called")
	}
}

func TestRESTBlockedRequestHeadersAreStripped(t *testing.T) {
	var called atomic.Bool
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called.Store(true)

		if r.Header.Get("Authorization") != "" {
			t.Fatalf("expected Authorization header to be stripped, got %q", r.Header.Get("Authorization"))
		}
		if r.Header.Get("Cookie") != "" {
			t.Fatalf("expected Cookie header to be stripped, got %q", r.Header.Get("Cookie"))
		}
		if r.Header.Get("Proxy-Authorization") != "" {
			t.Fatalf("expected Proxy-Authorization header to be stripped, got %q", r.Header.Get("Proxy-Authorization"))
		}
		if r.Header.Get("X-Keep") != "ok" {
			t.Fatalf("expected X-Keep header to be forwarded, got %q", r.Header.Get("X-Keep"))
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{"ok": true})
	}))
	defer server.Close()

	wkr, err := New(
		WithRESTModule(),
		WithHTTPPolicy(testHTTPPolicy(t, server.URL)),
	)
	if err != nil {
		t.Fatalf("new worker: %v", err)
	}

	out, err := wkr.DoQuery(context.Background(), Query{
		Text: `
			LET api = NET::REST::CLIENT({
				baseUrl: @baseUrl,
				encoding: "json"
			})
			LET res = QUERY ONE "/headers" IN api USING http WITH {
				headers: {
					Authorization: "Bearer token",
					Cookie: "sid=1",
					"Proxy-Authorization": "Basic abc",
					"X-Keep": "ok"
				}
			}
			RETURN res.ok
		`,
		Params: map[string]interface{}{
			"baseUrl": server.URL,
		},
	})
	if err != nil {
		t.Fatalf("do query: %v", err)
	}
	if !called.Load() {
		t.Fatal("expected HTTP server to be called")
	}
	if strings.TrimSpace(string(out.Raw)) != "true" {
		t.Fatalf("expected true output, got %s", out.Raw)
	}
}

func testHTTPPolicy(t *testing.T, rawURL string) HTTPPolicy {
	t.Helper()

	parsed, err := url.Parse(rawURL)
	if err != nil {
		t.Fatalf("parse server url: %v", err)
	}

	policy := DefaultHTTPPolicy()
	policy.AllowedHosts = []string{parsed.Host}
	policy.AllowLocalhost = true

	return policy
}
