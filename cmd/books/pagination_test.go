package main

import (
	"testing"

	"github.com/joaovictornsv/books-cli/internal/models"
	"github.com/spf13/cobra"
)

func TestPaginationFromFlags(t *testing.T) {
	t.Run("unchanged flags return nil", func(t *testing.T) {
		cmd := &cobra.Command{}
		page, limit := 1, models.DefaultPageLimit
		addPaginationFlags(cmd, &page, &limit)

		got, err := paginationFromFlags(cmd, &page, &limit)
		if err != nil {
			t.Fatal(err)
		}
		if got != nil {
			t.Fatalf("expected nil pagination, got %+v", got)
		}
	})

	t.Run("changed page uses default limit", func(t *testing.T) {
		cmd := &cobra.Command{}
		page, limit := 1, models.DefaultPageLimit
		addPaginationFlags(cmd, &page, &limit)
		if err := cmd.Flags().Set("page", "2"); err != nil {
			t.Fatal(err)
		}

		got, err := paginationFromFlags(cmd, &page, &limit)
		if err != nil {
			t.Fatal(err)
		}
		if got == nil || got.Page != 2 || got.Limit != models.DefaultPageLimit {
			t.Fatalf("got %+v, want page=2 limit=%d", got, models.DefaultPageLimit)
		}
	})

	t.Run("invalid page returns error", func(t *testing.T) {
		cmd := &cobra.Command{}
		page, limit := 1, models.DefaultPageLimit
		addPaginationFlags(cmd, &page, &limit)
		if err := cmd.Flags().Set("page", "0"); err != nil {
			t.Fatal(err)
		}

		_, err := paginationFromFlags(cmd, &page, &limit)
		if err == nil {
			t.Fatal("expected validation error for page 0")
		}
	})
}
