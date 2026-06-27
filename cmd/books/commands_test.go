package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSearchRejectsWhitespaceQuery(t *testing.T) {
	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetErr(&buf)
	rootCmd.SetArgs([]string{"search", "   "})

	err := rootCmd.Execute()
	if err == nil {
		t.Fatal("expected error for whitespace-only search query")
	}
	if !strings.Contains(err.Error(), "at least one search term is required") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSearchRejectsNoTerms(t *testing.T) {
	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetErr(&buf)
	rootCmd.SetArgs([]string{"search"})

	err := rootCmd.Execute()
	if err == nil {
		t.Fatal("expected error when no search terms provided")
	}
	if !strings.Contains(err.Error(), "at least one search term is required") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCollectSearchTerms(t *testing.T) {
	terms, err := collectSearchTerms([]string{"dune"}, []string{"hobbit", "  "})
	if err != nil {
		t.Fatal(err)
	}
	if len(terms) != 2 || terms[0] != "dune" || terms[1] != "hobbit" {
		t.Fatalf("unexpected terms: %v", terms)
	}

	_, err = collectSearchTerms(nil, []string{"  "})
	if err == nil {
		t.Fatal("expected error for empty terms")
	}
}

func TestConfigJSONOutput(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", home)
	t.Setenv("BOOKS_DB", "")

	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetErr(&buf)
	rootCmd.SetArgs([]string{"--json", "config"})

	if err := rootCmd.Execute(); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), `"database_path"`) {
		t.Fatalf("expected JSON config output, got: %s", buf.String())
	}
}

func TestAddValidatesBeforeDB(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", home)
	dbPath := filepath.Join(home, "books.db")
	t.Setenv("BOOKS_DB", dbPath)

	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetErr(&buf)
	rootCmd.SetArgs([]string{"add", "   "})

	err := rootCmd.Execute()
	if err == nil {
		t.Fatal("expected validation error for empty title")
	}
	if !strings.Contains(err.Error(), "title is required") {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, statErr := os.Stat(dbPath); statErr == nil {
		t.Fatal("expected database not to be created when validation fails early")
	}
}
