package Timekeeper

import (
	"log"
	"time"
	"sort"
	"chicha/pkg/config"
	"chicha/pkg/data"
)

// getCurrentRaceId returns the race ID of the current lap based on previous laps.
// If there are no previous laps, the race ID will be set to 1.
func getCurrentRaceId(currentTimekeeperTask Data.RawData, previousLaps []Data.Lap, config Config.Settings) (raceId uint, err error) {
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

// getMyCurrentRaceLapNumber returns the current lap number of the current racer based on previous laps and current timekeeper task.
func getMyCurrentRaceLapNumber(currentTimekeeperTask Data.RawData, previousLaps []Data.Lap, config Config.Settings) (lapNumber uint, err error) {

	var myPreviousLaps []Data.Lap

	// Get current Race ID:
	var currentRaceId uint
	currentRaceId, err = getCurrentRaceId(currentTimekeeperTask, previousLaps, config)
	if err != nil {
		log.Println(err)
		return
	}

	// If previousLaps slice is empty, current lap number is 0:
	if len(previousLaps) == 0 {
		//log.Println("len(previousLaps) == 0 (1)")
		lapNumber = 0
	} else {
		for _, previousLap := range previousLaps {
			// Gather together only my laps (by my TagId) and only from current race:
			if previousLap.TagId == currentTimekeeperTask.TagId {
				if previousLap.RaceId == currentRaceId {
					myPreviousLaps = append(myPreviousLaps, previousLap)
				}
			}
		}
		//log.Println("len(myPreviousLaps)", len(myPreviousLaps))
		//if my results empty:
		if len(myPreviousLaps) == 0 {
			//log.Println("len(myPreviousLaps) == 0 (2)")
			lapNumber = 0
		} else {
			//if my results are not empty:
			//log.Println("len(myPreviousLaps): ", len(myPreviousLaps))
			sort.Slice(myPreviousLaps, func(i, j int) bool {
				//sort descending by DiscoveryAverageUnixTime or DiscoveryMinimalUnixTime (depending on config.AVERAGE_RESULTS setting:
				if config.AVERAGE_RESULTS {
					return myPreviousLaps[i].DiscoveryAverageUnixTime > myPreviousLaps[j].DiscoveryAverageUnixTime
				} else {
					return myPreviousLaps[i].DiscoveryMinimalUnixTime > myPreviousLaps[j].DiscoveryMinimalUnixTime
				}
			})

			var location *time.Location
			location, err = time.LoadLocation(config.TIME_ZONE)
			if err != nil {
				//stop and return error
				log.Println("Error: time.LoadLocation(config.TIME_ZONE) -", err)
				return
			} else {
				var myLastLapTime time.Time
				if config.AVERAGE_RESULTS {
					myLastLapTime = time.UnixMilli(myPreviousLaps[0].DiscoveryAverageUnixTime).In(location)
				} else {
					myLastLapTime = time.UnixMilli(myPreviousLaps[0].DiscoveryMinimalUnixTime).In(location)
				}
				myCurrentLapTime := time.UnixMilli(currentTimekeeperTask.DiscoveryUnixTime).In(location)

				if myCurrentLapTime.After(myLastLapTime.Add(config.MINIMAL_LAP_TIME_DURATION)) {
					lapNumber = myPreviousLaps[0].LapNumber + 1
				} else {
					lapNumber = myPreviousLaps[0].LapNumber
					//log.Println(fmt.Sprintf("Error: Received lap result is outside time specs: > average duration = %s and less than minimal lap time = %s", config.AVERAGE_DURATION, config.MINIMAL_LAP_TIME_DURATION))
				}
			}
		}
	}
	return
}

