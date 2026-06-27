package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/joaovictornsv/books-cli/internal/db"
	"github.com/spf13/cobra"
)

var (
	searchAuthor string
	searchPage   int
	searchLimit  int
)

var searchCmd = &cobra.Command{
	Use:   "search [query]",
	Short: "Search books by title, description, and optional author",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		query := strings.TrimSpace(args[0])
		if query == "" {
			return fmt.Errorf("search query cannot be empty")
		}

		filter := db.SearchFilter{
			Query:  query,
			Author: searchAuthor,
		}

		pagination, err := paginationFromFlags(cmd, &searchPage, &searchLimit)
		if err != nil {
			return err
		}
		filter.Pagination = pagination

		return runWithRepo(cmd.Context(), func(ctx context.Context, repo *db.Repository) error {
			result, err := repo.Search(ctx, filter)
			if err != nil {
				return err
			}
			return printBooksResult(cmd, result, pagination)
		})
	},
}

func init() {
	searchCmd.Flags().StringVar(&searchAuthor, "author", "", "Filter by author substring")
	addPaginationFlags(searchCmd, &searchPage, &searchLimit)
	addFieldsFlag(searchCmd)
	rootCmd.AddCommand(searchCmd)
}
