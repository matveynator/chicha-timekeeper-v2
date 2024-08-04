package Timekeeper

import (
	"log"
	"sort"
	"time"
	"chicha/pkg/data"
	"chicha/pkg/config"
)

// getCurrentRaceId returns the race ID of the current lap based on previous laps.
// If there are no previous laps, the race ID will be set to 1.
func getCurrentRaceId(currentTimekeeperTask data.RawData, previousLaps []data.Lap, config Config.Settings) (raceId uint, err error) {
	if len(previousLaps) == 0 {
		raceId = 1 // If there are no previous laps, assign race ID 1.
	} else {
		previousLapsCopy := previousLaps // Create a copy of the previous laps to avoid altering the original slice.
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
				lastLapTime = time.UnixMilli(previousLapsCopy[0].DiscoveryAverageUnixTime).In(location)
			} else {
				lastLapTime = time.UnixMilli(previousLapsCopy[0].DiscoveryMinimalUnixTime).In(location)
			}
			currentTime := time.UnixMilli(currentTimekeeperTask.DiscoveryUnixTime).In(location)

			if currentTime.After(lastLapTime.Add(config.RACE_TIMEOUT_DURATION)) {
				// If previous race time expired (timed out), increment race ID:
				log.Println("race expired, starting new race.")
				raceId = previousLaps[0].RaceId + 1
			} else {
				// Return current race ID:
				raceId = previousLaps[0].RaceId
			}
		}
	}
	return
}

