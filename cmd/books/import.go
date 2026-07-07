package main

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/joaovictornsv/books-cli/internal/db"
	"github.com/joaovictornsv/books-cli/internal/models"
	"github.com/spf13/cobra"
)

var (
	importInput  string
	importDryRun bool
)

var importCmd = &cobra.Command{
	Use:   "import",
	Short: "Import books from JSON or CSV",
	RunE: func(cmd *cobra.Command, args []string) error {
		if importInput == "" {
			return fmt.Errorf("--input is required")
		}

		data, err := os.ReadFile(importInput)
		if err != nil {
			return fmt.Errorf("read input file: %w", err)
		}

		books, err := parseImportBooks(importInput, data)
		if err != nil {
			return err
		}

		return runWithRepo(cmd.Context(), func(ctx context.Context, repo *db.Repository) error {
			result, err := repo.ImportBooks(ctx, books, importDryRun)
			if err != nil {
				return err
			}
			if len(result.Errors) > 0 {
				return fmt.Errorf("import validation failed:\n%s", strings.Join(result.Errors, "\n"))
			}
			return formatter().PrintImport(cmd.OutOrStdout(), result)
		})
	},
}

func parseImportBooks(path string, data []byte) ([]models.Book, error) {
	switch strings.ToLower(filepath.Ext(path)) {
	case ".json":
		return db.UnmarshalBooksJSON(data)
	case ".csv":
		return db.ReadBooksCSV(bytes.NewReader(data))
	default:
		return nil, fmt.Errorf("unsupported input format %q: use .json or .csv", path)
	}
}

func init() {
	importCmd.Flags().StringVar(&importInput, "input", "", "Input file path (required)")
	importCmd.Flags().BoolVar(&importDryRun, "dry-run", false, "Validate without writing to the database")
	rootCmd.AddCommand(importCmd)
}
