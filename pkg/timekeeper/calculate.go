package Timekeeper

import (
	"log"
	"sort"
	"time"
	"chicha/pkg/data"
	"chicha/pkg/config"
)

// Check time in data is valid:
func checkLapIsValid(lapToCheck Data.RawData, previousLaps []Data.Lap, config Config.Settings) (valid bool, err error) {

	// Get location from config time zone:
	var location *time.Location
	location, err = time.LoadLocation(config.TIME_ZONE)
	if err != nil {
		// Stop and return error:
		valid = false
		return
	}

	// Calculate my current time:
	myCurrentTime := time.UnixMilli(lapToCheck.DiscoveryUnixTime).In(location)

	// Calculate current system time:
	currentSystemTime := time.Now().In(location)

	// Check if current system time is after lap time + config.MINIMAL_LAP_TIME_DURATION (some very old data received?):
	if currentSystemTime.After(myCurrentTime.Add(config.MINIMAL_LAP_TIME_DURATION)) {
		// Stop and return error:
		valid = false
		return
	}

	// Check: Do we have any previous race data?
	if len(previousLaps) == 0 {
		// This is the first race and lap 
		valid = true
		return
	} else {
		// Sort previous laps by DiscoveryUnixTime:
		previousLapsCopy := previousLaps // Create a copy of the previous laps to avoid altering the original slice.
		sort.Slice(previousLapsCopy, func(i, j int) bool {
			//sort descending by DiscoveryAverageUnixTime or DiscoveryMinimalUnixTime (depending on config.AVERAGE_RESULTS setting:
			if config.AVERAGE_RESULTS {
				return previousLapsCopy[i].DiscoveryAverageUnixTime > previousLapsCopy[j].DiscoveryAverageUnixTime
			} else {
				return previousLapsCopy[i].DiscoveryMinimalUnixTime > previousLapsCopy[j].DiscoveryMinimalUnixTime
			}
		})
		// Get anyones very last time added to memory: 
		var lastLapTime time.Time
		if config.AVERAGE_RESULTS {
			lastLapTime = time.UnixMilli(previousLapsCopy[0].DiscoveryAverageUnixTime).In(location)
		} else {
			lastLapTime = time.UnixMilli(previousLapsCopy[0].DiscoveryMinimalUnixTime).In(location)
		}
		// Check if my current time is after config.RACE_TIMEOUT_DURATION:
		if myCurrentTime.After(lastLapTime.Add(config.RACE_TIMEOUT_DURATION)) {
			// If previous race time expired (timed out), data is valid:
			log.Println("Time is valid as race expired and we will be starting new race.")
			valid = true
			return
		} else {
			// If race not expired, calculate lap time (valid if > minimal lap time before previous results) and not more than precision time.
			// Create slice with only my data in it:
			var myPreviousLaps []Data.Lap
			for _, myPreviousLap := range previousLapsCopy {
				if lapToCheck.TagId == myPreviousLap.TagId {
					// Add my previous data to slice:	
					myPreviousLaps = append(myPreviousLaps, myPreviousLap)
				}
			}
			// Check if I have some previous laps in this race:
			if len(myPreviousLaps) > 0 {
				// I have some previous laps in this race.

				// Sort my laps by discovery unix time descending (big -> small) depending on config.AVERAGE_RESULTS global setting:
				if config.AVERAGE_RESULTS {
					sortLapsDescByDiscoveryAverageUnixTime(myPreviousLaps)
				} else {
					sortLapsDescByDiscoveryMinimalUnixTime(myPreviousLaps)
				}
				// Get my last lap time:
				var myLastLapTime time.Time
				if config.AVERAGE_RESULTS {
					myLastLapTime = time.UnixMilli(myPreviousLaps[0].DiscoveryAverageUnixTime).In(location)
				} else {
					myLastLapTime = time.UnixMilli(myPreviousLaps[0].DiscoveryMinimalUnixTime).In(location)
				}
				// Data is invalid if my time is after config.AVERAGE_DURATION and before than config.MINIMAL_LAP_TIME_DURATION:
				if myCurrentTime.After(myLastLapTime.Add(config.AVERAGE_DURATION)) && myCurrentTime.Before(myLastLapTime.Add(config.MINIMAL_LAP_TIME_DURATION)) {
					// Time is more than config.AVERAGE_DURATION and less than config.MINIMAL_LAP_TIME_DURATION - is invalid:
					valid = false
					return
				} else {
					valid = true
					return
				}	
			} else {
				// Race is valid and I have no laps in it yet - set lap is valid:
				valid = true
				return
			}
		}
	}
	return
}


