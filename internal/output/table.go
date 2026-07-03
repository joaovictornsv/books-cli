package output

import (
	"fmt"
	"io"
	"sort"
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

func (TableFormatter) PrintCount(w io.Writer, total int) error {
	_, err := fmt.Fprintf(w, "total: %d\n", total)
	return err
}

func (TableFormatter) PrintStats(w io.Writer, stats models.LibraryStats) error {
	if _, err := fmt.Fprintf(w, "finished_this_year (%d): %d\n", stats.Year, stats.FinishedThisYear); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(w, "priority_wishlist: %d\n", stats.PriorityWishlist); err != nil {
		return err
	}

	if err := printStatsGroup(w, "by_status", stats.ByStatus); err != nil {
		return err
	}
	return printStatsGroup(w, "by_category", stats.ByCategory)
}

func printStatsGroup(w io.Writer, label string, counts map[string]int) error {
	if _, err := fmt.Fprintf(w, "\n%s:\n", label); err != nil {
		return err
	}
	keys := make([]string, 0, len(counts))
	for k := range counts {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	for _, k := range keys {
		if _, err := fmt.Fprintf(tw, "  %s\t%d\n", k, counts[k]); err != nil {
			return err
		}
	}
	return tw.Flush()
}

func (TableFormatter) PrintBackup(w io.Writer, source, dest string) error {
	_, err := fmt.Fprintf(w, "backed up %s to %s\n", source, dest)
	return err
}

func (TableFormatter) PrintSchema(w io.Writer, schema models.SchemaDocument) error {
	_, err := fmt.Fprintf(w, "%d statuses, %d categories, %d fields (use --json for full schema)\n",
		len(schema.Statuses), len(schema.Categories), len(schema.Fields))
	return err
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
