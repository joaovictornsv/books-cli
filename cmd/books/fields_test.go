package main

import (
	"testing"

	"github.com/joaovictornsv/books-cli/internal/db"
	"github.com/spf13/cobra"
)

func TestFieldsFromFlags(t *testing.T) {
	t.Run("unchanged flag returns nil", func(t *testing.T) {
		cmd := &cobra.Command{}
		addFieldsFlag(cmd)

		got, err := fieldsFromFlags(cmd)
		if err != nil {
			t.Fatal(err)
		}
		if got != nil {
			t.Fatalf("got %v, want nil", got)
		}
	})

	t.Run("parses valid fields", func(t *testing.T) {
		cmd := &cobra.Command{}
		addFieldsFlag(cmd)
		if err := cmd.Flags().Set(fieldsFlagName, "id,title"); err != nil {
			t.Fatal(err)
		}

		got, err := fieldsFromFlags(cmd)
		if err != nil {
			t.Fatal(err)
		}
		want := []string{"id", "title"}
		if len(got) != len(want) {
			t.Fatalf("got %v, want %v", got, want)
		}
		for i := range want {
			if got[i] != want[i] {
				t.Fatalf("got %v, want %v", got, want)
			}
		}
	})

	t.Run("invalid field returns error", func(t *testing.T) {
		cmd := &cobra.Command{}
		addFieldsFlag(cmd)
		if err := cmd.Flags().Set(fieldsFlagName, "nope"); err != nil {
			t.Fatal(err)
		}

		_, err := fieldsFromFlags(cmd)
		if err == nil {
			t.Fatal("expected validation error")
		}
	})
}

func TestPrintBooksResultRejectsFieldsWithoutJSON(t *testing.T) {
	jsonOutput = false
	t.Cleanup(func() {
		jsonOutput = false
	})

	cmd := &cobra.Command{}
	addFieldsFlag(cmd)
	if err := cmd.Flags().Set(fieldsFlagName, "id,title"); err != nil {
		t.Fatal(err)
	}

	err := printBooksResult(cmd, db.BooksResult{}, nil)
	if err == nil {
		t.Fatal("expected error when --fields is used without --json")
	}
}
