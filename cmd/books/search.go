package main

import (
	"context"

	"github.com/joaovictornsv/books-cli/internal/db"
	"github.com/spf13/cobra"
)

var searchAuthor string

var searchCmd = &cobra.Command{
	Use:   "search [query]",
	Short: "Search books by title and optional author",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		filter := db.SearchFilter{
			Query:  args[0],
			Author: searchAuthor,
		}

		return runWithRepo(func(ctx context.Context, repo *db.Repository) error {
			books, err := repo.Search(ctx, filter)
			if err != nil {
				return err
			}
			return formatter().PrintBooks(cmd.OutOrStdout(), books)
		})
	},
}

func init() {
	searchCmd.Flags().StringVar(&searchAuthor, "author", "", "Filter by author substring")
	rootCmd.AddCommand(searchCmd)
}
