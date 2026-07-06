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

func TestSortOrderByClause(t *testing.T) {
	tests := []struct {
		name    string
		sort    Sort
		want    string
		wantErr bool
	}{
		{
			name: "title asc",
			sort: Sort{Field: SortFieldTitle, Order: SortOrderAsc},
			want: " ORDER BY title ASC",
		},
		{
			name: "finished_at desc nulls last",
			sort: Sort{Field: SortFieldFinishedAt, Order: SortOrderDesc},
			want: " ORDER BY finished_at DESC NULLS LAST",
		},
		{
			name:    "invalid field",
			sort:    Sort{Field: "bad", Order: SortOrderAsc},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.sort.OrderByClause()
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error")
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if got != tt.want {
				t.Fatalf("got %q, want %q", got, tt.want)
			}
		})
	}
}
