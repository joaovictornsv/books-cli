package main

import (
	"context"

	"github.com/joaovictornsv/books-cli/internal/db"
	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update [id]",
	Short: "Update book fields",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ids, err := resolveUpdateTargets(cmd, args)
		if err != nil {
			return err
		}

		patch, err := buildUpdatePatch(cmd)
		if err != nil {
			return err
		}

		return runWithRepo(cmd.Context(), func(ctx context.Context, repo *db.Repository) error {
			if len(ids) == 1 {
				book, err := repo.Update(ctx, ids[0], patch)
				if err != nil {
					return handleRepoError(err)
				}
				return formatter().PrintBook(cmd.OutOrStdout(), book)
			}

			books, err := repo.UpdateMany(ctx, ids, patch)
			if err != nil {
				return handleRepoError(err)
			}
			return formatter().PrintBulkUpdate(cmd.OutOrStdout(), books)
		})
	},
}

func init() {
	addUpdateFlags(updateCmd)
	rootCmd.AddCommand(updateCmd)
}
