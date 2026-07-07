package main

import (
	"bytes"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestUpdateBooleanFlagsMutualExclusion(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{name: "priority and no-priority", args: []string{"update", "1", "--priority", "--no-priority"}},
		{name: "eligible and no-eligible", args: []string{"update", "1", "--eligible-to-sell", "--no-eligible-to-sell"}},
		{name: "sold and no-sold", args: []string{"update", "1", "--sold", "--no-sold"}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			resetUpdateCommandFlags(t)
			var buf bytes.Buffer
			rootCmd.SetOut(&buf)
			rootCmd.SetErr(&buf)
			rootCmd.SetArgs(tc.args)

			err := rootCmd.Execute()
			if err == nil {
				t.Fatal("expected error for conflicting boolean flags")
			}
			if !strings.Contains(err.Error(), "cannot use") {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestUpdateNoPriorityFlag(t *testing.T) {
	cmd := &cobra.Command{}
	cmd.Flags().Bool("no-priority", false, "")
	if err := cmd.Flags().Set("no-priority", "true"); err != nil {
		t.Fatal(err)
	}
	if !cmd.Flags().Changed("no-priority") {
		t.Fatal("expected no-priority to be changed")
	}
}
