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

		//remove duplicates
		previousLaps = removeLapDuplicates(previousLaps)

		// Sort laps by discovery unix time descending (big -> small):
		if config.AVERAGE_RESULTS {
			sortLapsDescByDiscoveryAverageUnixTime(previousLaps)
		} else {
			sortLapsDescByDiscoveryMinimalUnixTime(previousLaps)
		}

		for _, previousLap := range previousLaps {
			fmt.Printf("TagId=%s, DiscoveryMinimalUnixTime=%d, DiscoveryAverageUnixTime=%d, AverageResultsCount=%d, RaceId=%d, LapNumber=%d \n", previousLap.TagId, previousLap.DiscoveryMinimalUnixTime, previousLap.AverageResultsCount, previousLap.RaceId, previousLap.LapNumber)

		}

		//get current race Id
		var currentRaceId uint
		previousLapsCopy := previousLaps
		currentRaceId, err = getCurrentRaceId(currentTimekeeperTask, previousLapsCopy, config)
		if err != nil {
			log.Println(err)
			return
		} else {
			log.Println("currentRaceId:", currentRaceId)

			//get my current race lap number
			var myCurrentRaceLapNumber uint
			myCurrentRaceLapNumber, err = getMyCurrentRaceLapNumber(currentTimekeeperTask, previousLapsCopy, config)
			if err != nil {
				log.Println(err)
				return
			} else {
				log.Println("myCurrentRaceLapNumber:", myCurrentRaceLapNumber)

				//currentLap.Id = 1
				currentLap.TagId = currentTimekeeperTask.TagId
				currentLap.DiscoveryMinimalUnixTime = currentTimekeeperTask.DiscoveryUnixTime
				currentLap.DiscoveryAverageUnixTime = currentTimekeeperTask.DiscoveryUnixTime
				//currentLap.AverageResultsCount = 1
				currentLap.RaceId = currentRaceId
				//currentLap.RacePosition = 1
				//currentLap.TimeBehindTheLeader = 0
				currentLap.LapNumber = myCurrentRaceLapNumber
				//currentLap.LapPosition = 1
				//currentLap.LapIsCurrent = true
				//currentLap.LapIsStrange = false
				//currentLap.RaceFinished = false
				//currentLap.RaceTotalTime = 0
			}
		}
	}

	//add currentLap to currentLaps slice:
	currentLaps = append(previousLaps, currentLap)
	return
}
