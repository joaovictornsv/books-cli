package models

import (
	"testing"
)

func TestParseStatus(t *testing.T) {
	tests := []struct {
		input   string
		want    Status
		wantErr bool
	}{
		{"READ", StatusRead, false},
		{"reading", StatusReading, false},
		{"NOT_STARTED", StatusNotStarted, false},
		{"invalid", "", true},
		{"", "", true},
	}

	for _, tt := range tests {
		got, err := ParseStatus(tt.input)
		if tt.wantErr {
			if err == nil {
				t.Errorf("ParseStatus(%q) expected error", tt.input)
			}
			continue
		}
		if err != nil {
			t.Fatalf("ParseStatus(%q): %v", tt.input, err)
		}
		if got != tt.want {
			t.Errorf("ParseStatus(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestValidateBool01(t *testing.T) {
	if err := ValidateBool01("x", 0); err != nil {
		t.Fatal(err)
	}
	if err := ValidateBool01("x", 1); err != nil {
		t.Fatal(err)
	}
	if err := ValidateBool01("x", 2); err == nil {
		t.Fatal("expected error for value 2")
	}
}

func TestBookValidate(t *testing.T) {
	author := "Author"
	book := Book{
		Title:          "Title",
		Author:         &author,
		Status:         StatusNotStarted,
		PriorityToBuy:  0,
		EligibleToSell: 0,
		Sold:           0,
		AddedAt:        NowTimestamp(),
	}
	if err := book.Validate(); err != nil {
		t.Fatal(err)
	}

	book.Title = ""
	if err := book.Validate(); err == nil {
		t.Fatal("expected error for empty title")
	}
}

func TestBookValidateForCreate(t *testing.T) {
	book := Book{
		Title:          "Title",
		Status:         StatusNotStarted,
		PriorityToBuy:  0,
		EligibleToSell: 0,
		Sold:           0,
	}
	if err := book.ValidateForCreate(); err != nil {
		t.Fatal(err)
	}

	book.AddedAt = ""
	if err := book.ValidateForCreate(); err != nil {
		t.Fatal("ValidateForCreate should not require added_at")
	}
	if err := book.Validate(); err == nil {
		t.Fatal("expected full Validate to require added_at")
	}
}

func TestToBool01(t *testing.T) {
	if ToBool01(true) != 1 {
		t.Fatal("expected 1")
	}
	if ToBool01(false) != 0 {
		t.Fatal("expected 0")
	}
	if !FromBool01(1) || FromBool01(0) {
		t.Fatal("FromBool01 mismatch")
	}
}
