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
	listPage           int
	listLimit          int
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

		pagination, err := paginationFromFlags(cmd, &listPage, &listLimit)
		if err != nil {
			return err
		}
		filter.Pagination = pagination

		return runWithRepo(cmd.Context(), func(ctx context.Context, repo *db.Repository) error {
			result, err := repo.List(ctx, filter)
			if err != nil {
				return err
			}
			return printBooksResult(cmd, result, pagination)
		})
	},
}

func init() {
	listCmd.Flags().StringVar(&listStatus, "status", "", "Filter by status")
	listCmd.Flags().BoolVar(&listPriority, "priority", false, "Only priority-to-buy books")
	listCmd.Flags().BoolVar(&listEligibleToSell, "eligible-to-sell", false, "Only eligible-to-sell books")
	addPaginationFlags(listCmd, &listPage, &listLimit)
	addFieldsFlag(listCmd)
	rootCmd.AddCommand(listCmd)
}
