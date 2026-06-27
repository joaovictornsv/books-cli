package db

import "testing"

func TestEscapeLike(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"plain", "plain"},
		{`100%`, `100\%`},
		{`a_b`, `a\_b`},
		{`back\slash`, `back\\slash`},
	}

	for _, tt := range tests {
		if got := escapeLike(tt.input); got != tt.want {
			t.Errorf("escapeLike(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}
