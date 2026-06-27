package main

import (
	"context"

	"github.com/joaovictornsv/books-cli/internal/db"
	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:   "get [id]",
	Short: "Show one book by ID",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := parseID(args[0])
		if err != nil {
			return err
		}

		return runWithRepo(func(ctx context.Context, repo *db.Repository) error {
			book, err := repo.GetByID(ctx, id)
			if err != nil {
				return handleRepoError(err)
			}
			return formatter().PrintBook(cmd.OutOrStdout(), book)
		})
	},
}

func init() {
	rootCmd.AddCommand(getCmd)
}