func calculateRacePosition(currentLap Data.Lap, otherOldLaps []Data.Lap, raceType string) (racePosition uint, err error) {

	if ( raceType == "mass-start" ) {

	} else if (raceType == "delayed-start") {

	}

	return

}


// Calculate next Id:
func getNextId (laps []Data.Lap) (id int64) {
	// Set initial Id:
	id = 0
	// Search for highest Id:
	for _, lap := range laps { 
		if lap.Id > id {
			id = lap.Id
		}
	}
	// Return next Id:
	return id + 1
}

// Sort slice by discovery average unix time descending (big -> small):
func sortLapsDescByDiscoveryAverageUnixTime (lapsToSort []Data.Lap) {
	sort.Slice(lapsToSort, func(i, j int) bool {
		//log.Println("Average: sorting laps by DiscoveryAverageUnixTime")
		return lapsToSort[i].DiscoveryAverageUnixTime > lapsToSort[j].DiscoveryAverageUnixTime
	})
}

// Sort slice by discovery minimal unix time descending (big -> small):
func sortLapsDescByDiscoveryMinimalUnixTime (lapsToSort []Data.Lap) {
	sort.Slice(lapsToSort, func(i, j int) bool {
		//log.Println("Minimal: sorting laps by DiscoveryMinimalUnixTime")
		return lapsToSort[i].DiscoveryMinimalUnixTime > lapsToSort[j].DiscoveryMinimalUnixTime
	})
}

