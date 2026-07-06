package db

import (
	"context"
	"errors"
	"testing"

	"github.com/joaovictornsv/books-cli/internal/models"
)

func TestRepositoryListSort(t *testing.T) {
	database, err := OpenMemory()
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()

	repo := NewRepository(database)
	ctx := context.Background()

	for _, title := range []string{"Charlie", "Alpha", "Bravo"} {
		_, err = repo.Create(ctx, models.Book{
			Title:          title,
			Status:         models.StatusNotStarted,
			PriorityToBuy:  0,
			EligibleToSell: 0,
			Sold:           0,
			AddedAt:        models.NowTimestamp(),
		})
		if err != nil {
			t.Fatal(err)
		}
	}

	result, err := repo.List(ctx, ListFilter{
		Sort: models.Sort{Field: models.SortFieldTitle, Order: models.SortOrderAsc},
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Books) != 3 {
		t.Fatalf("expected 3 books, got %d", len(result.Books))
	}
	if result.Books[0].Title != "Alpha" || result.Books[2].Title != "Charlie" {
		t.Fatalf("unexpected order: %+v", result.Books)
	}
}

func TestRepositorySearchAuthorTerm(t *testing.T) {
	database, err := OpenMemory()
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()

	repo := NewRepository(database)
	ctx := context.Background()
	author := "Frank Herbert"
	description := "Epic science fiction saga"

	_, err = repo.Create(ctx, models.Book{
		Title:          "The Spice Chronicles",
		Author:         &author,
		Description:    &description,
		Status:         models.StatusNotStarted,
		PriorityToBuy:  0,
		EligibleToSell: 0,
		Sold:           0,
		AddedAt:        models.NowTimestamp(),
	})
	if err != nil {
		t.Fatal(err)
	}

	results, err := repo.Search(ctx, SearchFilter{Terms: []string{"herbert"}})
	if err != nil {
		t.Fatal(err)
	}
	if len(results.Books) != 1 {
		t.Fatalf("expected 1 result matching author term, got %d", len(results.Books))
	}
}

func TestRepositoryGetByTitle(t *testing.T) {
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
	_, err = repo.Create(ctx, models.Book{
		Title:          "Children of Dune",
		Status:         models.StatusNotStarted,
		PriorityToBuy:  0,
		EligibleToSell: 0,
		Sold:           0,
		AddedAt:        models.NowTimestamp(),
	})
	if err != nil {
		t.Fatal(err)
	}

	book, err := repo.GetByTitle(ctx, TitleFilter{Title: "Dune", Exact: true})
	if err != nil {
		t.Fatal(err)
	}
	if book.Title != "Dune" {
		t.Fatalf("got title %q", book.Title)
	}

	_, err = repo.GetByTitle(ctx, TitleFilter{Title: "missing"})
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected not found, got %v", err)
	}

	_, err = repo.GetByTitle(ctx, TitleFilter{Title: "dune"})
	if !errors.Is(err, ErrAmbiguousTitle) {
		t.Fatalf("expected ambiguous title, got %v", err)
	}

	book, err = repo.GetByTitle(ctx, TitleFilter{Title: "Dune", Author: "herbert"})
	if err != nil {
		t.Fatal(err)
	}
	if book.Title != "Dune" {
		t.Fatalf("got title %q", book.Title)
	}
}
