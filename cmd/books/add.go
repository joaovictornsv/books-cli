package main

import (
	"context"

	"github.com/joaovictornsv/books-cli/internal/db"
	"github.com/joaovictornsv/books-cli/internal/models"
	"github.com/spf13/cobra"
)

var (
	addAuthor         string
	addStatus         string
	addPriority       bool
	addEligibleToSell bool
	addNotes          string
)

var addCmd = &cobra.Command{
	Use:   "add [title]",
	Short: "Add a book",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		status, err := models.ParseStatus(addStatus)
		if err != nil {
			return err
		}

		book := models.Book{
			Title:          args[0],
			Status:         status,
			PriorityToBuy:  models.ToBool01(addPriority),
			EligibleToSell: models.ToBool01(addEligibleToSell),
			Sold:           0,
		}
		if addAuthor != "" {
			book.Author = &addAuthor
		}
		if addNotes != "" {
			book.Notes = &addNotes
		}

		return runWithRepo(func(ctx context.Context, repo *db.Repository) error {
			created, err := repo.Create(ctx, book)
			if err != nil {
				return err
			}
			return formatter().PrintBook(cmd.OutOrStdout(), created)
		})
	},
}

func init() {
	addCmd.Flags().StringVar(&addAuthor, "author", "", "Book author")
	addCmd.Flags().StringVar(&addStatus, "status", models.StatusNotStarted.String(), "Book status")
	addCmd.Flags().BoolVar(&addPriority, "priority", false, "Mark as priority to buy")
	addCmd.Flags().BoolVar(&addEligibleToSell, "eligible-to-sell", false, "Mark as eligible to sell")
	addCmd.Flags().StringVar(&addNotes, "notes", "", "Free-form notes")
	rootCmd.AddCommand(addCmd)
}
