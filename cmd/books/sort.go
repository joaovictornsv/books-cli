package main

import (
	"github.com/joaovictornsv/books-cli/internal/models"
	"github.com/spf13/cobra"
)

func addSortFlags(cmd *cobra.Command, sort, order *string) {
	cmd.Flags().StringVar(sort, "sort", string(models.SortFieldID), "Sort field (id, title, author, status, added_at, started_at, finished_at)")
	cmd.Flags().StringVar(order, "order", string(models.SortOrderAsc), "Sort order (asc or desc)")
}

func sortFromFlags(sort, order *string) (models.Sort, error) {
	field, err := models.ParseSortField(*sort)
	if err != nil {
		return models.Sort{}, err
	}
	sortOrder, err := models.ParseSortOrder(*order)
	if err != nil {
		return models.Sort{}, err
	}
	return models.Sort{Field: field, Order: sortOrder}, nil
}
