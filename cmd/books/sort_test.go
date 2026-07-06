package main

import (
	"testing"

	"github.com/joaovictornsv/books-cli/internal/models"
	"github.com/spf13/cobra"
)

func newSortCmd(t *testing.T) (*cobra.Command, *string, *string) {
	t.Helper()
	cmd := &cobra.Command{}
	sort, order := "", ""
	addSortFlags(cmd, &sort, &order)
	return cmd, &sort, &order
}

func TestSortFromFlags(t *testing.T) {
	cmd, sort, order := newSortCmd(t)

	got, err := sortFromFlags(cmd, sort, order)
	if err != nil {
		t.Fatal(err)
	}
	if got != models.DefaultSort() {
		t.Fatalf("got %+v", got)
	}

	if err := cmd.Flags().Set("sort", "title"); err != nil {
		t.Fatal(err)
	}
	if err := cmd.Flags().Set("order", "desc"); err != nil {
		t.Fatal(err)
	}

	got, err = sortFromFlags(cmd, sort, order)
	if err != nil {
		t.Fatal(err)
	}
	if got.Field != models.SortFieldTitle || got.Order != models.SortOrderDesc {
		t.Fatalf("got %+v", got)
	}
}

func TestSortFromFlagsInvalid(t *testing.T) {
	cmd, sort, order := newSortCmd(t)
	if err := cmd.Flags().Set("sort", "bad"); err != nil {
		t.Fatal(err)
	}
	_, err := sortFromFlags(cmd, sort, order)
	if err == nil {
		t.Fatal("expected error for invalid sort field")
	}
}
