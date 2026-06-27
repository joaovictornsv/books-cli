package main

import (
	"github.com/joaovictornsv/books-cli/internal/config"
	"github.com/joaovictornsv/books-cli/internal/output"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Show effective CLI configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Resolve()
		if err != nil {
			return err
		}
		return output.PrintConfigHuman(cmd.OutOrStdout(), cfg)
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
}
