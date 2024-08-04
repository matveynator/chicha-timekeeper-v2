package Timekeeper

import (
	"chicha/pkg/data"
)

// Calculate next Id:
func getNextLapId (laps []data.Lap) (id int64) {
	// Set initial Id:
	id = 0
	// Search for highest Id:
	for _, lap := range laps { 
		if lap.Id > id {
			id = lap.Id
		}
	}
	// Return next Id:
	return id + 1
}
