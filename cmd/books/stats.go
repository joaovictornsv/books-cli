package main

import (
	"context"
	"fmt"
	"time"

	"github.com/joaovictornsv/books-cli/internal/db"
	"github.com/spf13/cobra"
)

var statsYear int

var statsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Show library aggregates",
	RunE: func(cmd *cobra.Command, args []string) error {
		year := statsYear
		if !cmd.Flags().Changed("year") {
			year = time.Now().UTC().Year()
		}
		if year < 1 {
			return fmt.Errorf("invalid year %d", year)
		}

		return runWithRepo(cmd.Context(), func(ctx context.Context, repo *db.Repository) error {
			result, err := repo.Stats(ctx, year)
			if err != nil {
				return err
			}
			return formatter().PrintStats(cmd.OutOrStdout(), result)
		})
	},
}

func init() {
	statsCmd.Flags().IntVar(&statsYear, "year", 0, "Year for finished_this_year (default: current year)")
	rootCmd.AddCommand(statsCmd)
}
