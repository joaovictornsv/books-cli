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
	Books any `json:"books"`
	Total int `json:"total"`
	Page  *int `json:"page,omitempty"`
	Limit *int `json:"limit,omitempty"`
}

func (JSONFormatter) PrintBook(w io.Writer, book models.Book) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(book)
}

func (JSONFormatter) PrintBooks(w io.Writer, page BooksPage) error {
	books := page.Books
	if books == nil {
		books = []models.Book{}
	}

	var booksPayload any = books
	if len(page.Fields) > 0 {
		booksPayload = ProjectBooks(books, page.Fields)
	}

	resp := booksResponse{
		Books: booksPayload,
		Total: page.Total,
	}
	if page.Pagination != nil {
		resp.Page = &page.Pagination.Page
		resp.Limit = &page.Pagination.Limit
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(resp)
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

func (JSONFormatter) PrintCount(w io.Writer, total int) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(map[string]int{"total": total})
}

func (JSONFormatter) PrintStats(w io.Writer, stats models.LibraryStats) error {
	payload := map[string]any{
		"year":               stats.Year,
		"by_status":          stats.ByStatus,
		"by_category":        stats.ByCategory,
		"finished_this_year": stats.FinishedThisYear,
		"priority_wishlist":  stats.PriorityWishlist,
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(payload)
}

func (JSONFormatter) PrintBackup(w io.Writer, source, dest string) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(map[string]string{
		"source": source,
		"output": dest,
	})
}

func PrintConfigHuman(w io.Writer, cfg config.Config) error {
	_, err := fmt.Fprintf(w, "database_path: %s\nconfig_path: %s\nconfig_exists: %t\nsource: %s\n",
		cfg.DatabasePath, cfg.ConfigPath, cfg.ConfigExists, cfg.Source)
	return err
}
