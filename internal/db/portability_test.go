package db

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/joaovictornsv/books-cli/internal/models"
)

func TestUpdateManyRollsBackOnMissingID(t *testing.T) {
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
		EligibleToDonate: 0,
		Donated:           0,
		AddedAt:        models.NowTimestamp(),
	})
	if err != nil {
		t.Fatal(err)
	}

	read := models.StatusRead
	_, err = repo.UpdateMany(ctx, []int64{created.ID, 999}, models.BookPatch{Status: &read})
	if err == nil {
		t.Fatal("expected error for missing id in bulk update")
	}

	got, err := repo.GetByID(ctx, created.ID)
	if err != nil {
		t.Fatal(err)
	}
	if got.Status != models.StatusNotStarted {
		t.Fatalf("expected rollback, got status %q", got.Status)
	}
}

func TestUpdateManyUpdatesAllIDs(t *testing.T) {
	database, err := OpenMemory()
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()

	repo := NewRepository(database)
	ctx := context.Background()

	id1 := mustCreateBook(t, repo, ctx, "One")
	id2 := mustCreateBook(t, repo, ctx, "Two")

	read := models.StatusRead
	updated, err := repo.UpdateMany(ctx, []int64{id1, id2}, models.BookPatch{Status: &read})
	if err != nil {
		t.Fatal(err)
	}
	if len(updated) != 2 {
		t.Fatalf("expected 2 updated books, got %d", len(updated))
	}
	for _, book := range updated {
		if book.Status != models.StatusRead {
			t.Fatalf("book %d status = %q", book.ID, book.Status)
		}
	}
}

func TestExportImportJSONRoundTrip(t *testing.T) {
	database, err := OpenMemory()
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()

	repo := NewRepository(database)
	ctx := context.Background()
	author := "Frank Herbert"
	created, err := repo.Create(ctx, models.Book{
		Title:          "Dune",
		Author:         &author,
		Status:         models.StatusToBuy,
		PriorityToBuy:  1,
		EligibleToDonate: 0,
		Donated:           0,
		AddedAt:        models.NowTimestamp(),
	})
	if err != nil {
		t.Fatal(err)
	}

	books, err := repo.ListAll(ctx, true)
	if err != nil {
		t.Fatal(err)
	}
	data, err := MarshalBooksJSON(books)
	if err != nil {
		t.Fatal(err)
	}

	parsed, err := UnmarshalBooksJSON(data)
	if err != nil {
		t.Fatal(err)
	}
	if len(parsed) != 1 || parsed[0].Title != "Dune" {
		t.Fatalf("unexpected parsed books: %+v", parsed)
	}

	database2, err := OpenMemory()
	if err != nil {
		t.Fatal(err)
	}
	defer database2.Close()
	repo2 := NewRepository(database2)

	result, err := repo2.ImportBooks(ctx, parsed, false)
	if err != nil {
		t.Fatal(err)
	}
	if result.Created != 1 || result.Updated != 0 {
		t.Fatalf("unexpected import result: %+v", result)
	}

	got, err := repo2.GetByID(ctx, created.ID)
	if err != nil {
		t.Fatal(err)
	}
	if got.Title != "Dune" || got.PriorityToBuy != 1 {
		t.Fatalf("unexpected imported book: %+v", got)
	}
}

func TestImportDryRunDoesNotWrite(t *testing.T) {
	database, err := OpenMemory()
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()

	repo := NewRepository(database)
	ctx := context.Background()

	books := []models.Book{{
		Title:          "Dune",
		Status:         models.StatusToBuy,
		PriorityToBuy:  0,
		EligibleToDonate: 0,
		Donated:           0,
		AddedAt:        models.NowTimestamp(),
	}}

	result, err := repo.ImportBooks(ctx, books, true)
	if err != nil {
		t.Fatal(err)
	}
	if !result.DryRun || result.Created != 1 {
		t.Fatalf("unexpected dry-run result: %+v", result)
	}

	count, err := repo.Count(ctx, ListFilter{})
	if err != nil {
		t.Fatal(err)
	}
	if count != 0 {
		t.Fatalf("expected no books after dry run, got %d", count)
	}
}

func TestImportUpsertUpdatesExisting(t *testing.T) {
	database, err := OpenMemory()
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()

	repo := NewRepository(database)
	ctx := context.Background()

	created, err := repo.Create(ctx, models.Book{
		Title:          "Dune",
		Status:         models.StatusToBuy,
		PriorityToBuy:  0,
		EligibleToDonate: 0,
		Donated:           0,
		AddedAt:        models.NowTimestamp(),
	})
	if err != nil {
		t.Fatal(err)
	}

	created.Title = "Dune Messiah"
	result, err := repo.ImportBooks(ctx, []models.Book{created}, false)
	if err != nil {
		t.Fatal(err)
	}
	if result.Updated != 1 || result.Created != 0 {
		t.Fatalf("unexpected upsert result: %+v", result)
	}

	got, err := repo.GetByID(ctx, created.ID)
	if err != nil {
		t.Fatal(err)
	}
	if got.Title != "Dune Messiah" {
		t.Fatalf("title = %q", got.Title)
	}
}

func TestWriteAndReadBooksCSV(t *testing.T) {
	author := "Frank Herbert"
	book := models.Book{
		ID:             1,
		Title:          "Dune",
		Author:         &author,
		Status:         models.StatusToBuy,
		PriorityToBuy:  1,
		EligibleToDonate: 0,
		Donated:           0,
		AddedAt:        "2024-01-01T00:00:00Z",
	}

	var buf bytes.Buffer
	if err := WriteBooksCSV(&buf, []models.Book{book}); err != nil {
		t.Fatal(err)
	}

	parsed, err := ReadBooksCSV(strings.NewReader(buf.String()))
	if err != nil {
		t.Fatal(err)
	}
	if len(parsed) != 1 || parsed[0].Title != "Dune" || parsed[0].PriorityToBuy != 1 {
		t.Fatalf("unexpected csv round-trip: %+v", parsed)
	}
}

func mustCreateBook(t *testing.T, repo *Repository, ctx context.Context, title string) int64 {
	t.Helper()
	created, err := repo.Create(ctx, models.Book{
		Title:          title,
		Status:         models.StatusNotStarted,
		PriorityToBuy:  0,
		EligibleToDonate: 0,
		Donated:           0,
		AddedAt:        models.NowTimestamp(),
	})
	if err != nil {
		t.Fatal(err)
	}
	return created.ID
}
