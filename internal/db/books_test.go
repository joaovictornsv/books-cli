package db

import (
	"context"
	"errors"
	"fmt"
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
	if updated.StartedAt != nil {
		t.Fatalf("expected started_at to remain unset on status change, got %v", updated.StartedAt)
	}

	startedAt := "2024-06-01T12:00:00Z"
	updated, err = repo.Update(ctx, created.ID, models.BookPatch{StartedAt: &startedAt})
	if err != nil {
		t.Fatal(err)
	}
	if updated.StartedAt == nil || *updated.StartedAt != startedAt {
		t.Fatalf("started_at = %v, want %q", updated.StartedAt, startedAt)
	}

	read := models.StatusRead
	updated, err = repo.Update(ctx, created.ID, models.BookPatch{Status: &read})
	if err != nil {
		t.Fatal(err)
	}
	if updated.FinishedAt != nil {
		t.Fatalf("expected finished_at to remain unset on status change, got %v", updated.FinishedAt)
	}

	finishedAt := "2024-06-02T12:00:00Z"
	updated, err = repo.Update(ctx, created.ID, models.BookPatch{FinishedAt: &finishedAt})
	if err != nil {
		t.Fatal(err)
	}
	if updated.FinishedAt == nil || *updated.FinishedAt != finishedAt {
		t.Fatalf("finished_at = %v, want %q", updated.FinishedAt, finishedAt)
	}

	updated, err = repo.Update(ctx, created.ID, models.BookPatch{ClearStartedAt: true})
	if err != nil {
		t.Fatal(err)
	}
	if updated.StartedAt != nil {
		t.Fatalf("expected started_at cleared, got %v", updated.StartedAt)
	}

	archivedStatus := models.StatusArchived
	archived, err := repo.Update(ctx, created.ID, models.BookPatch{Status: &archivedStatus})
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
	if len(books.Books) != 0 {
		t.Fatalf("expected archived book hidden, got %d books", len(books.Books))
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

	results, err := repo.Search(ctx, SearchFilter{Terms: []string{"dune"}})
	if err != nil {
		t.Fatal(err)
	}
	if len(results.Books) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results.Books))
	}

	results, err = repo.Search(ctx, SearchFilter{Terms: []string{"dune"}, Author: "herbert"})
	if err != nil {
		t.Fatal(err)
	}
	if len(results.Books) != 1 {
		t.Fatalf("expected 1 result with author filter, got %d", len(results.Books))
	}
}

