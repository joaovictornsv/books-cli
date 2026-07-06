package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/joaovictornsv/books-cli/internal/db"
	"github.com/spf13/cobra"
)

var (
	getTitle  string
	getAuthor string
	getExact  bool
)

var getCmd = &cobra.Command{
	Use:   "get [id]",
	Short: "Show one book by ID or title",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		title := strings.TrimSpace(getTitle)
		hasTitle := title != ""
		hasID := len(args) == 1

		switch {
		case hasID && hasTitle:
			return fmt.Errorf("provide either an id or --title, not both")
		case !hasID && !hasTitle:
			return fmt.Errorf("id or --title is required")
		}

		return runWithRepo(cmd.Context(), func(ctx context.Context, repo *db.Repository) error {
			if hasID {
				id, err := parseID(args[0])
				if err != nil {
					return err
				}
				book, err := repo.GetByID(ctx, id)
				if err != nil {
					return handleRepoError(err)
				}
				return formatter().PrintBook(cmd.OutOrStdout(), book)
			}

			book, err := repo.GetByTitle(ctx, db.TitleFilter{
				Title:  title,
				Author: getAuthor,
				Exact:  getExact,
			})
			if err != nil {
				return handleRepoError(err)
			}
			return formatter().PrintBook(cmd.OutOrStdout(), book)
		})
	},
}

func init() {
	getCmd.Flags().StringVar(&getTitle, "title", "", "Look up book by title substring")
	getCmd.Flags().StringVar(&getAuthor, "author", "", "Filter by author substring when using --title")
	getCmd.Flags().BoolVar(&getExact, "exact", false, "Case-insensitive exact title match when using --title")
	rootCmd.AddCommand(getCmd)
}
