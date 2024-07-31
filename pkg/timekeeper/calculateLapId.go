package Timekeeper

import (
	"chicha/pkg/data"
)

// Calculate Data.Lap.Id based on previous laps (if available):
func calculateLapId(currentLap Data.Lap, laps []Data.Lap) (id int64) {
	// If no previous laps available:
	if len(laps) == 0 {
		// This is the first lap! - set id to 1:
		id = 1
	} else {
		// Previous data available (id will be >= 1)
		var lapFound bool = false
		for _, lap := range laps {
			// For the save TagId, RaceId and LapNumber set SAME id:
			if currentLap.TagId == lap.TagId && currentLap.RaceId == lap.RaceId && currentLap.LapNumber == lap.LapNumber {
				lapFound = true
				return lap.Id
			}
		}

		// If no such lap in previous data - calculate next lap id:
		if lapFound == false {
			var lapId int64 = 0
			// Search for highest lap id in previous data:
			for _, lap := range laps {
				if lap.Id > lapId {
					lapId = lap.Id
				}
			}
			// Return incremented highest id:
			id = lapId+1
		}
	}

	return
}

