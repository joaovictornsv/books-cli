package main

import (
	"context"

	"github.com/joaovictornsv/books-cli/internal/db"
	"github.com/spf13/cobra"
)

var (
	countStatus         string
	countCategory       string
	countPriority       bool
	countEligibleToSell bool
)

var countCmd = &cobra.Command{
	Use:   "count",
	Short: "Count books matching optional filters",
	RunE: func(cmd *cobra.Command, args []string) error {
		filter, err := listFilterFromFlags(cmd, &countStatus, &countCategory, &countPriority, &countEligibleToSell)
		if err != nil {
			return err
		}

		return runWithRepo(cmd.Context(), func(ctx context.Context, repo *db.Repository) error {
			total, err := repo.Count(ctx, filter)
			if err != nil {
				return err
			}
			return formatter().PrintCount(cmd.OutOrStdout(), total)
		})
	},
}

func init() {
	addListFilterFlags(countCmd, &countStatus, &countCategory, &countPriority, &countEligibleToSell)
	rootCmd.AddCommand(countCmd)
}
