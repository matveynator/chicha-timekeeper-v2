package Timekeeper

import (
	"chicha/pkg/config"
	"chicha/pkg/data"
)

// Инициализация текущего круга
func initializeCurrentLap(currentTimekeeperTask Data.RawData, previousLaps []Data.Lap, config Config.Settings) (Data.Lap, error) {
	var currentLap Data.Lap
	var err error

	if len(previousLaps) == 0 {
		currentLap = Data.Lap{
			Id:                        1,
			TagId:                     currentTimekeeperTask.TagId,
			DiscoveryMinimalUnixTime:  currentTimekeeperTask.DiscoveryUnixTime,
			DiscoveryAverageUnixTime:  currentTimekeeperTask.DiscoveryUnixTime,
			AverageResultsCount:       1,
			RaceId:                    1,
			RacePosition:              1,
			TimeBehindTheLeader:       0,
			LapNumber:                 0,
			LapPosition:               1,
			LapIsLatest:               true,
			RaceFinished:              false,
			RaceTotalTime:             0,
		}
	} else {
		currentLap = Data.Lap{
			TagId:                     currentTimekeeperTask.TagId,
			DiscoveryMinimalUnixTime:  currentTimekeeperTask.DiscoveryUnixTime,
			DiscoveryMaximalUnixTime:  currentTimekeeperTask.DiscoveryUnixTime,
			DiscoveryAverageUnixTime:  currentTimekeeperTask.DiscoveryUnixTime,
		}

		currentLap.RaceId, err = getCurrentRaceId(currentTimekeeperTask, previousLaps, config)
		if err != nil {
			return currentLap, err
		}

		currentLap.LapNumber, err = getMyCurrentRaceLapNumber(currentTimekeeperTask, previousLaps, config)
		if err != nil {
			return currentLap, err
		}
	}

	currentLap.Id = calculateLapId(currentLap, previousLaps)
	return currentLap, nil
}

