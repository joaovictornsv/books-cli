package main

import (
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func resetUpdateCommandFlags(t *testing.T) {
	t.Helper()
	updateCmd.Flags().VisitAll(func(f *pflag.Flag) {
		_ = f.Value.Set(f.DefValue)
		f.Changed = false
	})
}

func TestParseIDsList(t *testing.T) {
	ids, err := parseIDsList("1, 2,3")
	if err != nil {
		t.Fatal(err)
	}
	if len(ids) != 3 || ids[0] != 1 || ids[1] != 2 || ids[2] != 3 {
		t.Fatalf("unexpected ids: %v", ids)
	}

	_, err = parseIDsList(" ")
	if err == nil {
		t.Fatal("expected error for empty ids")
	}
}

func TestResolveUpdateTargetsMutualExclusion(t *testing.T) {
	cmd := &cobra.Command{}
	addUpdateFlags(cmd)
	if err := cmd.Flags().Set("ids", "1,2"); err != nil {
		t.Fatal(err)
	}

	_, err := resolveUpdateTargets(cmd, []string{"3"})
	if err == nil {
		t.Fatal("expected error when using positional id and --ids")
	}
	if !strings.Contains(err.Error(), "cannot use positional id and --ids together") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestBuildUpdatePatchRequiresFields(t *testing.T) {
	cmd := &cobra.Command{}
	addUpdateFlags(cmd)

	_, err := buildUpdatePatch(cmd)
	if err == nil {
		t.Fatal("expected error when no update flags provided")
	}
}
