package main

import (
	"context"
	"fmt"

	"github.com/joaovictornsv/books-cli/internal/db"
	"github.com/joaovictornsv/books-cli/internal/models"
	"github.com/spf13/cobra"
)

var (
	updateTitle            string
	updateAuthor           string
	updateCategory         string
	updateStatus           string
	updateNotes            string
	updateDescription      string
	updateStartedAt        string
	updateFinishedAt       string
	updatePriority         bool
	updateNoPriority       bool
	updateEligibleToSell   bool
	updateNoEligibleToSell bool
	updateSold             bool
	updateNoSold           bool
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
		if flags.Changed("started-at") {
			if updateStartedAt == "" {
				patch.ClearStartedAt = true
			} else {
				patch.StartedAt = &updateStartedAt
			}
		}
		if flags.Changed("finished-at") {
			if updateFinishedAt == "" {
				patch.ClearFinishedAt = true
			} else {
				patch.FinishedAt = &updateFinishedAt
			}
		}
		if flags.Changed("priority") {
			if flags.Changed("no-priority") {
				return fmt.Errorf("cannot use --priority and --no-priority together")
			}
			v := models.ToBool01(updatePriority)
			patch.PriorityToBuy = &v
		}
		if flags.Changed("no-priority") {
			v := 0
			patch.PriorityToBuy = &v
		}
		if flags.Changed("eligible-to-sell") {
			if flags.Changed("no-eligible-to-sell") {
				return fmt.Errorf("cannot use --eligible-to-sell and --no-eligible-to-sell together")
			}
			v := models.ToBool01(updateEligibleToSell)
			patch.EligibleToSell = &v
		}
		if flags.Changed("no-eligible-to-sell") {
			v := 0
			patch.EligibleToSell = &v
		}
		if flags.Changed("sold") {
			if flags.Changed("no-sold") {
				return fmt.Errorf("cannot use --sold and --no-sold together")
			}
			v := models.ToBool01(updateSold)
			patch.Sold = &v
		}
		if flags.Changed("no-sold") {
			v := 0
			patch.Sold = &v
		}

		if patch.Title == nil && patch.Author == nil && !patch.ClearAuthor &&
			patch.Category == nil && !patch.ClearCategory && patch.Status == nil &&
			patch.Notes == nil && patch.Description == nil && patch.StartedAt == nil &&
			!patch.ClearStartedAt && patch.FinishedAt == nil && !patch.ClearFinishedAt &&
			patch.PriorityToBuy == nil && patch.EligibleToSell == nil && patch.Sold == nil {
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
	updateCmd.Flags().StringVar(&updateStartedAt, "started-at", "", "Started reading at (RFC3339); pass empty string to clear")
	updateCmd.Flags().StringVar(&updateFinishedAt, "finished-at", "", "Finished reading at (RFC3339); pass empty string to clear")
	updateCmd.Flags().BoolVar(&updatePriority, "priority", false, "Set priority to buy")
	updateCmd.Flags().BoolVar(&updateNoPriority, "no-priority", false, "Clear priority to buy")
	updateCmd.Flags().BoolVar(&updateEligibleToSell, "eligible-to-sell", false, "Set eligible to sell")
	updateCmd.Flags().BoolVar(&updateNoEligibleToSell, "no-eligible-to-sell", false, "Clear eligible to sell")
	updateCmd.Flags().BoolVar(&updateSold, "sold", false, "Mark as sold")
	updateCmd.Flags().BoolVar(&updateNoSold, "no-sold", false, "Clear sold flag")
	rootCmd.AddCommand(updateCmd)
}
