package main

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/joaovictornsv/books-cli/internal/config"
	"github.com/joaovictornsv/books-cli/internal/db"
	"github.com/joaovictornsv/books-cli/internal/output"
	"github.com/spf13/cobra"
)

var (
	jsonOutput bool
	version    = "0.1.0"
)

var rootCmd = &cobra.Command{
	Use:           "books",
	Short:         "Manage a personal reading list",
	SilenceUsage:  true,
	SilenceErrors: true,
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&jsonOutput, "json", false, "Machine-readable JSON output")
	rootCmd.Version = version
}

func openRepo() (*db.Repository, func(), error) {
	cfg, err := config.Resolve()
	if err != nil {
		return nil, nil, err
	}

	database, err := db.Open(cfg.DatabasePath)
	if err != nil {
		return nil, nil, err
	}

	cleanup := func() {
		_ = database.Close()
	}
	return db.NewRepository(database), cleanup, nil
}

func formatter() output.Formatter {
	return output.New(jsonOutput)
}

func parseID(arg string) (int64, error) {
	id, err := strconv.ParseInt(arg, 10, 64)
	if err != nil || id <= 0 {
		return 0, fmt.Errorf("invalid id %q: must be a positive integer", arg)
	}
	return id, nil
}

func handleRepoError(err error) error {
	if errors.Is(err, db.ErrNotFound) {
		return fmt.Errorf("book not found")
	}
	return err
}

func runWithRepo(ctx context.Context, fn func(context.Context, *db.Repository) error) error {
	repo, cleanup, err := openRepo()
	if err != nil {
		return err
	}
	defer cleanup()
	return fn(ctx, repo)
}
