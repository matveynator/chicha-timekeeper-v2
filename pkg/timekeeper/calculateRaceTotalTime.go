package Timekeeper

import (
	"chicha/pkg/config"
	"chicha/pkg/data"
)

// Рассчитываем общее время гонки для текущего круга
func calculateRaceTotalTime(currentLap *Data.Lap, nonMatchingLaps []Data.Lap, config Config.Settings) {
	if currentLap.LapNumber == 0 {
		currentLap.RaceTotalTime = currentLap.LapTime
		currentLap.FasterOrSlowerThanPreviousLapTime = 0
		currentLap.LapIsStrange = false
	} else {
		for _, lap := range nonMatchingLaps {
			if currentLap.TagId == lap.TagId && currentLap.RaceId == lap.RaceId && lap.LapNumber == currentLap.LapNumber-1 {
				currentLap.RaceTotalTime = lap.RaceTotalTime + currentLap.LapTime

				if currentLap.LapNumber > 1 && !config.VARIABLE_DISTANCE_RACE {
					currentLap.FasterOrSlowerThanPreviousLapTime = currentLap.LapTime - lap.LapTime

					if (currentLap.LapTime - lap.LapTime) >= lap.LapTime {
						currentLap.LapIsStrange = true
					} else {
						currentLap.LapIsStrange = false
					}
				} else {
					currentLap.FasterOrSlowerThanPreviousLapTime = 0
					currentLap.LapIsStrange = false
				}
			}
		}
	}
}


