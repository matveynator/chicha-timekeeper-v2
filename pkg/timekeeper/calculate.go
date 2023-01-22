package Timekeeper

import (
	"log"
	"sort"
	"chicha/pkg/data"
	"chicha/pkg/config"
)

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
		//non empty slice (race allready running)
		//sort slice by discovery unix time descending (big -> small):
		sort.Slice(previousLaps, func(i, j int) bool {
			if config.AVERAGE_RESULTS {
				log.Println("doing.. config.AVERAGE_RESULTS true")
				return previousLaps[i].DiscoveryAverageUnixTime > previousLaps[j].DiscoveryAverageUnixTime
			} else {
				log.Println("doing.. config.AVERAGE_RESULTS false")
				return previousLaps[i].DiscoveryMinimalUnixTime > previousLaps[j].DiscoveryMinimalUnixTime
			}
		})

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
