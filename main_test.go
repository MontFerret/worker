package main

import (
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/MontFerret/worker/pkg/worker"
)

func TestResolveFSRootDefaultsToCurrentWorkingDirectory(t *testing.T) {
	wd := t.TempDir()
	t.Chdir(wd)

	got, err := resolveFSRoot("", false)
	if err != nil {
		t.Fatalf("resolve fs root: %v", err)
	}

	if got != wd {
		t.Fatalf("expected fs root to default to %q, got %q", wd, got)
	}
}

func TestResolveFSRootPreservesExplicitPath(t *testing.T) {
	got, err := resolveFSRoot("  ./data  ", true)
	if err != nil {
		t.Fatalf("resolve fs root: %v", err)
	}

	if got != "./data" {
		t.Fatalf("expected explicit fs root to be trimmed, got %q", got)
	}
}

func TestResolveFSRootRejectsExplicitBlankPath(t *testing.T) {
	_, err := resolveFSRoot(" \t\n ", true)
	if err == nil {
		t.Fatal("expected explicit blank fs root to fail")
	}

	if !strings.Contains(err.Error(), "fs root cannot be empty") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSplitCSVTrimsAndDropsEmptyValues(t *testing.T) {
	got := splitCSV(" example.com, , api.example.com:8443,\tinternal.example ")
	want := []string{"example.com", "api.example.com:8443", "internal.example"}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("expected %#v, got %#v", want, got)
	}
}

func TestNewHTTPPolicyFromConfigRejectsAllowlistAndAllowAll(t *testing.T) {
	_, err := newHTTPPolicyFromConfig(httpPolicyConfig{
		AllowedHosts:    "example.com",
		AllowAllHosts:   true,
		Timeout:         10 * time.Second,
		MaxRequestSize:  1024,
		MaxResponseSize: 2048,
		FollowRedirects: true,
	})
	if err == nil {
		t.Fatal("expected conflicting host policy to fail")
	}
	if !strings.Contains(err.Error(), "cannot both be set") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNewHTTPPolicyFromConfigUsesDefaultPolicyValues(t *testing.T) {
	defaults := worker.DefaultHTTPPolicy()
	policy, err := newHTTPPolicyFromConfig(httpPolicyConfig{
		BlockedRequestHeaders: strings.Join(defaults.BlockedRequestHeaders, ","),
		Timeout:               defaults.Timeout,
		MaxRequestSize:        defaults.MaxRequestSize,
		MaxResponseSize:       defaults.MaxResponseSize,
		MaxRedirects:          defaults.MaxRedirects,
		FollowRedirects:       defaults.FollowRedirects,
		AllowAllHosts:         defaults.AllowAllHosts,
		AllowLocalhost:        defaults.AllowLocalhost,
		AllowPrivateNetworks:  defaults.AllowPrivateNetworks,
	})
	if err != nil {
		t.Fatalf("new http policy: %v", err)
	}

	if !reflect.DeepEqual(policy.AllowedSchemes, defaults.AllowedSchemes) {
		t.Fatalf("expected default schemes %#v, got %#v", defaults.AllowedSchemes, policy.AllowedSchemes)
	}
	if !reflect.DeepEqual(policy.BlockedRequestHeaders, defaults.BlockedRequestHeaders) {
		t.Fatalf("expected default blocked headers %#v, got %#v", defaults.BlockedRequestHeaders, policy.BlockedRequestHeaders)
	}
	if policy.Timeout != defaults.Timeout {
		t.Fatalf("expected timeout %s, got %s", defaults.Timeout, policy.Timeout)
	}
	if policy.MaxRequestSize != defaults.MaxRequestSize {
		t.Fatalf("expected max request size %d, got %d", defaults.MaxRequestSize, policy.MaxRequestSize)
	}
	if policy.MaxResponseSize != defaults.MaxResponseSize {
		t.Fatalf("expected max response size %d, got %d", defaults.MaxResponseSize, policy.MaxResponseSize)
	}
	if policy.MaxRedirects != defaults.MaxRedirects {
		t.Fatalf("expected max redirects %d, got %d", defaults.MaxRedirects, policy.MaxRedirects)
	}
	if policy.FollowRedirects != defaults.FollowRedirects {
		t.Fatalf("expected follow redirects %t, got %t", defaults.FollowRedirects, policy.FollowRedirects)
	}
	if policy.AllowLocalhost != defaults.AllowLocalhost {
		t.Fatalf("expected allow localhost %t, got %t", defaults.AllowLocalhost, policy.AllowLocalhost)
	}
	if policy.AllowPrivateNetworks != defaults.AllowPrivateNetworks {
		t.Fatalf("expected allow private networks %t, got %t", defaults.AllowPrivateNetworks, policy.AllowPrivateNetworks)
	}
}
