package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/joaovictornsv/books-cli/internal/db"
	"github.com/spf13/cobra"
)

var (
	checkTitle  string
	checkAuthor string
	checkExact  bool
)

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Find likely duplicate books by title before adding",
	Long: `Find books that may duplicate a title you plan to add.

Matches title only (not description). Use --exact for case-insensitive exact title match;
default is case-insensitive substring on title. Optional --author narrows results.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		title := strings.TrimSpace(checkTitle)
		if title == "" {
			return fmt.Errorf("--title is required")
		}

		filter := db.CheckFilter{
			Title:  title,
			Author: checkAuthor,
			Exact:  checkExact,
		}

		return runWithRepo(cmd.Context(), func(ctx context.Context, repo *db.Repository) error {
			result, err := repo.Check(ctx, filter)
			if err != nil {
				return err
			}
			return printBooksResult(cmd, result, nil)
		})
	},
}

func init() {
	checkCmd.Flags().StringVar(&checkTitle, "title", "", "Title to check for duplicates (required)")
	checkCmd.Flags().StringVar(&checkAuthor, "author", "", "Filter by author substring")
	checkCmd.Flags().BoolVar(&checkExact, "exact", false, "Case-insensitive exact title match")
	rootCmd.AddCommand(checkCmd)
}
