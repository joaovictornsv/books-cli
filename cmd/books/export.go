package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/joaovictornsv/books-cli/internal/db"
	"github.com/spf13/cobra"
)

var (
	exportFormat         string
	exportOutput         string
	exportIncludeArchived bool
)

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export the library to JSON or CSV",
	RunE: func(cmd *cobra.Command, args []string) error {
		format := strings.ToLower(strings.TrimSpace(exportFormat))
		if format == "" {
			return fmt.Errorf("--format is required")
		}
		if exportOutput == "" {
			return fmt.Errorf("--output is required")
		}
		if format != "json" && format != "csv" {
			return fmt.Errorf("invalid format %q: must be json or csv", exportFormat)
		}

		return runWithRepo(cmd.Context(), func(ctx context.Context, repo *db.Repository) error {
			books, err := repo.ListAll(ctx, exportIncludeArchived)
			if err != nil {
				return err
			}

			var writer *os.File
			outputPath := exportOutput
			if outputPath == "-" {
				writer = os.Stdout
			} else {
				writer, err = os.Create(outputPath)
				if err != nil {
					return fmt.Errorf("create output file: %w", err)
				}
				defer writer.Close()
			}

			switch format {
			case "json":
				data, err := db.MarshalBooksJSON(books)
				if err != nil {
					return err
				}
				if _, err := writer.Write(data); err != nil {
					return fmt.Errorf("write json: %w", err)
				}
			case "csv":
				if err := db.WriteBooksCSV(writer, books); err != nil {
					return err
				}
			}

			if outputPath == "-" {
				return nil
			}
			return formatter().PrintExport(cmd.OutOrStdout(), outputPath, format, len(books))
		})
	},
}

func init() {
	exportCmd.Flags().StringVar(&exportFormat, "format", "", "Output format: json or csv")
	exportCmd.Flags().StringVar(&exportOutput, "output", "", "Output file path, or - for stdout")
	exportCmd.Flags().BoolVar(&exportIncludeArchived, "include-archived", false, "Include archived books")
	rootCmd.AddCommand(exportCmd)
}
