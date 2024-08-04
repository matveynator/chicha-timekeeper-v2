package Timekeeper

import (
	"log"
	"sort"
	"time"
	"chicha/pkg/data"
	"chicha/pkg/config"
)

// Plan:

//currentLap.Id = +
//currentLap.TagId = +
//currentLap.DiscoveryMinimalUnixTime = +
//currentLap.DiscoveryMaximalUnixTime = +
//currentLap.DiscoveryAverageUnixTime = +
//currentLap.AverageResultsCount = +

//currentLap.RaceId = +
//currentLap.RaceTotalTime = +
//currentLap.RacePosition = +

//currentLap.TimeBehindTheLeader = +

//currentLap.LapNumber = +
//currentLap.LapTime = +
//currentLap.LapPosition = +
//currentLap.LapIsLatest = +

// Если гонка не VARIABLE_DISTANCE_RACE (false):
//currentLap.FasterOrSlowerThanPreviousLapTime = +
//currentLap.LapIsStrange = +
//currentLap.BestLapTime = +
//currentLap.BestLapNumber = +
//currentLap.FastestLapInThisRace = +


// getMyCurrentRaceLapNumber returns the current lap number of the current racer based on previous laps and current timekeeper task.
func getMyCurrentRaceLapNumber(currentTimekeeperTask data.RawData, previousLaps []data.Lap, config Config.Settings) (lapNumber uint, err error) {

	var myPreviousLaps []data.Lap

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


