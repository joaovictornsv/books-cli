package output

import (
	"fmt"
	"io"
	"strings"
	"text/tabwriter"

	"github.com/joaovictornsv/books-cli/internal/config"
	"github.com/joaovictornsv/books-cli/internal/models"
)

type TableFormatter struct{}

func (TableFormatter) PrintBook(w io.Writer, book models.Book) error {
	return printBooksTable(w, BooksPage{Books: []models.Book{book}, Total: 1})
}

func (TableFormatter) PrintBooks(w io.Writer, page BooksPage) error {
	if err := printBooksTable(w, page); err != nil {
		return err
	}
	if page.Pagination != nil {
		totalPages := (page.Total + page.Pagination.Limit - 1) / page.Pagination.Limit
		if totalPages < 1 {
			totalPages = 1
		}
		_, err := fmt.Fprintf(w, "\npage %d of %d (%d total)\n",
			page.Pagination.Page, totalPages, page.Total)
		return err
	}
	return nil
}

func (TableFormatter) PrintConfig(w io.Writer, cfg config.Config) error {
	return PrintConfigHuman(w, cfg)
}

func printBooksTable(w io.Writer, page BooksPage) error {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	if _, err := fmt.Fprintln(tw, "ID\tTITLE\tAUTHOR\tCATEGORY\tSTATUS\tPRIORITY\tSELL\tSOLD\tADDED"); err != nil {
		return err
	}
	for _, book := range page.Books {
		_, err := fmt.Fprintf(tw, "%d\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\n",
			book.ID,
			book.Title,
			derefOr(book.Author, "-"),
			categoryLabel(book.Category),
			book.Status,
			boolMark(book.PriorityToBuy),
			boolMark(book.EligibleToSell),
			boolMark(book.Sold),
			book.AddedAt,
		)
		if err != nil {
			return err
		}
	}
	return tw.Flush()
}

func boolMark(v int) string {
	if models.FromBool01(v) {
		return "Y"
	}
	return "-"
}

func derefOr(v *string, fallback string) string {
	if v == nil || strings.TrimSpace(*v) == "" {
		return fallback
	}
	return *v
}

func categoryLabel(v *models.Category) string {
	if v == nil {
		return "-"
	}
	return v.String()
}
