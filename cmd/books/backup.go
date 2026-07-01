package main

import (
	"fmt"
	"path/filepath"

	"github.com/joaovictornsv/books-cli/internal/config"
	"github.com/joaovictornsv/books-cli/internal/db"
	"github.com/spf13/cobra"
)

var (
	backupOutput string
	backupForce  bool
)

var backupCmd = &cobra.Command{
	Use:   "backup",
	Short: "Copy the SQLite database file",
	RunE: func(cmd *cobra.Command, args []string) error {
		if backupOutput == "" {
			return fmt.Errorf("--output is required")
		}

		cfg, err := config.Resolve()
		if err != nil {
			return err
		}

		source := cfg.DatabasePath
		dest := filepath.Clean(backupOutput)

		if err := db.Backup(cmd.Context(), source, dest, backupForce); err != nil {
			return err
		}

		return formatter().PrintBackup(cmd.OutOrStdout(), source, dest)
	},
}

func init() {
	backupCmd.Flags().StringVar(&backupOutput, "output", "", "Destination file path (required)")
	backupCmd.Flags().BoolVar(&backupForce, "force", false, "Overwrite destination if it exists")
	rootCmd.AddCommand(backupCmd)
}
