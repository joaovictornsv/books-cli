package models

import (
	"fmt"
	"strings"
	"time"
)

type Book struct {
	ID              int64   `json:"id"`
	Title           string  `json:"title"`
	Author          *string `json:"author,omitempty"`
	Status          Status  `json:"status"`
	PriorityToBuy   int     `json:"priority_to_buy"`
	EligibleToSell  int     `json:"eligible_to_sell"`
	Sold            int     `json:"sold"`
	Notes           *string `json:"notes,omitempty"`
	AddedAt         string  `json:"added_at"`
	StartedAt       *string `json:"started_at,omitempty"`
	FinishedAt      *string `json:"finished_at,omitempty"`
}

type BookPatch struct {
	Title          *string
	Author         *string
	ClearAuthor    bool
	Status         *Status
	PriorityToBuy  *int
	EligibleToSell *int
	Sold           *int
	Notes          *string
}

func ToBool01(v bool) int {
	if v {
		return 1
	}
	return 0
}

func FromBool01(v int) bool {
	return v == 1
}

func ValidateBool01(name string, v int) error {
	if v != 0 && v != 1 {
		return fmt.Errorf("%s must be 0 or 1", name)
	}
	return nil
}

func (b *Book) ValidateForCreate() error {
	if strings.TrimSpace(b.Title) == "" {
		return fmt.Errorf("title is required")
	}
	if !b.Status.Valid() {
		return fmt.Errorf("invalid status %q", b.Status)
	}
	if err := ValidateBool01("priority_to_buy", b.PriorityToBuy); err != nil {
		return err
	}
	if err := ValidateBool01("eligible_to_sell", b.EligibleToSell); err != nil {
		return err
	}
	if err := ValidateBool01("sold", b.Sold); err != nil {
		return err
	}
	return nil
}

func (b *Book) Validate() error {
	if err := b.ValidateForCreate(); err != nil {
		return err
	}
	if err := validateTimestamp("added_at", b.AddedAt, true); err != nil {
		return err
	}
	if b.StartedAt != nil {
		if err := validateTimestamp("started_at", *b.StartedAt, false); err != nil {
			return err
		}
	}
	if b.FinishedAt != nil {
		if err := validateTimestamp("finished_at", *b.FinishedAt, false); err != nil {
			return err
		}
	}
	return nil
}

func (p *BookPatch) Validate() error {
	if p.Title != nil && strings.TrimSpace(*p.Title) == "" {
		return fmt.Errorf("title cannot be empty")
	}
	if p.Status != nil && !p.Status.Valid() {
		return fmt.Errorf("invalid status %q", *p.Status)
	}
	if p.PriorityToBuy != nil {
		if err := ValidateBool01("priority_to_buy", *p.PriorityToBuy); err != nil {
			return err
		}
	}
	if p.EligibleToSell != nil {
		if err := ValidateBool01("eligible_to_sell", *p.EligibleToSell); err != nil {
			return err
		}
	}
	if p.Sold != nil {
		if err := ValidateBool01("sold", *p.Sold); err != nil {
			return err
		}
	}
	return nil
}

func NowTimestamp() string {
	return time.Now().UTC().Format(time.RFC3339)
}

func validateTimestamp(name, value string, required bool) error {
	if value == "" {
		if required {
			return fmt.Errorf("%s is required", name)
		}
		return nil
	}
	if _, err := time.Parse(time.RFC3339, value); err != nil {
		return fmt.Errorf("%s must be a valid ISO 8601 timestamp: %w", name, err)
	}
	return nil
}

