package models

import "testing"

func TestParseCategory(t *testing.T) {
	tests := []struct {
		input   string
		want    Category
		wantErr bool
	}{
		{"FICTION", CategoryFiction, false},
		{"biography", CategoryBiography, false},
		{"PERSONAL_DEVELOPMENT", CategoryPersonalDevelopment, false},
		{"invalid", "", true},
		{"", "", true},
	}

	for _, tt := range tests {
		got, err := ParseCategory(tt.input)
		if tt.wantErr {
			if err == nil {
				t.Errorf("ParseCategory(%q) expected error", tt.input)
			}
			continue
		}
		if err != nil {
			t.Fatalf("ParseCategory(%q): %v", tt.input, err)
		}
		if got != tt.want {
			t.Errorf("ParseCategory(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}
