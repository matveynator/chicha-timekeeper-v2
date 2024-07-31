package Timekeeper

import (
	"chicha/pkg/data"
)

// Обновление времен круга
func updateLapTimes(currentLap *Data.Lap, previousLap Data.Lap, currentTimekeeperTask Data.RawData) {
	if currentLap.DiscoveryMinimalUnixTime > previousLap.DiscoveryMinimalUnixTime {
		currentLap.DiscoveryMinimalUnixTime = previousLap.DiscoveryMinimalUnixTime
	}

	if currentLap.DiscoveryMaximalUnixTime < previousLap.DiscoveryMaximalUnixTime {
		currentLap.DiscoveryMaximalUnixTime = previousLap.DiscoveryMaximalUnixTime
	}

	currentLap.DiscoveryAverageUnixTime = (currentTimekeeperTask.DiscoveryUnixTime + previousLap.DiscoveryAverageUnixTime) / 2

	if previousLap.AverageResultsCount < 1 {
		previousLap.AverageResultsCount = 1
	}

	currentLap.AverageResultsCount = previousLap.AverageResultsCount + 1
}