func TestRepositorySearchDescription(t *testing.T) {
	database, err := OpenMemory()
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()

	repo := NewRepository(database)
	ctx := context.Background()
	description := "Epic science fiction saga set on the desert planet Arrakis"

	_, err = repo.Create(ctx, models.Book{
		Title:          "Dune",
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

	results, err := repo.Search(ctx, SearchFilter{Terms: []string{"arrakis"}})
	if err != nil {
		t.Fatal(err)
	}
	if len(results.Books) != 1 {
		t.Fatalf("expected 1 result searching description, got %d", len(results.Books))
	}
}

func TestRepositorySearchMultipleTerms(t *testing.T) {
	database, err := OpenMemory()
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()

	repo := NewRepository(database)
	ctx := context.Background()

	for _, title := range []string{"Dune", "The Hobbit"} {
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

	results, err := repo.Search(ctx, SearchFilter{Terms: []string{"dune", "hobbit"}})
	if err != nil {
		t.Fatal(err)
	}
	if len(results.Books) != 2 {
		t.Fatalf("expected 2 results for OR search, got %d", len(results.Books))
	}

	results, err = repo.Search(ctx, SearchFilter{Terms: []string{"dune", "nonexistent"}})
	if err != nil {
		t.Fatal(err)
	}
	if len(results.Books) != 1 {
		t.Fatalf("expected 1 result when only one term matches, got %d", len(results.Books))
	}
}

func TestRepositorySearchCategory(t *testing.T) {
	database, err := OpenMemory()
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()

	repo := NewRepository(database)
	ctx := context.Background()
	fiction := models.CategoryFiction
	history := models.CategoryHistory

	_, err = repo.Create(ctx, models.Book{
		Title:          "Dune",
		Category:       &fiction,
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
		Title:          "Sapiens",
		Category:       &history,
		Status:         models.StatusNotStarted,
		PriorityToBuy:  0,
		EligibleToSell: 0,
		Sold:           0,
		AddedAt:        models.NowTimestamp(),
	})
	if err != nil {
		t.Fatal(err)
	}

	results, err := repo.Search(ctx, SearchFilter{Terms: []string{"e"}, Category: &fiction})
	if err != nil {
		t.Fatal(err)
	}
	if len(results.Books) != 1 || results.Books[0].Title != "Dune" {
		t.Fatalf("expected only Dune, got %+v", results.Books)
	}
}

func TestRepositoryCheck(t *testing.T) {
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
		Status:         models.StatusToBuy,
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
		Status:         models.StatusToBuy,
		PriorityToBuy:  0,
		EligibleToSell: 0,
		Sold:           0,
		AddedAt:        models.NowTimestamp(),
	})
	if err != nil {
		t.Fatal(err)
	}
	_, err = repo.Create(ctx, models.Book{
		Title:          "Dune Messiah",
		Status:         models.StatusArchived,
		PriorityToBuy:  0,
		EligibleToSell: 0,
		Sold:           0,
		AddedAt:        models.NowTimestamp(),
	})
	if err != nil {
		t.Fatal(err)
	}

	results, err := repo.Check(ctx, CheckFilter{Title: "dune"})
	if err != nil {
		t.Fatal(err)
	}
	if len(results.Books) != 2 {
		t.Fatalf("expected 2 substring matches, got %d", len(results.Books))
	}

	results, err = repo.Check(ctx, CheckFilter{Title: "Dune", Exact: true})
	if err != nil {
		t.Fatal(err)
	}
	if len(results.Books) != 1 || results.Books[0].Title != "Dune" {
		t.Fatalf("expected exact Dune match, got %+v", results.Books)
	}

	results, err = repo.Check(ctx, CheckFilter{Title: "Dune", Author: "herbert"})
	if err != nil {
		t.Fatal(err)
	}
	if len(results.Books) != 1 {
		t.Fatalf("expected 1 match with author filter, got %d", len(results.Books))
	}

	results, err = repo.Check(ctx, CheckFilter{Title: "Dune Messiah"})
	if err != nil {
		t.Fatal(err)
	}
	if len(results.Books) != 0 {
		t.Fatalf("expected archived book excluded, got %d", len(results.Books))
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
	if len(books.Books) != 1 {
		t.Fatalf("expected 1 book, got %d", len(books.Books))
	}
}

func TestRepositoryCategory(t *testing.T) {
	database, err := OpenMemory()
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()

	repo := NewRepository(database)
	ctx := context.Background()
	category := models.CategoryBiography

	created, err := repo.Create(ctx, models.Book{
		Title:          "Elon Musk",
		Category:       &category,
		Status:         models.StatusToBuy,
		PriorityToBuy:  0,
		EligibleToSell: 0,
		Sold:           0,
		AddedAt:        models.NowTimestamp(),
	})
	if err != nil {
		t.Fatal(err)
	}
	if created.Category == nil || *created.Category != models.CategoryBiography {
		t.Fatalf("expected BIOGRAPHY category, got %v", created.Category)
	}

	updated, err := repo.Update(ctx, created.ID, models.BookPatch{ClearCategory: true})
	if err != nil {
		t.Fatal(err)
	}
	if updated.Category != nil {
		t.Fatalf("expected category cleared, got %v", updated.Category)
	}

	fiction := models.CategoryFiction
	updated, err = repo.Update(ctx, created.ID, models.BookPatch{Category: &fiction})
	if err != nil {
		t.Fatal(err)
	}
	if updated.Category == nil || *updated.Category != models.CategoryFiction {
		t.Fatalf("expected FICTION category, got %v", updated.Category)
	}
}

func TestRepositoryClearAuthor(t *testing.T) {
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
		Status:         models.StatusNotStarted,
		PriorityToBuy:  0,
		EligibleToSell: 0,
		Sold:           0,
		AddedAt:        models.NowTimestamp(),
	})
	if err != nil {
		t.Fatal(err)
	}

	updated, err := repo.Update(ctx, created.ID, models.BookPatch{ClearAuthor: true})
	if err != nil {
		t.Fatal(err)
	}
	if updated.Author != nil {
		t.Fatalf("expected author cleared, got %v", updated.Author)
	}
}

