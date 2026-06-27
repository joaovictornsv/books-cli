package models

// ApplyStatusSideEffects updates reading timestamps when status changes.
// Pass an explicit now timestamp (RFC3339) so callers can test deterministically.
func ApplyStatusSideEffects(after Book, now string) Book {
	result := after

	if after.Status == StatusReading {
		if result.StartedAt == nil {
			ts := now
			result.StartedAt = &ts
		}
	} else if after.Status != StatusRead {
		result.StartedAt = nil
	}

	if after.Status == StatusRead {
		if result.FinishedAt == nil {
			ts := now
			result.FinishedAt = &ts
		}
	} else {
		result.FinishedAt = nil
	}

	return result
}
