package main

import (
	"context"

	"github.com/joaovictornsv/books-cli/internal/db"
	"github.com/spf13/cobra"
)

var (
	listStatus         string
	listCategory       string
	listPriority       bool
	listEligibleToDonate bool
	listPage           int
	listLimit          int
	listSort           string
	listOrder          string
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List books with optional filters",
	RunE: func(cmd *cobra.Command, args []string) error {
		filter, err := listFilterFromFlags(cmd, &listStatus, &listCategory, &listPriority, &listEligibleToDonate)
		if err != nil {
			return err
		}

		pagination, err := paginationFromFlags(cmd, &listPage, &listLimit)
		if err != nil {
			return err
		}
		filter.Pagination = pagination

		sort, err := sortFromFlags(&listSort, &listOrder)
		if err != nil {
			return err
		}
		filter.Sort = sort

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
	addListFilterFlags(listCmd, &listStatus, &listCategory, &listPriority, &listEligibleToDonate)
	addPaginationFlags(listCmd, &listPage, &listLimit)
	addSortFlags(listCmd, &listSort, &listOrder)
	addFieldsFlag(listCmd)
	rootCmd.AddCommand(listCmd)
}
