package output

import (
	"io"

	"github.com/joaovictornsv/books-cli/internal/config"
	"github.com/joaovictornsv/books-cli/internal/models"
)

type BooksPage struct {
	Books      []models.Book
	Total      int
	Pagination *models.Pagination
	Fields     []string
}

type Formatter interface {
	PrintBook(w io.Writer, book models.Book) error
	PrintBooks(w io.Writer, page BooksPage) error
	PrintConfig(w io.Writer, cfg config.Config) error
}

func New(json bool) Formatter {
	if json {
		return JSONFormatter{}
	}
	return TableFormatter{}
}
