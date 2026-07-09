package main

import (
	"bytes"
	"encoding/json"
	"testing"
)

func TestVersionJSONOutput(t *testing.T) {
	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetErr(&buf)
	rootCmd.SetArgs([]string{"--json", "version"})

	if err := rootCmd.Execute(); err != nil {
		t.Fatal(err)
	}

	var resp struct {
		Version   string `json:"version"`
		Commit    string `json:"commit"`
		GoVersion string `json:"go_version"`
	}
	if err := json.Unmarshal(buf.Bytes(), &resp); err != nil {
		t.Fatalf("decode version JSON: %v\noutput: %s", err, buf.String())
	}
	if resp.Version == "" || resp.Commit == "" || resp.GoVersion == "" {
		t.Fatalf("unexpected version JSON: %+v", resp)
	}
}