func TestRepositoryDelete(t *testing.T) {
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

	deleted, err := repo.Delete(ctx, created.ID)
	if err != nil {
		t.Fatal(err)
	}
	if deleted.Title != "Dune" {
		t.Fatalf("got title %q", deleted.Title)
	}

	_, err = repo.GetByID(ctx, created.ID)
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound after delete, got %v", err)
	}

	books, err := repo.List(ctx, ListFilter{IncludeArchived: true})
	if err != nil {
		t.Fatal(err)
	}
	if len(books.Books) != 0 {
		t.Fatalf("expected 0 books after delete, got %d", len(books.Books))
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

func TestRepositoryPagination(t *testing.T) {
	database, err := OpenMemory()
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()

	repo := NewRepository(database)
	ctx := context.Background()

	for i := 1; i <= 5; i++ {
		_, err := repo.Create(ctx, models.Book{
			Title:          fmt.Sprintf("Book %d", i),
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

	page2 := models.Pagination{Page: 2, Limit: 2}
	result, err := repo.List(ctx, ListFilter{Pagination: &page2})
	if err != nil {
		t.Fatal(err)
	}
	if result.Total != 5 {
		t.Fatalf("expected total 5, got %d", result.Total)
	}
	if len(result.Books) != 2 {
		t.Fatalf("expected 2 books on page 2, got %d", len(result.Books))
	}
	if result.Books[0].Title != "Book 3" {
		t.Fatalf("expected Book 3 first on page 2, got %q", result.Books[0].Title)
	}
}

func TestRepositoryCountAndStats(t *testing.T) {
	database, err := OpenMemory()
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()

	repo := NewRepository(database)
	ctx := context.Background()

	fiction := models.CategoryFiction
	priority := true
	_, err = repo.Create(ctx, models.Book{
		Title:          "Dune",
		Category:       &fiction,
		Status:         models.StatusRead,
		PriorityToBuy:  0,
		EligibleToSell: 0,
		Sold:           0,
		AddedAt:        models.NowTimestamp(),
		FinishedAt:     strPtr("2026-03-01T12:00:00Z"),
	})
	if err != nil {
		t.Fatal(err)
	}
	_, err = repo.Create(ctx, models.Book{
		Title:          "Foundation",
		Category:       &fiction,
		Status:         models.StatusToBuy,
		PriorityToBuy:  1,
		EligibleToSell: 0,
		Sold:           0,
		AddedAt:        models.NowTimestamp(),
	})
	if err != nil {
		t.Fatal(err)
	}
	_, err = repo.Create(ctx, models.Book{
		Title:          "Clean Code",
		Status:         models.StatusNotStarted,
		PriorityToBuy:  0,
		EligibleToSell: 0,
		Sold:           0,
		AddedAt:        models.NowTimestamp(),
	})
	if err != nil {
		t.Fatal(err)
	}

	total, err := repo.Count(ctx, ListFilter{})
	if err != nil {
		t.Fatal(err)
	}
	if total != 3 {
		t.Fatalf("expected total 3, got %d", total)
	}

	read := models.StatusRead
	readCount, err := repo.Count(ctx, ListFilter{Status: &read})
	if err != nil {
		t.Fatal(err)
	}
	if readCount != 1 {
		t.Fatalf("expected 1 read book, got %d", readCount)
	}

	fictionCount, err := repo.Count(ctx, ListFilter{Category: &fiction})
	if err != nil {
		t.Fatal(err)
	}
	if fictionCount != 2 {
		t.Fatalf("expected 2 fiction books, got %d", fictionCount)
	}

	priorityCount, err := repo.Count(ctx, ListFilter{PriorityToBuy: &priority})
	if err != nil {
		t.Fatal(err)
	}
	if priorityCount != 1 {
		t.Fatalf("expected 1 priority book, got %d", priorityCount)
	}

	stats, err := repo.Stats(ctx, 2026)
	if err != nil {
		t.Fatal(err)
	}
	if stats.ByStatus[string(models.StatusRead)] != 1 {
		t.Fatalf("expected 1 read in stats, got %v", stats.ByStatus)
	}
	if stats.ByCategory[string(models.CategoryFiction)] != 2 {
		t.Fatalf("expected 2 fiction in stats, got %v", stats.ByCategory)
	}
	if stats.FinishedThisYear != 1 {
		t.Fatalf("expected 1 finished in 2026, got %d", stats.FinishedThisYear)
	}
	if stats.Year != 2026 {
		t.Fatalf("expected year 2026, got %d", stats.Year)
	}
	if stats.PriorityWishlist != 1 {
		t.Fatalf("expected 1 priority wishlist, got %d", stats.PriorityWishlist)
	}
}

func strPtr(s string) *string {
	return &s
}
