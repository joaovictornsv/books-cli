package main

import (
	"fmt"

	"github.com/joaovictornsv/books-cli/internal/db"
	"github.com/joaovictornsv/books-cli/internal/models"
	"github.com/joaovictornsv/books-cli/internal/output"
	"github.com/spf13/cobra"
)

func paginationFromFlags(cmd *cobra.Command, page, limit *int) (*models.Pagination, error) {
	pagination := models.Pagination{Page: *page, Limit: *limit}
	if err := pagination.Validate(); err != nil {
		return nil, err
	}
	return &pagination, nil
}

func addPaginationFlags(cmd *cobra.Command, page, limit *int) {
	cmd.Flags().IntVar(page, "page", 1, "Page number (1-based)")
	cmd.Flags().IntVar(limit, "limit", models.DefaultPageLimit, fmt.Sprintf("Results per page (max %d)", models.MaxPageLimit))
}

func printBooksResult(cmd *cobra.Command, result db.BooksResult, pagination *models.Pagination) error {
	return formatter().PrintBooks(cmd.OutOrStdout(), output.BooksPage{
		Books:      result.Books,
		Total:      result.Total,
		Pagination: pagination,
	})
}
