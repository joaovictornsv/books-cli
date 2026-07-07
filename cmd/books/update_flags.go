package main

import (
	"fmt"
	"strings"

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
	updateIDs              string
)

func addUpdateFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&updateTitle, "title", "", "New title")
	cmd.Flags().StringVar(&updateAuthor, "author", "", "New author")
	cmd.Flags().StringVar(&updateCategory, "category", "", "New category")
	cmd.Flags().StringVar(&updateStatus, "status", "", "New status")
	cmd.Flags().StringVar(&updateNotes, "notes", "", "New notes")
	cmd.Flags().StringVar(&updateDescription, "description", "", "New description")
	cmd.Flags().StringVar(&updateStartedAt, "started-at", "", "Started reading at (RFC3339); pass empty string to clear")
	cmd.Flags().StringVar(&updateFinishedAt, "finished-at", "", "Finished reading at (RFC3339); pass empty string to clear")
	cmd.Flags().BoolVar(&updatePriority, "priority", false, "Set priority to buy")
	cmd.Flags().BoolVar(&updateNoPriority, "no-priority", false, "Clear priority to buy")
	cmd.Flags().BoolVar(&updateEligibleToSell, "eligible-to-sell", false, "Set eligible to sell")
	cmd.Flags().BoolVar(&updateNoEligibleToSell, "no-eligible-to-sell", false, "Clear eligible to sell")
	cmd.Flags().BoolVar(&updateSold, "sold", false, "Mark as sold")
	cmd.Flags().BoolVar(&updateNoSold, "no-sold", false, "Clear sold flag")
	cmd.Flags().StringVar(&updateIDs, "ids", "", "Comma-separated book IDs to update in bulk")
}

func buildUpdatePatch(cmd *cobra.Command) (models.BookPatch, error) {
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
				return models.BookPatch{}, err
			}
			patch.Category = &category
		}
	}
	if flags.Changed("status") {
		status, err := models.ParseStatus(updateStatus)
		if err != nil {
			return models.BookPatch{}, err
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
			return models.BookPatch{}, fmt.Errorf("cannot use --priority and --no-priority together")
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
			return models.BookPatch{}, fmt.Errorf("cannot use --eligible-to-sell and --no-eligible-to-sell together")
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
			return models.BookPatch{}, fmt.Errorf("cannot use --sold and --no-sold together")
		}
		v := models.ToBool01(updateSold)
		patch.Sold = &v
	}
	if flags.Changed("no-sold") {
		v := 0
		patch.Sold = &v
	}

	if isPatchEmpty(patch) {
		return models.BookPatch{}, fmt.Errorf("no fields to update: pass at least one flag")
	}
	return patch, nil
}

func isPatchEmpty(patch models.BookPatch) bool {
	return patch.Title == nil && patch.Author == nil && !patch.ClearAuthor &&
		patch.Category == nil && !patch.ClearCategory && patch.Status == nil &&
		patch.Notes == nil && patch.Description == nil && patch.StartedAt == nil &&
		!patch.ClearStartedAt && patch.FinishedAt == nil && !patch.ClearFinishedAt &&
		patch.PriorityToBuy == nil && patch.EligibleToSell == nil && patch.Sold == nil
}

func parseIDsList(raw string) ([]int64, error) {
	parts := strings.Split(raw, ",")
	ids := make([]int64, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		id, err := parseID(part)
		if err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	if len(ids) == 0 {
		return nil, fmt.Errorf("invalid --ids %q: must contain at least one positive integer", raw)
	}
	return ids, nil
}

func resolveUpdateTargets(cmd *cobra.Command, args []string) ([]int64, error) {
	hasIDs := cmd.Flags().Changed("ids")
	hasPositional := len(args) > 0

	switch {
	case hasIDs && hasPositional:
		return nil, fmt.Errorf("cannot use positional id and --ids together")
	case !hasIDs && !hasPositional:
		return nil, fmt.Errorf("provide a book id or --ids")
	case hasIDs:
		return parseIDsList(updateIDs)
	default:
		id, err := parseID(args[0])
		if err != nil {
			return nil, err
		}
		return []int64{id}, nil
	}
}
