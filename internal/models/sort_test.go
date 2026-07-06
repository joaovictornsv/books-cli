package models

import "testing"

func TestParseSortField(t *testing.T) {
	field, err := ParseSortField("title")
	if err != nil {
		t.Fatal(err)
	}
	if field != SortFieldTitle {
		t.Fatalf("got %q", field)
	}

	_, err = ParseSortField("invalid")
	if err == nil {
		t.Fatal("expected error for invalid sort field")
	}
}

func TestParseSortOrder(t *testing.T) {
	order, err := ParseSortOrder("DESC")
	if err != nil {
		t.Fatal(err)
	}
	if order != SortOrderDesc {
		t.Fatalf("got %q", order)
	}

	_, err = ParseSortOrder("up")
	if err == nil {
		t.Fatal("expected error for invalid sort order")
	}
}

func TestSortWithDefaults(t *testing.T) {
	got := Sort{}.WithDefaults()
	want := DefaultSort()
	if got != want {
		t.Fatalf("got %+v, want %+v", got, want)
	}

	got = Sort{Field: SortFieldTitle}.WithDefaults()
	if got.Field != SortFieldTitle || got.Order != SortOrderAsc {
		t.Fatalf("got %+v", got)
	}
}
