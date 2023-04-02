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
		log.Println("Average: sorting laps by DiscoveryAverageUnixTime")
		return lapsToSort[i].DiscoveryAverageUnixTime > lapsToSort[j].DiscoveryAverageUnixTime
	})
}

// Sort slice by discovery minimal unix time descending (big -> small):
func sortLapsDescByDiscoveryMinimalUnixTime (lapsToSort []Data.Lap) {
	sort.Slice(lapsToSort, func(i, j int) bool {
		log.Println("Minimal: sorting laps by DiscoveryMinimalUnixTime")
		return lapsToSort[i].DiscoveryMinimalUnixTime > lapsToSort[j].DiscoveryMinimalUnixTime
	})
}

// Find minimal discovery unix time for each laps with same TagId, RaceId, LapNumber:
func setMinimalDiscoveryUnixTime(laps []Data.Lap) []Data.Lap {

	otherLaps := laps
	for _, lap := range laps {
		for _, otherLap := range otherLaps {
			if lap.TagId == otherLap.TagId && lap.RaceId == otherLap.RaceId && lap.LapNumber == otherLap.LapNumber {
				if lap.DiscoveryMinimalUnixTime > otherLap.DiscoveryMinimalUnixTime {
					// Replace DiscoveryMinimalUnixTime with smaller value:
					lap.DiscoveryMinimalUnixTime = otherLap.DiscoveryMinimalUnixTime
				}
			}
		}
	}

	return laps
}

// Remove duplicates in laps data:
func removeLapDuplicates(elements []Data.Lap) []Data.Lap {
	encountered := map[string]bool{}
	result := []Data.Lap{}

	for _, v := range elements {
		key := fmt.Sprintf("%v:%v:%v", v.TagId, v.RaceId, v.LapNumber)
		if encountered[key] == true {
			// Do not add duplicate element.
		} else {
			// Record this element as encountered.
			encountered[key] = true

			// Find the minimum DiscoveryMinimalUnixTime.
			minDiscoveryMinimalUnixTime := v.DiscoveryMinimalUnixTime
			for _, w := range elements {
				if w.TagId == v.TagId && w.RaceId == v.RaceId && w.LapNumber == v.LapNumber {
					if w.DiscoveryMinimalUnixTime < minDiscoveryMinimalUnixTime {
						minDiscoveryMinimalUnixTime = w.DiscoveryMinimalUnixTime
					}
				}
			}
			v.DiscoveryMinimalUnixTime = minDiscoveryMinimalUnixTime

			// Calculate the average DiscoveryAverageUnixTime.
			var totalDiscoveryAverageUnixTime int64
			encounteredDiscoveryAverageUnixTimes := map[int64]bool{}
			for _, w := range elements {
				if w.TagId == v.TagId && w.RaceId == v.RaceId && w.LapNumber == v.LapNumber {
					totalDiscoveryAverageUnixTime += w.DiscoveryAverageUnixTime
					encounteredDiscoveryAverageUnixTimes[w.DiscoveryAverageUnixTime] = true
				}
			}
			fmt.Println("totalDiscoveryAverageUnixTime=",totalDiscoveryAverageUnixTime)
			fmt.Println("len(encounteredDiscoveryAverageUnixTimes)=", len(encounteredDiscoveryAverageUnixTimes))
			v.DiscoveryAverageUnixTime = totalDiscoveryAverageUnixTime / int64(len(encounteredDiscoveryAverageUnixTimes))
			v.AverageResultsCount = uint(len(encounteredDiscoveryAverageUnixTimes))

			// Append to result slice.
			result = append(result, v)
		}
	}

	// Return the new slice with unique elements.
	return result
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
	currentLaps = setMinimalDiscoveryUnixTime(currentLaps)

	// Remove duplicates
	currentLaps = removeLapDuplicates(currentLaps)

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
