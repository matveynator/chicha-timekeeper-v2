package Timekeeper

import (
	"log"
	"time"
	"sort"
	"chicha/pkg/config"
	"chicha/pkg/data"
)

func getCurrentRaceId(currentTimekeeperTask Data.RawData, previousLaps []Data.Lap, config Config.Settings) (raceId uint, err error) {
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
				lastLapTime = time.UnixMilli(previousLapsCopy[0].DiscoveryAverageUnixTime).In(location)
			} else {
				lastLapTime = time.UnixMilli(previousLapsCopy[0].DiscoveryMinimalUnixTime).In(location)
			}
			currentTime := time.UnixMilli(currentTimekeeperTask.DiscoveryUnixTime).In(location)

			if currentTime.After(lastLapTime.Add(config.RACE_TIMEOUT_DURATION)) {
				//if previous race time expired (time outed) - increment race id:
				log.Println("race expired, starting new race.")
				raceId = previousLaps[0].RaceId + 1	
			} else {
				//return current race id:
				raceId = previousLaps[0].RaceId
			}
		}
	}
	return
}

func getCurrentLapNumber(currentTimekeeperTask Data.RawData, previousLaps []Data.Lap, config Config.Settings) (lapNumber uint, err error) {
	if cap(previousLaps) == 0 {
		lapNumber = 0
	} else {

		//get only my results from slice:
		var myPreviousLaps []Data.Lap 
		for _, previousLap := range previousLaps {
			if previousLap.TagId == currentTimekeeperTask.TagId {
				myPreviousLaps = append(myPreviousLaps, previousLap)
			}
		}
		//if my results available:
		if cap(myPreviousLaps)==0 {
			log.Println("cap(myPreviousLaps):", cap(myPreviousLaps))
			lapNumber = 0
		} else {
		  log.Println("cap(myPreviousLaps):", cap(myPreviousLaps))
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
				return
			} else {
				var lastLapTime time.Time
				if config.AVERAGE_RESULTS {
					lastLapTime = time.UnixMilli(myPreviousLaps[0].DiscoveryAverageUnixTime).In(location)
				} else {
					lastLapTime = time.UnixMilli(myPreviousLaps[0].DiscoveryMinimalUnixTime).In(location)
				}
				currentTime := time.UnixMilli(currentTimekeeperTask.DiscoveryUnixTime).In(location)

				if currentTime.After(lastLapTime.Add(config.MINIMAL_LAP_TIME_DURATION)) {
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

