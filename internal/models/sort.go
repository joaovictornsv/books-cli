package models

import (
	"fmt"
	"strings"
)

type SortField string

const (
	SortFieldID         SortField = "id"
	SortFieldTitle      SortField = "title"
	SortFieldAuthor     SortField = "author"
	SortFieldStatus     SortField = "status"
	SortFieldAddedAt    SortField = "added_at"
	SortFieldStartedAt  SortField = "started_at"
	SortFieldFinishedAt SortField = "finished_at"
)

type SortOrder string

const (
	SortOrderAsc  SortOrder = "asc"
	SortOrderDesc SortOrder = "desc"
)

type Sort struct {
	Field SortField
	Order SortOrder
}

func DefaultSort() Sort {
	return Sort{Field: SortFieldID, Order: SortOrderAsc}
}

func ParseSortField(s string) (SortField, error) {
	field := SortField(strings.ToLower(strings.TrimSpace(s)))
	switch field {
	case SortFieldID, SortFieldTitle, SortFieldAuthor, SortFieldStatus,
		SortFieldAddedAt, SortFieldStartedAt, SortFieldFinishedAt:
		return field, nil
	default:
		return "", fmt.Errorf("invalid sort field %q: must be one of id, title, author, status, added_at, started_at, finished_at", s)
	}
}

func ParseSortOrder(s string) (SortOrder, error) {
	order := SortOrder(strings.ToLower(strings.TrimSpace(s)))
	switch order {
	case SortOrderAsc, SortOrderDesc:
		return order, nil
	default:
		return "", fmt.Errorf("invalid sort order %q: must be asc or desc", s)
	}
}

func (s Sort) Validate() error {
	if s.Field == "" && s.Order == "" {
		return nil
	}
	if _, err := ParseSortField(string(s.Field)); err != nil {
		return err
	}
	if _, err := ParseSortOrder(string(s.Order)); err != nil {
		return err
	}
	return nil
}

func (s Sort) WithDefaults() Sort {
	if s.Field == "" && s.Order == "" {
		return DefaultSort()
	}
	out := s
	if out.Field == "" {
		out.Field = SortFieldID
	}
	if out.Order == "" {
		out.Order = SortOrderAsc
	}
	return out
}

func (f SortField) Nullable() bool {
	switch f {
	case SortFieldAuthor, SortFieldStartedAt, SortFieldFinishedAt:
		return true
	default:
		return false
	}
}
