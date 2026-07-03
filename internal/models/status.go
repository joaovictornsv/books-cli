package models

import (
	"fmt"
	"strings"
)

type Status string

const (
	StatusRead       Status = "READ"
	StatusReading    Status = "READING"
	StatusNotStarted Status = "NOT_STARTED"
	StatusToBuy      Status = "TO_BUY"
	StatusArchived   Status = "ARCHIVED"
)

func ParseStatus(s string) (Status, error) {
	status := Status(strings.ToUpper(strings.TrimSpace(s)))
	if !status.Valid() {
		return "", fmt.Errorf("invalid status %q: must be one of READ, READING, NOT_STARTED, TO_BUY, ARCHIVED", s)
	}
	return status, nil
}

func (s Status) Valid() bool {
	switch s {
	case StatusRead, StatusReading, StatusNotStarted, StatusToBuy, StatusArchived:
		return true
	default:
		return false
	}
}

func (s Status) String() string {
	return string(s)
}
