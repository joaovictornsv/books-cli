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
		EligibleToSell: 0,
		Sold:           0,
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
	if err := formatter.PrintBooks(&buf, []models.Book{book}); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), `"total": 1`) {
		t.Fatalf("unexpected json: %s", buf.String())
	}
}

func TestTableFormatter(t *testing.T) {
	book := models.Book{
		ID:             1,
		Title:          "Dune",
		Status:         models.StatusReading,
		PriorityToBuy:  1,
		EligibleToSell: 0,
		Sold:           0,
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
}
