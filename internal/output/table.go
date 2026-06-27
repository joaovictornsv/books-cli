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
	return printBooksTable(w, []models.Book{book})
}

func (TableFormatter) PrintBooks(w io.Writer, books []models.Book) error {
	return printBooksTable(w, books)
}

func (TableFormatter) PrintConfig(w io.Writer, cfg config.Config) error {
	return PrintConfigHuman(w, cfg)
}

func printBooksTable(w io.Writer, books []models.Book) error {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	_, _ = fmt.Fprintln(tw, "ID\tTITLE\tAUTHOR\tSTATUS\tPRIORITY\tSELL\tSOLD\tADDED")
	for _, book := range books {
		_, err := fmt.Fprintf(tw, "%d\t%s\t%s\t%s\t%s\t%s\t%s\t%s\n",
			book.ID,
			book.Title,
			derefOr(book.Author, "-"),
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
