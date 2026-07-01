package main

import (
	"github.com/joaovictornsv/books-cli/internal/db"
	"github.com/joaovictornsv/books-cli/internal/models"
	"github.com/spf13/cobra"
)

func addListFilterFlags(cmd *cobra.Command, status, category *string, priority, eligibleToSell *bool) {
	cmd.Flags().StringVar(status, "status", "", "Filter by status")
	if category != nil {
		cmd.Flags().StringVar(category, "category", "", "Filter by category")
	}
	cmd.Flags().BoolVar(priority, "priority", false, "Only priority-to-buy books")
	cmd.Flags().BoolVar(eligibleToSell, "eligible-to-sell", false, "Only eligible-to-sell books")
}

func listFilterFromFlags(cmd *cobra.Command, status *string, category *string, priority, eligibleToSell *bool) (db.ListFilter, error) {
	filter := db.ListFilter{}

	if cmd.Flags().Changed("status") {
		parsed, err := models.ParseStatus(*status)
		if err != nil {
			return db.ListFilter{}, err
		}
		filter.Status = &parsed
	}
	if category != nil && cmd.Flags().Changed("category") {
		parsed, err := models.ParseCategory(*category)
		if err != nil {
			return db.ListFilter{}, err
		}
		filter.Category = &parsed
	}
	if cmd.Flags().Changed("priority") {
		filter.PriorityToBuy = priority
	}
	if cmd.Flags().Changed("eligible-to-sell") {
		filter.EligibleToSell = eligibleToSell
	}
	return filter, nil
}
