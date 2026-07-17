package main

import (
	"github.com/joaovictornsv/books-cli/internal/db"
	"github.com/joaovictornsv/books-cli/internal/models"
	"github.com/spf13/cobra"
)

func addListFilterFlags(cmd *cobra.Command, status, category *string, priority, eligibleToDonate *bool) {
	cmd.Flags().StringVar(status, "status", "", "Filter by status")
	if category != nil {
		cmd.Flags().StringVar(category, "category", "", "Filter by category")
	}
	cmd.Flags().BoolVar(priority, "priority", false, "Only priority-to-buy books")
	cmd.Flags().BoolVar(eligibleToDonate, "eligible-to-donate", false, "Only eligible-to-donate books")
}

func listFilterFromFlags(cmd *cobra.Command, status *string, category *string, priority, eligibleToDonate *bool) (db.ListFilter, error) {
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
	if cmd.Flags().Changed("eligible-to-donate") {
		filter.EligibleToDonate = eligibleToDonate
	}
	return filter, nil
}

func categoryFromFlag(cmd *cobra.Command, category *string) (*models.Category, error) {
	if !cmd.Flags().Changed("category") {
		return nil, nil
	}
	parsed, err := models.ParseCategory(*category)
	if err != nil {
		return nil, err
	}
	return &parsed, nil
}
