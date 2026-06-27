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
	if !strings.Contains(err.Error(), "search query cannot be empty") {
		t.Fatalf("unexpected error: %v", err)
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
