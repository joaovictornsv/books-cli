package output

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/joaovictornsv/books-cli/internal/models"
)

func TestParseFields(t *testing.T) {
	t.Run("valid fields", func(t *testing.T) {
		got, err := ParseFields("id, title, status")
		if err != nil {
			t.Fatal(err)
		}
		want := []string{"id", "title", "status"}
		if len(got) != len(want) {
			t.Fatalf("got %v, want %v", got, want)
		}
		for i := range want {
			if got[i] != want[i] {
				t.Fatalf("got %v, want %v", got, want)
			}
		}
	})

	t.Run("deduplicates fields", func(t *testing.T) {
		got, err := ParseFields("id,id,title")
		if err != nil {
			t.Fatal(err)
		}
		want := []string{"id", "title"}
		if len(got) != len(want) {
			t.Fatalf("got %v, want %v", got, want)
		}
	})

	t.Run("empty value", func(t *testing.T) {
		if _, err := ParseFields(""); err == nil {
			t.Fatal("expected error for empty fields")
		}
	})

	t.Run("invalid field", func(t *testing.T) {
		_, err := ParseFields("id,foo")
		if err == nil || !strings.Contains(err.Error(), `invalid field "foo"`) {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}

func TestProjectBookPreservesFieldOrder(t *testing.T) {
	author := "Herbert"
	book := models.Book{
		ID:     1,
		Title:  "Dune",
		Author: &author,
		Status: models.StatusReading,
	}

	projected := ProjectBook(book, []string{"title", "id", "author"})
	raw, err := json.Marshal(projected)
	if err != nil {
		t.Fatal(err)
	}
	if string(raw) != `{"title":"Dune","id":1,"author":"Herbert"}` {
		t.Fatalf("unexpected projection order: %s", raw)
	}
}

func TestJSONFormatterWithFields(t *testing.T) {
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
	if err := formatter.PrintBooks(&buf, BooksPage{
		Books:  []models.Book{book},
		Total:  1,
		Fields: []string{"id", "title", "status"},
	}); err != nil {
		t.Fatal(err)
	}

	output := buf.String()
	if !strings.Contains(output, `"books": [`) {
		t.Fatalf("unexpected json: %s", output)
	}
	if strings.Contains(output, `"added_at"`) {
		t.Fatalf("expected projected fields only, got: %s", output)
	}
	if !strings.Contains(output, `"title": "Dune"`) {
		t.Fatalf("unexpected json: %s", output)
	}
}
