package models

import "testing"

func TestPaginationValidate(t *testing.T) {
	if err := (Pagination{Page: 1, Limit: 20}).Validate(); err != nil {
		t.Fatal(err)
	}
	if err := (Pagination{Page: 0, Limit: 20}).Validate(); err == nil {
		t.Fatal("expected error for page 0")
	}
	if err := (Pagination{Page: 1, Limit: 0}).Validate(); err == nil {
		t.Fatal("expected error for limit 0")
	}
}

func TestPaginationOffset(t *testing.T) {
	p := Pagination{Page: 2, Limit: 10}
	if p.Offset() != 10 {
		t.Fatalf("got offset %d, want 10", p.Offset())
	}
}
