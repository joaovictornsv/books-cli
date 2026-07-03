package main

import (
	"github.com/joaovictornsv/books-cli/internal/models"
	"github.com/spf13/cobra"
)

var schemaCmd = &cobra.Command{
	Use:   "schema",
	Short: "Show machine-readable enums and field semantics",
	RunE: func(cmd *cobra.Command, args []string) error {
		return formatter().PrintSchema(cmd.OutOrStdout(), models.Schema())
	},
}

func init() {
	rootCmd.AddCommand(schemaCmd)
}
