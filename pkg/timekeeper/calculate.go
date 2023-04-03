package Timekeeper

import (
	"log"
	"fmt"
	"sort"
	"time"
	"chicha/pkg/data"
	"chicha/pkg/config"
)

// Check data is valid:
func checkLapIsValid(currentTimekeeperTask Data.RawData, previousLaps []Data.Lap, config Config.Settings) (isValid bool, err error) {
	if len(previousLaps) == 0 {
		isValid = true
		return
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
				// If previous race time expired (timed out), data is valid:
				log.Println("race expired, starting new race.")
				isValid = true
				return
			} else {
				// If race not expired, calculate lap time (valid if > minimal lap time before previous results) and not more than precision time.

				// Create slice with only my data in it:
				var myPreviousLaps []Data.Lap
				for _, myPreviousLap := range previousLapsCopy {
					if currentTimekeeperTask.TagId == myPreviousLap.TagId {
						// Add my previous data to slice:	
						myPreviousLaps = append(myPreviousLaps, myPreviousLap)
					}
				}

				if len(myPreviousLaps) > 0 {
					// Sort laps by discovery unix time descending (big -> small) depending on config.AVERAGE_RESULTS global setting:
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
					if currentTime.After(myLastLapTime.Add(config.AVERAGE_DURATION)) && currentTime.Before(myLastLapTime.Add(config.MINIMAL_LAP_TIME_DURATION)) {
						// Time is more than config.AVERAGE_DURATION and less than config.MINIMAL_LAP_TIME_DURATION - is invalid:
						isValid = false
						return
					} else {
						isValid = true
						return
					}	
				} else {
					// The first race and lap is always valid:
					isValid = true
					return
				}

			}
		}
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

// Calculate Id, DiscoveryMinimalUnixTime, DiscoveryAverageUnixTime, AverageResultsCount  for each laps with same TagId, RaceId, LapNumber. Remove duplicates.
func calculateIdMinimalAverageTimeRemoveDupes(lap Data.Lap, previousLaps []Data.Lap) (updatedLaps []Data.Lap) {

	for _, previousLap := range previousLaps {
		if lap.TagId == previousLap.TagId && lap.RaceId == previousLap.RaceId && lap.LapNumber == previousLap.LapNumber {
			if lap.DiscoveryMinimalUnixTime > previousLap.DiscoveryMinimalUnixTime {
				// Replace DiscoveryMinimalUnixTime with smaller value:
				lap.DiscoveryMinimalUnixTime = previousLap.DiscoveryMinimalUnixTime
			}
			// Calculate DiscoveryAverageUnixTime and AverageResultsCount:
			lap.DiscoveryAverageUnixTime = (lap.DiscoveryAverageUnixTime + previousLap.DiscoveryAverageUnixTime)/2

			// Set first AverageResultsCount = 1:
			if previousLap.AverageResultsCount == 0 {
				previousLap.AverageResultsCount = 1
			}
			lap.AverageResultsCount = previousLap.AverageResultsCount + 1


			// Calculate Id:
			if previousLap.Id > 0 {
				// Set old Id if available:
				lap.Id = previousLap.Id
			} else {
				// Calculate next Id if data is new:
				lap.Id = getNextId (previousLaps)
			}


		} else {
			// Recreate data slice with old data not related to current:
			updatedLaps = append(updatedLaps, previousLap)
		}
	}

	updatedLaps = append(updatedLaps, lap)

	return updatedLaps
}

func calculateRaceInMemory (currentTimekeeperTask Data.RawData, previousLaps []Data.Lap, config Config.Settings) (currentLaps []Data.Lap, err error) {

	var currentLap Data.Lap


	//log.Printf("calculateRaceInMemory received a new job %s\n", currentTimekeeperTask.TagId)

	if len(previousLaps) == 0 {
		//empty slice
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

		var dataIsValid bool
		//	1. Check if current lap data is within lap timeouts and race timeouts and combine with the rest of the race data:
		dataIsValid, err = checkLapIsValid(currentTimekeeperTask, previousLaps, config) 
		if err != nil {
			log.Printf("Data validation error: %s \n", err)
			return
		} else {
			if !dataIsValid {
				// Data is invalid - drop it and return previousLaps:
				log.Println("Received invalid lap data: Time is after AVERAGE_DURATION and before MINIMAL_LAP_TIME_DURATION. Skipping it.")
				currentLaps = previousLaps
				return
			}
		}

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


		// currentLap.Id = 
		//currentLap.TagId = 
		//currentLap.DiscoveryMinimalUnixTime = 
		//currentLap.DiscoveryAverageUnixTime = 
		//currentLap.AverageResultsCount = 1
		//currentLap.RaceId = 
		//currentLap.RacePosition = 
		//currentLap.TimeBehindTheLeader = 
		//currentLap.LapNumber =
		//currentLap.LapPosition = 1
		//currentLap.LapIsCurrent = true
		//currentLap.LapIsStrange = false
		//currentLap.RaceFinished = false
		//currentLap.RaceTotalTime = 0

	}
	// 5. Calculate Id, DiscoveryMinimalUnixTime, DiscoveryAverageUnixTime, AverageResultsCount  for each laps with same TagId, RaceId, LapNumber. Remove duplicates.
	currentLaps = calculateIdMinimalAverageTimeRemoveDupes(currentLap, previousLaps)

	// 6. Sort laps by discovery unix time descending (big -> small) depending on config.AVERAGE_RESULTS global setting:
	if config.AVERAGE_RESULTS {
		sortLapsDescByDiscoveryAverageUnixTime(currentLaps)
	} else {
		sortLapsDescByDiscoveryMinimalUnixTime(currentLaps)
	}

	// 7. Echo results before return:
	for _, lap := range currentLaps {
		fmt.Printf("Id=%d, TagId=%s, DiscoveryMinimalUnixTime=%d, DiscoveryAverageUnixTime=%d, AverageResultsCount=%d, RaceId=%d, LapNumber=%d, \n", lap.Id, lap.TagId, lap.DiscoveryMinimalUnixTime, lap.DiscoveryAverageUnixTime, lap.AverageResultsCount, lap.RaceId, lap.LapNumber)
	}

	// 8. Return currentLaps slice or error.
	return
}
