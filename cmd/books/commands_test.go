package main

import (
	"bytes"
	"encoding/json"
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

	var resp struct {
		DatabasePath string `json:"database_path"`
		Source       string `json:"source"`
	}
	if err := json.Unmarshal(buf.Bytes(), &resp); err != nil {
		t.Fatalf("decode config JSON: %v\noutput: %s", err, buf.String())
	}
	if resp.DatabasePath == "" || resp.Source == "" {
		t.Fatalf("unexpected config JSON: %+v", resp)
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

func TestCountJSONOutput(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", home)
	dbPath := filepath.Join(home, "books.db")
	t.Setenv("BOOKS_DB", dbPath)

	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetErr(&buf)
	rootCmd.SetArgs([]string{"add", "Dune", "--json"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatal(err)
	}

	buf.Reset()
	rootCmd.SetArgs([]string{"count", "--json"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatal(err)
	}

	var resp struct {
		Total int `json:"total"`
	}
	if err := json.Unmarshal(buf.Bytes(), &resp); err != nil {
		t.Fatalf("decode count JSON: %v\noutput: %s", err, buf.String())
	}
	if resp.Total != 1 {
		t.Fatalf("expected total 1, got %+v", resp)
	}
}

func TestStatsJSONOutput(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", home)
	dbPath := filepath.Join(home, "books.db")
	t.Setenv("BOOKS_DB", dbPath)

	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetErr(&buf)
	rootCmd.SetArgs([]string{"add", "Dune", "--status", "READ", "--json"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatal(err)
	}

	buf.Reset()
	rootCmd.SetArgs([]string{"stats", "--json"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatal(err)
	}

	var resp struct {
		Year             int            `json:"year"`
		ByStatus         map[string]int `json:"by_status"`
		FinishedThisYear int            `json:"finished_this_year"`
	}
	if err := json.Unmarshal(buf.Bytes(), &resp); err != nil {
		t.Fatalf("decode stats JSON: %v\noutput: %s", err, buf.String())
	}
	if resp.Year == 0 || resp.ByStatus == nil {
		t.Fatalf("unexpected stats JSON: %+v", resp)
	}
}

func TestBackupCopiesDatabase(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", home)
	dbPath := filepath.Join(home, "books.db")
	t.Setenv("BOOKS_DB", dbPath)
	backupPath := filepath.Join(home, "backup.db")

	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetErr(&buf)
	rootCmd.SetArgs([]string{"add", "Dune", "--json"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatal(err)
	}

	buf.Reset()
	rootCmd.SetArgs([]string{"backup", "--output", backupPath, "--json"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatal(err)
	}

	var resp struct {
		Source string `json:"source"`
		Output string `json:"output"`
	}
	if err := json.Unmarshal(buf.Bytes(), &resp); err != nil {
		t.Fatalf("decode backup JSON: %v\noutput: %s", err, buf.String())
	}
	if resp.Source != dbPath || resp.Output != backupPath {
		t.Fatalf("unexpected backup JSON: %+v", resp)
	}
	if _, err := os.Stat(backupPath); err != nil {
		t.Fatalf("backup file missing: %v", err)
	}

	buf.Reset()
	rootCmd.SetArgs([]string{"backup", "--output", backupPath})
	err := rootCmd.Execute()
	if err == nil {
		t.Fatal("expected error when destination exists without --force")
	}
	if !strings.Contains(err.Error(), "already exists") {
		t.Fatalf("unexpected error: %v", err)
	}

	buf.Reset()
	rootCmd.SetArgs([]string{"backup", "--output", backupPath, "--force", "--json"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatal(err)
	}
}
