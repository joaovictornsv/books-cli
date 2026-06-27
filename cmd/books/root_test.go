package main

import "testing"

func TestParseID(t *testing.T) {
	tests := []struct {
		input   string
		want    int64
		wantErr bool
	}{
		{"1", 1, false},
		{"42", 42, false},
		{"0", 0, true},
		{"-1", 0, true},
		{"abc", 0, true},
		{"", 0, true},
	}

	for _, tt := range tests {
		got, err := parseID(tt.input)
		if tt.wantErr {
			if err == nil {
				t.Errorf("parseID(%q) expected error", tt.input)
			}
			continue
		}
		if err != nil {
			t.Fatalf("parseID(%q): %v", tt.input, err)
		}
		if got != tt.want {
			t.Errorf("parseID(%q) = %d, want %d", tt.input, got, tt.want)
		}
	}
}
