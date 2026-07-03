package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/joaovictornsv/books-cli/internal/db"
	"github.com/spf13/cobra"
)

var deleteYes bool

var deleteCmd = &cobra.Command{
	Use:   "delete [id]",
	Short: "Permanently remove a book from the database",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := parseID(args[0])
		if err != nil {
			return err
		}

		return runWithRepo(cmd.Context(), func(ctx context.Context, repo *db.Repository) error {
			if jsonOutput && !deleteYes {
				return fmt.Errorf("delete requires --yes when using --json")
			}

			book, err := repo.GetByID(ctx, id)
			if err != nil {
				return handleRepoError(err)
			}

			if !deleteYes && isInteractiveTerminal() {
				prompt := fmt.Sprintf(`Delete "%s" (id %d)? [y/N] `, book.Title, book.ID)
				if _, err := fmt.Fprint(cmd.OutOrStdout(), prompt); err != nil {
					return err
				}
				reader := bufio.NewReader(os.Stdin)
				line, err := reader.ReadString('\n')
				if err != nil {
					return fmt.Errorf("read confirmation: %w", err)
				}
				answer := strings.ToLower(strings.TrimSpace(line))
				if answer != "y" && answer != "yes" {
					return fmt.Errorf("delete cancelled")
				}
			}

			deleted, err := repo.Delete(ctx, id)
			if err != nil {
				return handleRepoError(err)
			}
			return formatter().PrintBook(cmd.OutOrStdout(), deleted)
		})
	},
}

func isInteractiveTerminal() bool {
	fi, err := os.Stdin.Stat()
	if err != nil {
		return false
	}
	return (fi.Mode() & os.ModeCharDevice) != 0
}

func init() {
	deleteCmd.Flags().BoolVarP(&deleteYes, "yes", "y", false, "Confirm deletion without prompting")
	rootCmd.AddCommand(deleteCmd)
}
