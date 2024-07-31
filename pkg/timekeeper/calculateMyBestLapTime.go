package Timekeeper

import (
	"chicha/pkg/data"
)

//  Подсчитываем мое лучшее время круга:
func calculateMyBestLapTime(currentLap Data.Lap, otherOldLaps []Data.Lap) (bestLapTime int64, bestLapNumber uint, err error) {
	if	currentLap.LapNumber == 0 {
		bestLapTime = 0
		bestLapNumber = 0
	} else if currentLap.LapNumber == 1 {
		bestLapTime = currentLap.LapTime
		bestLapNumber = 1
	} else {
		// Make a selection that includes only the complete laps that belong to me:
		var mySameRaceFullLaps []Data.Lap
		for _, otherOldLap := range otherOldLaps {
			if otherOldLap.RaceId == currentLap.RaceId &&  otherOldLap.TagId == currentLap.TagId && otherOldLap.LapNumber > 0  {
				mySameRaceFullLaps = append(mySameRaceFullLaps, otherOldLap)
			}
		}
		mySameRaceFullLaps = append(mySameRaceFullLaps, currentLap)

		// Sort by lap time ascending small->big :
		sortLapsAscLapTime(mySameRaceFullLaps)

		// Assign fastest lap to bestLapTime and bestLapNumber:
		bestLapTime = mySameRaceFullLaps[0].LapTime
		bestLapNumber = mySameRaceFullLaps[0].LapNumber
	}
	return
}

