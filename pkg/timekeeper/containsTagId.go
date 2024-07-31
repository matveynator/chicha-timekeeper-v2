package Timekeeper

import (
	"chicha/pkg/data"
)

// Check slice contains TagID:
func containsTagId(laps []Data.Lap, tagId string) bool {
	for _, lap := range laps {
		if lap.TagId == tagId {
			return true
		}
	}
	return false
}
