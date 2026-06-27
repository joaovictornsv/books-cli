package main

import (
	"github.com/joaovictornsv/books-cli/internal/db"
	"github.com/joaovictornsv/books-cli/internal/models"
	"github.com/joaovictornsv/books-cli/internal/output"
	"github.com/spf13/cobra"
)

func paginationFromFlags(cmd *cobra.Command, page, limit *int) (*models.Pagination, error) {
	if !cmd.Flags().Changed("page") && !cmd.Flags().Changed("limit") {
		return nil, nil
	}

	p := *page
	l := *limit
	if !cmd.Flags().Changed("page") {
		p = 1
	}
	if !cmd.Flags().Changed("limit") {
		l = models.DefaultPageLimit
	}

	pagination := models.Pagination{Page: p, Limit: l}
	if err := pagination.Validate(); err != nil {
		return nil, err
	}
	return &pagination, nil
}

func addPaginationFlags(cmd *cobra.Command, page, limit *int) {
	cmd.Flags().IntVar(page, "page", 1, "Page number (1-based)")
	cmd.Flags().IntVar(limit, "limit", models.DefaultPageLimit, "Results per page")
}

func printBooksResult(cmd *cobra.Command, result db.BooksResult, pagination *models.Pagination) error {
	return formatter().PrintBooks(cmd.OutOrStdout(), output.BooksPage{
		Books:      result.Books,
		Total:      result.Total,
		Pagination: pagination,
	})
}
