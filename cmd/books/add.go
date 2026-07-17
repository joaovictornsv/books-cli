package main

import (
	"context"

	"github.com/joaovictornsv/books-cli/internal/db"
	"github.com/joaovictornsv/books-cli/internal/models"
	"github.com/spf13/cobra"
)

var (
	addAuthor           string
	addCategory         string
	addStatus           string
	addPriority         bool
	addEligibleToDonate bool
	addNotes            string
	addDescription      string
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
			Title:            args[0],
			Status:           status,
			PriorityToBuy:    models.ToBool01(addPriority),
			EligibleToDonate: models.ToBool01(addEligibleToDonate),
			Donated:          0,
		}
		if addAuthor != "" {
			book.Author = &addAuthor
		}
		if addCategory != "" {
			category, err := models.ParseCategory(addCategory)
			if err != nil {
				return err
			}
			book.Category = &category
		}
		if addNotes != "" {
			book.Notes = &addNotes
		}
		if addDescription != "" {
			book.Description = &addDescription
		}
		if err := book.ValidateForCreate(); err != nil {
			return err
		}

		return runWithRepo(cmd.Context(), func(ctx context.Context, repo *db.Repository) error {
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
	addCmd.Flags().StringVar(&addCategory, "category", "", "Book category")
	addCmd.Flags().StringVar(&addStatus, "status", models.StatusNotStarted.String(), "Book status")
	addCmd.Flags().BoolVar(&addPriority, "priority", false, "Mark as priority to buy")
	addCmd.Flags().BoolVar(&addEligibleToDonate, "eligible-to-donate", false, "Mark as eligible to donate")
	addCmd.Flags().StringVar(&addNotes, "notes", "", "Free-form notes")
	addCmd.Flags().StringVar(&addDescription, "description", "", "Book description (e.g. from the web)")
	rootCmd.AddCommand(addCmd)
}
