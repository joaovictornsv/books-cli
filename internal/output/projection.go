package output

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/joaovictornsv/books-cli/internal/models"
)

var bookFieldNames = []string{
	"id",
	"title",
	"author",
	"category",
	"status",
	"priority_to_buy",
	"eligible_to_sell",
	"sold",
	"notes",
	"description",
	"added_at",
	"started_at",
	"finished_at",
}

var bookFieldSet = func() map[string]struct{} {
	set := make(map[string]struct{}, len(bookFieldNames))
	for _, name := range bookFieldNames {
		set[name] = struct{}{}
	}
	return set
}()

func ParseFields(raw string) ([]string, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil, fmt.Errorf("at least one field is required")
	}

	parts := strings.Split(raw, ",")
	fields := make([]string, 0, len(parts))
	seen := make(map[string]struct{}, len(parts))
	for _, part := range parts {
		field := strings.TrimSpace(part)
		if field == "" {
			return nil, fmt.Errorf("at least one field is required")
		}
		if _, ok := bookFieldSet[field]; !ok {
			return nil, fmt.Errorf(
				"invalid field %q: must be one of %s",
				field,
				strings.Join(bookFieldNames, ", "),
			)
		}
		if _, dup := seen[field]; dup {
			continue
		}
		seen[field] = struct{}{}
		fields = append(fields, field)
	}
	return fields, nil
}

type orderedJSON struct {
	keys   []string
	values map[string]any
}

func (o orderedJSON) MarshalJSON() ([]byte, error) {
	buf := &bytes.Buffer{}
	buf.WriteByte('{')
	for i, key := range o.keys {
		if i > 0 {
			buf.WriteByte(',')
		}
		keyJSON, err := json.Marshal(key)
		if err != nil {
			return nil, err
		}
		buf.Write(keyJSON)
		buf.WriteByte(':')
		valueJSON, err := json.Marshal(o.values[key])
		if err != nil {
			return nil, err
		}
		buf.Write(valueJSON)
	}
	buf.WriteByte('}')
	return buf.Bytes(), nil
}

func ProjectBook(book models.Book, fields []string) orderedJSON {
	values := make(map[string]any, len(fields))
	for _, field := range fields {
		values[field] = bookFieldValue(book, field)
	}
	return orderedJSON{keys: fields, values: values}
}

func ProjectBooks(books []models.Book, fields []string) []orderedJSON {
	projected := make([]orderedJSON, len(books))
	for i, book := range books {
		projected[i] = ProjectBook(book, fields)
	}
	return projected
}

func bookFieldValue(book models.Book, field string) any {
	switch field {
	case "id":
		return book.ID
	case "title":
		return book.Title
	case "author":
		if book.Author == nil {
			return nil
		}
		return *book.Author
	case "category":
		if book.Category == nil {
			return nil
		}
		return string(*book.Category)
	case "status":
		return string(book.Status)
	case "priority_to_buy":
		return book.PriorityToBuy
	case "eligible_to_sell":
		return book.EligibleToSell
	case "sold":
		return book.Sold
	case "notes":
		if book.Notes == nil {
			return nil
		}
		return *book.Notes
	case "description":
		if book.Description == nil {
			return nil
		}
		return *book.Description
	case "added_at":
		return book.AddedAt
	case "started_at":
		if book.StartedAt == nil {
			return nil
		}
		return *book.StartedAt
	case "finished_at":
		if book.FinishedAt == nil {
			return nil
		}
		return *book.FinishedAt
	default:
		return nil
	}
}
