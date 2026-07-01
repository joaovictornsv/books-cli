package main

import (
	"testing"

	"github.com/joaovictornsv/books-cli/internal/models"
	"github.com/spf13/cobra"
)

type listFilterFixture struct {
	cmd              *cobra.Command
	status, category string
	priority         bool
	eligible         bool
}

func newListFilterCmd(t *testing.T) *listFilterFixture {
	t.Helper()
	f := &listFilterFixture{cmd: &cobra.Command{}}
	addListFilterFlags(f.cmd, &f.status, &f.category, &f.priority, &f.eligible)
	return f
}

func TestListFilterFromFlags(t *testing.T) {
	t.Run("unchanged flags produce empty filter", func(t *testing.T) {
		f := newListFilterCmd(t)

		got, err := listFilterFromFlags(f.cmd, &f.status, &f.category, &f.priority, &f.eligible)
		if err != nil {
			t.Fatal(err)
		}
		if got.Status != nil || got.Category != nil || got.PriorityToBuy != nil || got.EligibleToSell != nil {
			t.Fatalf("expected empty filter, got %+v", got)
		}
	})

	t.Run("status filter", func(t *testing.T) {
		f := newListFilterCmd(t)
		if err := f.cmd.Flags().Set("status", "READ"); err != nil {
			t.Fatal(err)
		}

		got, err := listFilterFromFlags(f.cmd, &f.status, &f.category, &f.priority, &f.eligible)
		if err != nil {
			t.Fatal(err)
		}
		if got.Status == nil || *got.Status != models.StatusRead {
			t.Fatalf("expected READ status, got %+v", got.Status)
		}
	})

	t.Run("invalid status", func(t *testing.T) {
		f := newListFilterCmd(t)
		if err := f.cmd.Flags().Set("status", "INVALID"); err != nil {
			t.Fatal(err)
		}

		_, err := listFilterFromFlags(f.cmd, &f.status, &f.category, &f.priority, &f.eligible)
		if err == nil {
			t.Fatal("expected error for invalid status")
		}
	})

	t.Run("category filter", func(t *testing.T) {
		f := newListFilterCmd(t)
		if err := f.cmd.Flags().Set("category", "FICTION"); err != nil {
			t.Fatal(err)
		}

		got, err := listFilterFromFlags(f.cmd, &f.status, &f.category, &f.priority, &f.eligible)
		if err != nil {
			t.Fatal(err)
		}
		if got.Category == nil || *got.Category != models.CategoryFiction {
			t.Fatalf("expected FICTION category, got %+v", got.Category)
		}
	})

	t.Run("nil category pointer skips category filter", func(t *testing.T) {
		cmd := &cobra.Command{}
		status, priority, eligible := "", false, false
		addListFilterFlags(cmd, &status, nil, &priority, &eligible)

		got, err := listFilterFromFlags(cmd, &status, nil, &priority, &eligible)
		if err != nil {
			t.Fatal(err)
		}
		if got.Category != nil {
			t.Fatalf("expected category to be ignored, got %+v", got.Category)
		}
	})

	t.Run("priority true", func(t *testing.T) {
		f := newListFilterCmd(t)
		if err := f.cmd.Flags().Set("priority", "true"); err != nil {
			t.Fatal(err)
		}

		got, err := listFilterFromFlags(f.cmd, &f.status, &f.category, &f.priority, &f.eligible)
		if err != nil {
			t.Fatal(err)
		}
		if got.PriorityToBuy == nil || !*got.PriorityToBuy {
			t.Fatalf("expected priority true, got %+v", got.PriorityToBuy)
		}
	})

	t.Run("priority false", func(t *testing.T) {
		f := newListFilterCmd(t)
		if err := f.cmd.Flags().Set("priority", "false"); err != nil {
			t.Fatal(err)
		}

		got, err := listFilterFromFlags(f.cmd, &f.status, &f.category, &f.priority, &f.eligible)
		if err != nil {
			t.Fatal(err)
		}
		if got.PriorityToBuy == nil || *got.PriorityToBuy {
			t.Fatalf("expected priority false, got %+v", got.PriorityToBuy)
		}
	})

	t.Run("eligible to sell", func(t *testing.T) {
		f := newListFilterCmd(t)
		if err := f.cmd.Flags().Set("eligible-to-sell", "true"); err != nil {
			t.Fatal(err)
		}

		got, err := listFilterFromFlags(f.cmd, &f.status, &f.category, &f.priority, &f.eligible)
		if err != nil {
			t.Fatal(err)
		}
		if got.EligibleToSell == nil || !*got.EligibleToSell {
			t.Fatalf("expected eligible to sell, got %+v", got.EligibleToSell)
		}
	})
}
