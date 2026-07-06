package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/joaovictornsv/books-cli/internal/db"
	"github.com/spf13/cobra"
)

var (
	searchTerms    []string
	searchAuthor   string
	searchCategory string
	searchPage     int
	searchLimit    int
	searchSort     string
	searchOrder    string
)

var searchCmd = &cobra.Command{
	Use:   "search [query]",
	Short: "Search books by title, description, author, and optional filters",
	Long: `Search books by case-insensitive substring in title, description, or author.

Pass a positional query and/or repeatable --term flags. Multiple terms are combined
with OR (a book matches if any term hits title, description, or author).`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		terms, err := collectSearchTerms(args, searchTerms)
		if err != nil {
			return err
		}

		filter := db.SearchFilter{
			Terms:  terms,
			Author: searchAuthor,
		}

		category, err := categoryFromFlag(cmd, &searchCategory)
		if err != nil {
			return err
		}
		filter.Category = category

		pagination, err := paginationFromFlags(cmd, &searchPage, &searchLimit)
		if err != nil {
			return err
		}
		filter.Pagination = pagination

		sort, err := sortFromFlags(cmd, &searchSort, &searchOrder)
		if err != nil {
			return err
		}
		filter.Sort = sort

		return runWithRepo(cmd.Context(), func(ctx context.Context, repo *db.Repository) error {
			result, err := repo.Search(ctx, filter)
			if err != nil {
				return err
			}
			return printBooksResult(cmd, result, pagination)
		})
	},
}

func collectSearchTerms(args, flagTerms []string) ([]string, error) {
	terms := make([]string, 0, len(args)+len(flagTerms))
	for _, arg := range args {
		if term := strings.TrimSpace(arg); term != "" {
			terms = append(terms, term)
		}
	}
	for _, term := range flagTerms {
		if term = strings.TrimSpace(term); term != "" {
			terms = append(terms, term)
		}
	}
	if len(terms) == 0 {
		return nil, fmt.Errorf("at least one search term is required (positional query or --term)")
	}
	return terms, nil
}

func init() {
	searchCmd.Flags().StringArrayVar(&searchTerms, "term", nil, "Search term substring (repeatable; terms are OR'd)")
	searchCmd.Flags().StringVar(&searchAuthor, "author", "", "Filter by author substring")
	searchCmd.Flags().StringVar(&searchCategory, "category", "", "Filter by category")
	addPaginationFlags(searchCmd, &searchPage, &searchLimit)
	addSortFlags(searchCmd, &searchSort, &searchOrder)
	addFieldsFlag(searchCmd)
	rootCmd.AddCommand(searchCmd)
}
