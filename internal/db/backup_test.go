package db

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/joaovictornsv/books-cli/internal/models"
)

func TestBackupCreatesConsistentCopy(t *testing.T) {
	dir := t.TempDir()
	source := filepath.Join(dir, "books.db")
	dest := filepath.Join(dir, "backup.db")

	database, err := Open(source)
	if err != nil {
		t.Fatal(err)
	}
	repo := NewRepository(database)
	ctx := context.Background()

	_, err = repo.Create(ctx, models.Book{
		Title:          "Dune",
		Status:         models.StatusRead,
		PriorityToBuy:  0,
		EligibleToDonate: 0,
		Donated:           0,
		AddedAt:        models.NowTimestamp(),
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := database.Close(); err != nil {
		t.Fatal(err)
	}

	if err := Backup(ctx, source, dest, false); err != nil {
		t.Fatal(err)
	}

	backupDB, err := Open(dest)
	if err != nil {
		t.Fatal(err)
	}
	defer backupDB.Close()

	total, err := NewRepository(backupDB).Count(ctx, ListFilter{})
	if err != nil {
		t.Fatal(err)
	}
	if total != 1 {
		t.Fatalf("expected 1 book in backup, got %d", total)
	}
}

func TestBackupRejectsExistingDestination(t *testing.T) {
	dir := t.TempDir()
	source := filepath.Join(dir, "books.db")
	dest := filepath.Join(dir, "backup.db")

	database, err := Open(source)
	if err != nil {
		t.Fatal(err)
	}
	if err := database.Close(); err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()
	if err := Backup(ctx, source, dest, false); err != nil {
		t.Fatal(err)
	}

	err = Backup(ctx, source, dest, false)
	if err == nil {
		t.Fatal("expected error when destination exists")
	}
}

func TestBackupForceOverwrites(t *testing.T) {
	dir := t.TempDir()
	source := filepath.Join(dir, "books.db")
	dest := filepath.Join(dir, "backup.db")

	database, err := Open(source)
	if err != nil {
		t.Fatal(err)
	}
	if err := database.Close(); err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()
	if err := Backup(ctx, source, dest, false); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(dest, []byte("stale"), 0o600); err != nil {
		t.Fatal(err)
	}

	if err := Backup(ctx, source, dest, true); err != nil {
		t.Fatal(err)
	}

	info, err := os.Stat(dest)
	if err != nil {
		t.Fatal(err)
	}
	if info.Size() <= 5 {
		t.Fatalf("expected backup to replace stale file, size=%d", info.Size())
	}
}

func TestBackupMissingSource(t *testing.T) {
	dir := t.TempDir()
	err := Backup(context.Background(), filepath.Join(dir, "missing.db"), filepath.Join(dir, "backup.db"), false)
	if err == nil {
		t.Fatal("expected error for missing source")
	}
}
