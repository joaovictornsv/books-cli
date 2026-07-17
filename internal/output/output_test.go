package output

import (
	"bytes"
	"strings"
	"testing"

	"github.com/joaovictornsv/books-cli/internal/models"
)

func TestJSONFormatter(t *testing.T) {
	book := models.Book{
		ID:             1,
		Title:          "Dune",
		Status:         models.StatusReading,
		PriorityToBuy:  1,
		EligibleToDonate: 0,
		Donated:           0,
		AddedAt:        "2024-01-01T00:00:00Z",
	}

	var buf bytes.Buffer
	formatter := JSONFormatter{}
	if err := formatter.PrintBook(&buf, book); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), `"title": "Dune"`) {
		t.Fatalf("unexpected json: %s", buf.String())
	}

	buf.Reset()
	if err := formatter.PrintBooks(&buf, BooksPage{Books: []models.Book{book}, Total: 1}); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), `"total": 1`) {
		t.Fatalf("unexpected json: %s", buf.String())
	}

	buf.Reset()
	page := 1
	limit := 10
	if err := formatter.PrintBooks(&buf, BooksPage{
		Books:      []models.Book{book},
		Total:      25,
		Pagination: &models.Pagination{Page: page, Limit: limit},
	}); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), `"page": 1`) || !strings.Contains(buf.String(), `"limit": 10`) {
		t.Fatalf("unexpected paginated json: %s", buf.String())
	}

	buf.Reset()
	if err := formatter.PrintCount(&buf, 42); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), `"total": 42`) {
		t.Fatalf("unexpected count json: %s", buf.String())
	}

	buf.Reset()
	if err := formatter.PrintStats(&buf, models.LibraryStats{
		Year:             2025,
		ByStatus:         map[string]int{"READ": 2},
		ByCategory:       map[string]int{"FICTION": 1},
		FinishedThisYear: 3,
		PriorityWishlist: 1,
	}); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	for _, want := range []string{`"year": 2025`, `"by_status"`, `"finished_this_year": 3`, `"priority_wishlist": 1`} {
		if !strings.Contains(out, want) {
			t.Fatalf("expected %q in stats json, got: %s", want, out)
		}
	}

	buf.Reset()
	if err := formatter.PrintBackup(&buf, "/src/books.db", "/dest/backup.db"); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), `"source": "/src/books.db"`) || !strings.Contains(buf.String(), `"output": "/dest/backup.db"`) {
		t.Fatalf("unexpected backup json: %s", buf.String())
	}
}

func TestTableFormatter(t *testing.T) {
	book := models.Book{
		ID:             1,
		Title:          "Dune",
		Status:         models.StatusReading,
		PriorityToBuy:  1,
		EligibleToDonate: 0,
		Donated:           0,
		AddedAt:        "2024-01-01T00:00:00Z",
	}

	var buf bytes.Buffer
	table := TableFormatter{}
	if err := table.PrintBook(&buf, book); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), "Dune") {
		t.Fatalf("unexpected table: %s", buf.String())
	}

	buf.Reset()
	if err := table.PrintCount(&buf, 7); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), "total: 7") {
		t.Fatalf("unexpected count table: %s", buf.String())
	}

	buf.Reset()
	if err := table.PrintStats(&buf, models.LibraryStats{
		Year:             2024,
		ByStatus:         map[string]int{"READ": 1},
		FinishedThisYear: 2,
		PriorityWishlist: 3,
	}); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	for _, want := range []string{"finished_this_year (2024): 2", "priority_wishlist: 3", "by_status:", "READ"} {
		if !strings.Contains(out, want) {
			t.Fatalf("expected %q in stats table, got: %s", want, out)
		}
	}

	buf.Reset()
	if err := table.PrintBackup(&buf, "/src/books.db", "/dest/backup.db"); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), "backed up /src/books.db to /dest/backup.db") {
		t.Fatalf("unexpected backup table: %s", buf.String())
	}
}
