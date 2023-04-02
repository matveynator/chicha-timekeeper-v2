package Timekeeper

import (
	"log"
	"fmt"
	"sort"
	"chicha/pkg/data"
	"chicha/pkg/config"
)

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

// Calculate DiscoveryAverageUnixTime and AverageResultsCount time for each laps with same TagId, RaceId, LapNumber, find and set minimal discovery unix time for each laps with same TagId, RaceId, LapNumber, Remove duplicates in laps data:
func calculateMinimalAndAverageTime(lap Data.Lap, previousLaps []Data.Lap) (updatedLaps []Data.Lap) {

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

		} else {
			// Recreate updated data slice:
			updatedLaps = append(updatedLaps, previousLap)
		}
	}

	//add new data to updated laps:
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

		// Calculate data id: (currentLap.Id)

		// Set well known data receved from rfid (currentLap.TagId, currentLap.DiscoveryMinimalUnixTime, currentLap.DiscoveryAverageUnixTime):
		currentLap.TagId = currentTimekeeperTask.TagId
		currentLap.DiscoveryMinimalUnixTime = currentTimekeeperTask.DiscoveryUnixTime
		currentLap.DiscoveryAverageUnixTime = currentTimekeeperTask.DiscoveryUnixTime

		// Calculate currentRaceId:
		previousLapsCopy := previousLaps
		currentLap.RaceId, err = getCurrentRaceId(currentTimekeeperTask, previousLapsCopy, config)
		if err != nil {
			log.Println(err)
			return
		}

		// Calculate race lap number (currentLap.LapNumber):
		currentLap.LapNumber, err = getMyCurrentRaceLapNumber(currentTimekeeperTask, previousLapsCopy, config)
		if err != nil {
			log.Println(err)
			return
		}

		//currentLap.Id = 
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

	currentLaps = calculateMinimalAndAverageTime(currentLap, previousLaps)

	// Sort laps by discovery unix time descending (big -> small) depending on config.AVERAGE_RESULTS global setting:
	if config.AVERAGE_RESULTS {
		sortLapsDescByDiscoveryAverageUnixTime(currentLaps)
	} else {
		sortLapsDescByDiscoveryMinimalUnixTime(currentLaps)
	}

	// Echo results before return:
	for _, lap := range currentLaps {
		fmt.Printf("TagId=%s, DiscoveryMinimalUnixTime=%d, DiscoveryAverageUnixTime=%d, AverageResultsCount=%d, RaceId=%d, LapNumber=%d, \n", lap.TagId, lap.DiscoveryMinimalUnixTime, lap.DiscoveryAverageUnixTime, lap.AverageResultsCount, lap.RaceId, lap.LapNumber)
	}


	return
}
