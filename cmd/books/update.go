package main

import (
	"context"
	"fmt"

	"github.com/joaovictornsv/books-cli/internal/db"
	"github.com/joaovictornsv/books-cli/internal/models"
	"github.com/spf13/cobra"
)

var (
	updateTitle          string
	updateAuthor         string
	updateCategory       string
	updateStatus         string
	updateNotes          string
	updateDescription    string
	updatePriority       bool
	updateEligibleToSell bool
	updateSold           bool
)

var updateCmd = &cobra.Command{
	Use:   "update [id]",
	Short: "Update book fields",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := parseID(args[0])
		if err != nil {
			return err
		}

		patch := models.BookPatch{}
		flags := cmd.Flags()

		if flags.Changed("title") {
			patch.Title = &updateTitle
		}
		if flags.Changed("author") {
			if updateAuthor == "" {
				patch.ClearAuthor = true
			} else {
				patch.Author = &updateAuthor
			}
		}
		if flags.Changed("category") {
			if updateCategory == "" {
				patch.ClearCategory = true
			} else {
				category, err := models.ParseCategory(updateCategory)
				if err != nil {
					return err
				}
				patch.Category = &category
			}
		}
		if flags.Changed("status") {
			status, err := models.ParseStatus(updateStatus)
			if err != nil {
				return err
			}
			patch.Status = &status
		}
		if flags.Changed("notes") {
			patch.Notes = &updateNotes
		}
		if flags.Changed("description") {
			patch.Description = &updateDescription
		}
		if flags.Changed("priority") {
			v := models.ToBool01(updatePriority)
			patch.PriorityToBuy = &v
		}
		if flags.Changed("eligible-to-sell") {
			v := models.ToBool01(updateEligibleToSell)
			patch.EligibleToSell = &v
		}
		if flags.Changed("sold") {
			v := models.ToBool01(updateSold)
			patch.Sold = &v
		}

		if patch.Title == nil && patch.Author == nil && !patch.ClearAuthor &&
			patch.Category == nil && !patch.ClearCategory && patch.Status == nil &&
			patch.Notes == nil && patch.Description == nil && patch.PriorityToBuy == nil &&
			patch.EligibleToSell == nil && patch.Sold == nil {
			return fmt.Errorf("no fields to update: pass at least one flag")
		}

		return runWithRepo(cmd.Context(), func(ctx context.Context, repo *db.Repository) error {
			book, err := repo.Update(ctx, id, patch)
			if err != nil {
				return handleRepoError(err)
			}
			return formatter().PrintBook(cmd.OutOrStdout(), book)
		})
	},
}

func init() {
	updateCmd.Flags().StringVar(&updateTitle, "title", "", "New title")
	updateCmd.Flags().StringVar(&updateAuthor, "author", "", "New author")
	updateCmd.Flags().StringVar(&updateCategory, "category", "", "New category")
	updateCmd.Flags().StringVar(&updateStatus, "status", "", "New status")
	updateCmd.Flags().StringVar(&updateNotes, "notes", "", "New notes")
	updateCmd.Flags().StringVar(&updateDescription, "description", "", "New description")
	updateCmd.Flags().BoolVar(&updatePriority, "priority", false, "Set priority to buy")
	updateCmd.Flags().BoolVar(&updateEligibleToSell, "eligible-to-sell", false, "Set eligible to sell")
	updateCmd.Flags().BoolVar(&updateSold, "sold", false, "Mark as sold")
	rootCmd.AddCommand(updateCmd)
}
