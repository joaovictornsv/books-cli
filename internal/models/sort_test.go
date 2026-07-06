package models

import "testing"

func TestParseSortField(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    SortField
		wantErr bool
	}{
		{name: "title", input: "title", want: SortFieldTitle},
		{name: "mixed case", input: "Title", want: SortFieldTitle},
		{name: "whitespace", input: "  author  ", want: SortFieldAuthor},
		{name: "finished_at", input: "finished_at", want: SortFieldFinishedAt},
		{name: "empty", input: "", wantErr: true},
		{name: "invalid", input: "invalid", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseSortField(tt.input)
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

func TestParseSortOrder(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    SortOrder
		wantErr bool
	}{
		{name: "asc", input: "asc", want: SortOrderAsc},
		{name: "desc upper", input: "DESC", want: SortOrderDesc},
		{name: "whitespace", input: "  desc  ", want: SortOrderDesc},
		{name: "empty", input: "", wantErr: true},
		{name: "invalid", input: "up", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseSortOrder(tt.input)
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
