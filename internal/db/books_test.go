package db

import (
	"context"
	"errors"
	"testing"

	"github.com/joaovictornsv/books-cli/internal/models"
)

func TestRepositoryCRUD(t *testing.T) {
	database, err := OpenMemory()
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()

	repo := NewRepository(database)
	ctx := context.Background()

	created, err := repo.Create(ctx, models.Book{
		Title:          "Dune",
		Status:         models.StatusNotStarted,
		PriorityToBuy:  0,
		EligibleToSell: 0,
		Sold:           0,
		AddedAt:        models.NowTimestamp(),
	})
	if err != nil {
		t.Fatal(err)
	}

	got, err := repo.GetByID(ctx, created.ID)
	if err != nil {
		t.Fatal(err)
	}
	if got.Title != "Dune" {
		t.Fatalf("got title %q", got.Title)
	}

	reading := models.StatusReading
	updated, err := repo.Update(ctx, created.ID, models.BookPatch{Status: &reading})
	if err != nil {
		t.Fatal(err)
	}
	if updated.StartedAt == nil {
		t.Fatal("expected started_at to be set")
	}

	read := models.StatusRead
	updated, err = repo.Update(ctx, created.ID, models.BookPatch{Status: &read})
	if err != nil {
		t.Fatal(err)
	}
	if updated.FinishedAt == nil {
		t.Fatal("expected finished_at to be set")
	}

	archived, err := repo.Archive(ctx, created.ID)
	if err != nil {
		t.Fatal(err)
	}
	if archived.Status != models.StatusArchived {
		t.Fatalf("got status %q", archived.Status)
	}

	books, err := repo.List(ctx, ListFilter{})
	if err != nil {
		t.Fatal(err)
	}
	if len(books) != 0 {
		t.Fatalf("expected archived book hidden, got %d books", len(books))
	}
}

func TestRepositorySearch(t *testing.T) {
	database, err := OpenMemory()
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()

	repo := NewRepository(database)
	ctx := context.Background()
	author := "Frank Herbert"

	_, err = repo.Create(ctx, models.Book{
		Title:          "Dune",
		Author:         &author,
		Status:         models.StatusNotStarted,
		PriorityToBuy:  0,
		EligibleToSell: 0,
		Sold:           0,
		AddedAt:        models.NowTimestamp(),
	})
	if err != nil {
		t.Fatal(err)
	}

	results, err := repo.Search(ctx, SearchFilter{Query: "dune"})
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}

	results, err = repo.Search(ctx, SearchFilter{Query: "dune", Author: "herbert"})
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result with author filter, got %d", len(results))
	}
}

func TestRepositoryListFilters(t *testing.T) {
	database, err := OpenMemory()
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()

	repo := NewRepository(database)
	ctx := context.Background()

	_, err = repo.Create(ctx, models.Book{
		Title:          "Buy Me",
		Status:         models.StatusToBuy,
		PriorityToBuy:  1,
		EligibleToSell: 0,
		Sold:           0,
		AddedAt:        models.NowTimestamp(),
	})
	if err != nil {
		t.Fatal(err)
	}

	priority := true
	toBuy := models.StatusToBuy
	books, err := repo.List(ctx, ListFilter{Status: &toBuy, PriorityToBuy: &priority})
	if err != nil {
		t.Fatal(err)
	}
	if len(books) != 1 {
		t.Fatalf("expected 1 book, got %d", len(books))
	}
}

func TestGetByIDNotFound(t *testing.T) {
	database, err := OpenMemory()
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()

	repo := NewRepository(database)
	_, err = repo.GetByID(context.Background(), 999)
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestMigrationsIdempotent(t *testing.T) {
	database, err := OpenMemory()
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()

	if err := migrate(database.SQL()); err != nil {
		t.Fatal(err)
	}
}
