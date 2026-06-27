package main

import (
	"context"

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
	Short: "Search books by title and optional author",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		filter := db.SearchFilter{
			Query:  args[0],
			Author: searchAuthor,
		}

		pagination, err := paginationFromFlags(cmd, &searchPage, &searchLimit)
		if err != nil {
			return err
		}
		filter.Pagination = pagination

		return runWithRepo(func(ctx context.Context, repo *db.Repository) error {
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
	rootCmd.AddCommand(searchCmd)
}
