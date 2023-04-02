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

// Find and set minimal discovery unix time for each laps with same TagId, RaceId, LapNumber:
func calculateMinimalTime(laps []Data.Lap) (minimalLaps []Data.Lap) {

	for _, lap := range laps {
		for _, otherLap := range laps {
			if lap.TagId == otherLap.TagId && lap.RaceId == otherLap.RaceId && lap.LapNumber == otherLap.LapNumber {
				if lap.DiscoveryMinimalUnixTime > otherLap.DiscoveryMinimalUnixTime {
					// Replace DiscoveryMinimalUnixTime with smaller value:
					lap.DiscoveryMinimalUnixTime = otherLap.DiscoveryMinimalUnixTime
				}
			}
		}
		minimalLaps = append(minimalLaps, lap)
	}

	return minimalLaps
}


// Calculate DiscoveryAverageUnixTime and AverageResultsCount time for each laps with same TagId, RaceId, LapNumber:
func calculateAverageTimeAndResultsCount(laps []Data.Lap) (averageLaps []Data.Lap) {

	// Create slice copy:
	for _, lap := range laps {
		// Set initial AverageResultsCount equals zero
		lap.AverageResultsCount = 0
		for _, otherLap := range laps {
			if lap.TagId == otherLap.TagId && lap.RaceId == otherLap.RaceId && lap.LapNumber == otherLap.LapNumber {
				// Calculate DiscoveryAverageUnixTime and AverageResultsCount:	
				lap.DiscoveryAverageUnixTime = (lap.DiscoveryAverageUnixTime + otherLap.DiscoveryAverageUnixTime)/2
				lap.AverageResultsCount = lap.AverageResultsCount + 1
			}
		}
		// Prepare slice to return:
		averageLaps = append(averageLaps, lap)
	}

	return averageLaps
}

// Remove duplicates in laps data:
func removeDuplicates(laps []Data.Lap) (uniqLaps []Data.Lap) {

	for lapIndex, lap := range laps {
		for otherLapIndex, otherLap := range laps {
			if lap.TagId == otherLap.TagId && lap.RaceId == otherLap.RaceId && lap.LapNumber == otherLap.LapNumber && lap.DiscoveryAverageUnixTime == otherLap.DiscoveryAverageUnixTime && lap.DiscoveryMinimalUnixTime == otherLap.DiscoveryMinimalUnixTime {
				if lapIndex == otherLapIndex {
					// Append to uniq lap slice:
					uniqLaps = append(uniqLaps, lap)
				}
			}
		}
	}

	// Return the new slice with unique elements.
	return uniqLaps
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

		// 1: calculate id

		// 2: calculate race id

		// 3: calcutale lap number

	}
	// Add currentLap to previousLaps slice:
	currentLaps = append(previousLaps, currentLap)


	// Find minimal discovery unix time for each laps with same TagId, RaceId, LapNumber:
	currentLaps = calculateMinimalTime(currentLaps)

	// Calculate DiscoveryAverageUnixTime and AverageResultsCount time for each laps with same TagId, RaceId, LapNumber:
	currentLaps = calculateAverageTimeAndResultsCount(currentLaps)

	// Remove duplicates
	currentLaps = removeDuplicates(currentLaps)

	// Sort laps by discovery unix time descending (big -> small) depending on config.AVERAGE_RESULTS global setting:
	if config.AVERAGE_RESULTS {
		sortLapsDescByDiscoveryAverageUnixTime(currentLaps)
	} else {
		sortLapsDescByDiscoveryMinimalUnixTime(currentLaps)
	}


	// Echo results before return:
	for _, currentLap := range currentLaps {
		fmt.Printf("TagId=%s, DiscoveryMinimalUnixTime=%d, DiscoveryAverageUnixTime=%d, AverageResultsCount=%d, RaceId=%d, LapNumber=%d, \n", currentLap.TagId, currentLap.DiscoveryMinimalUnixTime, currentLap.DiscoveryAverageUnixTime, currentLap.AverageResultsCount, currentLap.RaceId, currentLap.LapNumber)
	}

	return
}
