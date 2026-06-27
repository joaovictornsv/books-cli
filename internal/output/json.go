package output

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/joaovictornsv/books-cli/internal/config"
	"github.com/joaovictornsv/books-cli/internal/models"
)

type JSONFormatter struct{}

type booksResponse struct {
	Books []models.Book `json:"books"`
	Total int           `json:"total"`
}

func (JSONFormatter) PrintBook(w io.Writer, book models.Book) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(book)
}

func (JSONFormatter) PrintBooks(w io.Writer, books []models.Book) error {
	if books == nil {
		books = []models.Book{}
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(booksResponse{Books: books, Total: len(books)})
}

func (JSONFormatter) PrintConfig(w io.Writer, cfg config.Config) error {
	payload := map[string]any{
		"database_path": cfg.DatabasePath,
		"config_path":   cfg.ConfigPath,
		"config_exists": cfg.ConfigExists,
		"source":        cfg.Source,
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(payload)
}

func PrintConfigHuman(w io.Writer, cfg config.Config) error {
	_, err := fmt.Fprintf(w, "database_path: %s\nconfig_path: %s\nconfig_exists: %t\nsource: %s\n",
		cfg.DatabasePath, cfg.ConfigPath, cfg.ConfigExists, cfg.Source)
	return err
}
