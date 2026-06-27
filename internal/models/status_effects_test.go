package models

import "testing"

func TestApplyStatusSideEffects(t *testing.T) {
	now := "2024-06-01T12:00:00Z"
	existingStart := "2024-05-01T00:00:00Z"

	t.Run("enter reading sets started_at", func(t *testing.T) {
		after := Book{Status: StatusReading}
		got := ApplyStatusSideEffects(after, now)
		if got.StartedAt == nil || *got.StartedAt != now {
			t.Fatalf("started_at = %v, want %q", got.StartedAt, now)
		}
		if got.FinishedAt != nil {
			t.Fatal("expected finished_at to remain unset")
		}
	})

	t.Run("enter read sets finished_at and keeps started_at", func(t *testing.T) {
		after := Book{Status: StatusRead, StartedAt: &existingStart}
		got := ApplyStatusSideEffects(after, now)
		if got.StartedAt == nil || *got.StartedAt != existingStart {
			t.Fatalf("started_at = %v, want %q", got.StartedAt, existingStart)
		}
		if got.FinishedAt == nil || *got.FinishedAt != now {
			t.Fatalf("finished_at = %v, want %q", got.FinishedAt, now)
		}
	})

	t.Run("leave read clears finished_at", func(t *testing.T) {
		finished := "2024-05-02T00:00:00Z"
		after := Book{Status: StatusReading, StartedAt: &existingStart, FinishedAt: &finished}
		got := ApplyStatusSideEffects(after, now)
		if got.FinishedAt != nil {
			t.Fatalf("expected finished_at cleared, got %v", got.FinishedAt)
		}
	})

	t.Run("leave reading clears started_at unless read", func(t *testing.T) {
		after := Book{Status: StatusNotStarted, StartedAt: &existingStart}
		got := ApplyStatusSideEffects(after, now)
		if got.StartedAt != nil {
			t.Fatalf("expected started_at cleared, got %v", got.StartedAt)
		}
	})
}
