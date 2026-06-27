package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestResolvePrecedence(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", home)

	cfgDir := filepath.Join(home, ".config", "books")
	if err := os.MkdirAll(cfgDir, 0o755); err != nil {
		t.Fatal(err)
	}
	cfgPath := filepath.Join(cfgDir, "config.toml")
	if err := os.WriteFile(cfgPath, []byte(`database = "/from/config.toml"`), 0o644); err != nil {
		t.Fatal(err)
	}

	t.Setenv("BOOKS_DB", "/from/env")
	cfg, err := Resolve()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.DatabasePath != "/from/env" {
		t.Fatalf("got %q, want env path", cfg.DatabasePath)
	}
	if cfg.Source != SourceEnv {
		t.Fatalf("got source %q, want env", cfg.Source)
	}

	t.Setenv("BOOKS_DB", "")
	cfg, err = Resolve()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.DatabasePath != "/from/config.toml" {
		t.Fatalf("got %q, want config file path", cfg.DatabasePath)
	}
	if cfg.Source != SourceConfigFile {
		t.Fatalf("got source %q, want config_file", cfg.Source)
	}

	if err := os.Remove(cfgPath); err != nil {
		t.Fatal(err)
	}
	cfg, err = Resolve()
	if err != nil {
		t.Fatal(err)
	}
	wantDefault := filepath.Join(home, ".local", "share", "books", "books.db")
	if cfg.DatabasePath != wantDefault {
		t.Fatalf("got %q, want default %q", cfg.DatabasePath, wantDefault)
	}
	if cfg.Source != SourceDefault {
		t.Fatalf("got source %q, want default", cfg.Source)
	}
}
