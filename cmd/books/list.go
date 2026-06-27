package main

import (
	"context"

	"github.com/joaovictornsv/books-cli/internal/db"
	"github.com/joaovictornsv/books-cli/internal/models"
	"github.com/spf13/cobra"
)

var (
	listStatus         string
	listPriority       bool
	listEligibleToSell bool
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List books with optional filters",
	RunE: func(cmd *cobra.Command, args []string) error {
		filter := db.ListFilter{}

		if cmd.Flags().Changed("status") {
			status, err := models.ParseStatus(listStatus)
			if err != nil {
				return err
			}
			filter.Status = &status
		}
		if listPriority {
			filter.PriorityToBuy = &listPriority
		}
		if listEligibleToSell {
			filter.EligibleToSell = &listEligibleToSell
		}

		return runWithRepo(func(ctx context.Context, repo *db.Repository) error {
			books, err := repo.List(ctx, filter)
			if err != nil {
				return err
			}
			return formatter().PrintBooks(cmd.OutOrStdout(), books)
		})
	},
}

func init() {
	listCmd.Flags().StringVar(&listStatus, "status", "", "Filter by status")
	listCmd.Flags().BoolVar(&listPriority, "priority", false, "Only priority-to-buy books")
	listCmd.Flags().BoolVar(&listEligibleToSell, "eligible-to-sell", false, "Only eligible-to-sell books")
	rootCmd.AddCommand(listCmd)
}
