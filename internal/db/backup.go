package db

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func Backup(ctx context.Context, source, dest string, force bool) error {
	source = filepath.Clean(source)
	dest = filepath.Clean(dest)

	if _, err := os.Stat(source); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("database not found at %s", source)
		}
		return fmt.Errorf("stat database: %w", err)
	}

	if _, err := os.Stat(dest); err == nil {
		if !force {
			return fmt.Errorf("destination %s already exists (use --force to overwrite)", dest)
		}
		if err := os.Remove(dest); err != nil {
			return fmt.Errorf("remove existing backup: %w", err)
		}
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("stat destination: %w", err)
	}

	if dir := filepath.Dir(dest); dir != "." {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return fmt.Errorf("create backup directory: %w", err)
		}
	}

	database, err := Open(source)
	if err != nil {
		return err
	}
	defer database.Close()

	escaped := strings.ReplaceAll(dest, "'", "''")
	if _, err := database.SQL().ExecContext(ctx, fmt.Sprintf(`VACUUM INTO '%s'`, escaped)); err != nil {
		return fmt.Errorf("backup database: %w", err)
	}
	return nil
}
