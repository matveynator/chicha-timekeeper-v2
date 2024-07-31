package Timekeeper

import (
	"chicha/pkg/config"
	"chicha/pkg/data"
	"fmt"
)

// Обновление текущего круга на основе предыдущих данных
func updateCurrentLap(currentLap Data.Lap, previousLaps []Data.Lap, currentTimekeeperTask Data.RawData, config Config.Settings) ([]Data.Lap, error) {
	var nonMatchingLaps []Data.Lap
	var err error

	for _, previousLap := range previousLaps {
		if currentLap.TagId == previousLap.TagId && currentLap.RaceId == previousLap.RaceId && currentLap.LapNumber == previousLap.LapNumber {
			updateLapTimes(&currentLap, previousLap, currentTimekeeperTask)
		} else {
			nonMatchingLaps = append(nonMatchingLaps, previousLap)
		}
	}

	currentLap.LapTime, err = calculateLapTime(currentLap, nonMatchingLaps, config)
	if err != nil {
		return nil, fmt.Errorf("Ошибка в расчете времени круга: %s", err)
	}

	if !config.VARIABLE_DISTANCE_RACE {
		currentLap.BestLapTime, currentLap.BestLapNumber, err = calculateMyBestLapTime(currentLap, nonMatchingLaps)
		if err != nil {
			return nil, fmt.Errorf("Ошибка в расчете лучшего времени круга: %s", err)
		}
	}

	calculateRaceTotalTime(&currentLap, nonMatchingLaps, config)

	updatedLaps := append(nonMatchingLaps, currentLap)

  if config.AVERAGE_RESULTS {
    sortLapsAscByDiscoveryAverageUnixTime(updatedLaps)
  } else {
    sortLapsAscByDiscoveryMinimalUnixTime(updatedLaps)
  } 

	return updatedLaps, nil
}

