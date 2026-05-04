package main

import (
	"strings"
	"testing"
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
