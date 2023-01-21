package Timekeeper

import (
	"log"
	"sort"
	"chicha/pkg/data"
	"chicha/pkg/config"
)

func calculateRaceInMemory (currentTimekeeperTask Data.RawData, previousLaps []Data.Lap, config Config.Settings) (currentLaps []Data.Lap, err error) {

	log.Printf("calculateRaceInMemory received a new job %s\n", currentTimekeeperTask.TagId)

	if cap(previousLaps) == 0 {
		//empty slice
		//	1st - check in db for data from current race:
		var currentLap Data.Lap

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

		currentLaps = append(currentLaps, currentLap)
	} else {


		//non empty slice (race allready running)
		
		//sort slice by DiscoveryUnixTime descending:
		sort.Slice(previousLaps, func(i, j int) bool {
			return previousLaps[i].DiscoveryMinimalUnixTime > previousLaps[j].DiscoveryMinimalUnixTime
		})

		//get current RACE Id
		var currentRaceId uint
		currentRaceId, err = getCurrentRaceId(previousLaps, config)
		if err != nil {
			return
		} else {
			log.Println("currentRaceId:", currentRaceId)
		}


		var currentLap Data.Lap

		currentLap.TagId = currentTimekeeperTask.TagId
		currentLap.DiscoveryMinimalUnixTime = currentTimekeeperTask.DiscoveryUnixTime
		currentLap.DiscoveryAverageUnixTime = currentTimekeeperTask.DiscoveryUnixTime




		  currentLaps = append(previousLaps, currentLap )

	}

	return

}
