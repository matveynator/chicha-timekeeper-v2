package Timekeeper

import (
	"time"
	"sort"
	"chicha/pkg/config"
	"chicha/pkg/data"
)

func getCurrentRaceId(previousLaps []Data.Lap, config Config.Settings) (raceId uint, err error) {
	if cap(previousLaps) == 0 {
		raceId = 1
	} else {
		previousLapsCopy := previousLaps
		sort.Slice(previousLapsCopy, func(i, j int) bool {
			//sort descending by DiscoveryAverageUnixTime or DiscoveryMinimalUnixTime (depending on config.AVERAGE_RESULTS setting:
			if config.AVERAGE_RESULTS {
				return previousLapsCopy[i].DiscoveryAverageUnixTime > previousLapsCopy[j].DiscoveryAverageUnixTime
			} else {
				return previousLapsCopy[i].DiscoveryMinimalUnixTime > previousLapsCopy[j].DiscoveryMinimalUnixTime
			}
		})

		var location *time.Location
		location, err = time.LoadLocation(config.TIME_ZONE)
		if err != nil {
			//stop and return error
			return
		} else {
			var lastLapTime time.Time 
			if config.AVERAGE_RESULTS {
				lastLapTime = time.Unix(previousLapsCopy[0].DiscoveryAverageUnixTime, 0).In(location)
			} else {
				lastLapTime = time.Unix(previousLapsCopy[0].DiscoveryMinimalUnixTime, 0).In(location)
			}
			currentTime := time.Now().In(location)

			if currentTime.After(lastLapTime.Add(config.RACE_TIMEOUT_DURATION)) {
				//if previous race time expired (time outed) - increment race id:
				raceId = previousLaps[0].RaceId + 1	
			} else {
				//return current race id:
				raceId = previousLaps[0].RaceId
			}
		}
	}
	return
}
