package models

import "fmt"

const DefaultPageLimit = 20

type Pagination struct {
	Page  int
	Limit int
}

func (p Pagination) Enabled() bool {
	return p.Limit > 0
}

func (p Pagination) Offset() int {
	return (p.Page - 1) * p.Limit
}

func (p Pagination) Validate() error {
	if p.Page < 1 {
		return fmt.Errorf("page must be >= 1")
	}
	if p.Limit < 1 {
		return fmt.Errorf("limit must be >= 1")
	}
	return nil
}