func calculateRaceInMemory (currentTimekeeperTask Data.RawData, previousLaps []Data.Lap, config Config.Settings) (currentLaps []Data.Lap, err error) {

	var currentLap Data.Lap

	var lapIsValid bool
	//  1. Check if current lap data is within lap timeouts and race timeouts and not very old (some expired data?) and combine with the rest of the race data:
	lapIsValid, err = checkLapIsValid(currentTimekeeperTask, previousLaps, config)
	if err != nil {
		log.Printf("Data validation error: %s \n", err)
		return
	} else {
		if !lapIsValid {
			// Data is invalid - drop it and return previousLaps:
			log.Println("Received invalid lap data. Skipping it.")
			currentLaps = previousLaps
			return
		}
	}

	// Check do we have any data allready in memory?
	if len(previousLaps) == 0 {
		// There is no race and lap in the memory - we will create the first race and lap:
		currentLap.Id = 1
		currentLap.TagId = currentTimekeeperTask.TagId
		currentLap.DiscoveryMinimalUnixTime = currentTimekeeperTask.DiscoveryUnixTime
		currentLap.DiscoveryAverageUnixTime = currentTimekeeperTask.DiscoveryUnixTime
		currentLap.AverageResultsCount = 1
		currentLap.RaceId = 1
		currentLap.RacePosition = 1
		currentLap.TimeBehindTheLeader = 0
		currentLap.LapNumber = 0
		currentLap.LapPosition = 1
		currentLap.LapIsCurrent = true
		currentLap.LapIsStrange = false
		currentLap.RaceFinished = false
		currentLap.RaceTotalTime = 0
	} else {
		// Not an empty slice (race allready running)

		// 2. Set well known data receved from rfid (currentLap.TagId, currentLap.DiscoveryMinimalUnixTime, currentLap.DiscoveryAverageUnixTime):
		currentLap.TagId = currentTimekeeperTask.TagId
		currentLap.DiscoveryMinimalUnixTime = currentTimekeeperTask.DiscoveryUnixTime
	  currentLap.DiscoveryAverageUnixTime = currentTimekeeperTask.DiscoveryUnixTime

		// 3. Calculate currentRaceId:
		previousLapsCopy := previousLaps
		currentLap.RaceId, err = getCurrentRaceId(currentTimekeeperTask, previousLapsCopy, config)
		if err != nil {
			log.Println(err)
			return
		}

		// 4. Calculate race lap number (currentLap.LapNumber):
		currentLap.LapNumber, err = getMyCurrentRaceLapNumber(currentTimekeeperTask, previousLapsCopy, config)
		if err != nil {
			log.Println(err)
			return
		}

		// currentLap.Id = + 
		//currentLap.TagId = + 
		//currentLap.DiscoveryMinimalUnixTime = + 
		//currentLap.DiscoveryAverageUnixTime = +
		//currentLap.AverageResultsCount = +
		//currentLap.RaceId = +
		//currentLap.RacePosition = 
		//currentLap.TimeBehindTheLeader = 
		//currentLap.LapNumber = 
		//currentLap.LapTime = ?
		//currentLap.LapPosition = 
  	//currentLap.BestLapTime =
	  //currentLap.BestLapNumber = 
		//currentLap.RaceTotalTime = 

		//currentLap.LapIsCurrent = 
		//currentLap.LapIsStrange = 
		//currentLap.RaceFinished = 
		//currentLap.RaceTotalTime = 

	}


	// 5. Calculate Id, DiscoveryMinimalUnixTime, DiscoveryAverageUnixTime, AverageResultsCount  for each laps with same TagId, RaceId, LapNumber. Remove duplicates BEGIN:.
	// Create slice copy with data from old laps not related to current lap BEGIN:
	var otherOldLaps []Data.Lap
	// Create slice copy with data from old laps not related to current lap END.
	for _, previousLap := range previousLaps {


		// Rewrite my lap data in memory with updated new data (DiscoveryMinimalUnixTime, DiscoveryAverageUnixTime, AverageResultsCount) removing duplicates BEGIN:
		if currentLap.TagId == previousLap.TagId && currentLap.RaceId == previousLap.RaceId && currentLap.LapNumber == previousLap.LapNumber {

			// Calculate DiscoveryMinimalUnixTime BEGIN:
			if currentLap.DiscoveryMinimalUnixTime > previousLap.DiscoveryMinimalUnixTime {
				currentLap.DiscoveryMinimalUnixTime = previousLap.DiscoveryMinimalUnixTime
			}
			// END.

			// Calculate DiscoveryAverageUnixTime BEGIN:
			currentLap.DiscoveryAverageUnixTime = (currentLap.DiscoveryAverageUnixTime + previousLap.DiscoveryAverageUnixTime)/2
			// END.

			// Calculate AverageResultsCount BEGIN:
			if previousLap.AverageResultsCount == 0 {
				previousLap.AverageResultsCount = 1
			} 
			currentLap.AverageResultsCount = previousLap.AverageResultsCount + 1
			// END.


			// Calculate Lap.Id BEGIN:
			if previousLap.Id > 0 {
				// Set old Id if available:
				currentLap.Id = previousLap.Id
			} else {
				// Calculate next Lap.Id if no old data available:
				currentLap.Id = getNextId (previousLaps)
			}
			// END.
		}


		// Rewrite my lap data in memory with updated new data (DiscoveryMinimalUnixTime, DiscoveryAverageUnixTime, AverageResultsCount) removing duplicates END.
		// Recreate data slice with data not related to current lap BEGIN:
		if currentLap.TagId != previousLap.TagId || currentLap.RaceId != previousLap.RaceId || currentLap.LapNumber != previousLap.LapNumber {
			// Recreate data slice with old race data not related to current lap :
			otherOldLaps = append(otherOldLaps, previousLap)
		}
		// END.


		// Calculate RacePosition BEGIN:
		currentLap.RacePosition, err = calculateRacePosition(currentLap, otherOldLaps, config.RACE_TYPE)
		if err != nil {
			log.Printf("calculateRacePosition error: %s \n", err)
			return
		} 
		// END.

	}
	// Remove lap duplicates in memory: otherOldLaps + CurrentLap 
	currentLaps = append(otherOldLaps, currentLap)


	// 6. Calculate Race Position:


	// 7. Sort laps by discovery unix time descending (big -> small) depending on config.AVERAGE_RESULTS global setting:
	if config.AVERAGE_RESULTS {
		sortLapsDescByDiscoveryAverageUnixTime(currentLaps)
	} else {
		sortLapsDescByDiscoveryMinimalUnixTime(currentLaps)
	}

	// X. Echo results before return:
	for _, lap := range currentLaps {
		log.Printf("Id=%d, TagId=%s, DiscoveryMinimalUnixTime=%d, DiscoveryAverageUnixTime=%d, AverageResultsCount=%d, RaceId=%d, LapNumber=%d, \n", lap.Id, lap.TagId, lap.DiscoveryMinimalUnixTime, lap.DiscoveryAverageUnixTime, lap.AverageResultsCount, lap.RaceId, lap.LapNumber)
	}

	// 8. Return currentLaps slice or error.
	return
}
