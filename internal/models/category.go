package models

import (
	"fmt"
	"strings"
)

type Category string

const (
	CategoryTheology            Category = "THEOLOGY"
	CategoryFiction             Category = "FICTION"
	CategorySoftware            Category = "SOFTWARE"
	CategoryPhilosophy          Category = "PHILOSOPHY"
	CategoryHistory             Category = "HISTORY"
	CategoryPersonalDevelopment Category = "PERSONAL_DEVELOPMENT"
	CategoryFinanceBusiness     Category = "FINANCE_BUSINESS"
	CategoryScience             Category = "SCIENCE"
	CategoryPoliticsCulture     Category = "POLITICS_CULTURE"
	CategoryBiography           Category = "BIOGRAPHY"
	CategoryOther               Category = "OTHER"
)

func ParseCategory(s string) (Category, error) {
	category := Category(strings.ToUpper(strings.TrimSpace(s)))
	if !category.Valid() {
		return "", fmt.Errorf(
			"invalid category %q: must be one of THEOLOGY, FICTION, SOFTWARE, PHILOSOPHY, HISTORY, PERSONAL_DEVELOPMENT, FINANCE_BUSINESS, SCIENCE, POLITICS_CULTURE, BIOGRAPHY, OTHER",
			s,
		)
	}
	return category, nil
}

func (c Category) Valid() bool {
	switch c {
	case CategoryTheology, CategoryFiction, CategorySoftware, CategoryPhilosophy,
		CategoryHistory, CategoryPersonalDevelopment, CategoryFinanceBusiness,
		CategoryScience, CategoryPoliticsCulture, CategoryBiography, CategoryOther:
		return true
	default:
		return false
	}
}

func (c Category) String() string {
	return string(c)
}
