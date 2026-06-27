package output

import (
	"io"

	"github.com/joaovictornsv/books-cli/internal/config"
	"github.com/joaovictornsv/books-cli/internal/models"
)

type Formatter interface {
	PrintBook(w io.Writer, book models.Book) error
	PrintBooks(w io.Writer, books []models.Book) error
	PrintConfig(w io.Writer, cfg config.Config) error
}

func New(json bool) Formatter {
	if json {
		return JSONFormatter{}
	}
	return TableFormatter{}
}
