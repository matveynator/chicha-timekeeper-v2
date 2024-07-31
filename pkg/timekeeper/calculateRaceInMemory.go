package Timekeeper

import (
	"log"
	"chicha/pkg/data"
	"chicha/pkg/config"
)

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
		currentLap.LapIsLatest = true
		currentLap.RaceFinished = false
		currentLap.RaceTotalTime = 0
	} else {
		// Not an empty slice (race allready running)

		// 2. Set well known data receved from rfid (currentLap.TagId, currentLap.DiscoveryMinimalUnixTime, currentLap.DiscoveryAverageUnixTime):
		currentLap.TagId = currentTimekeeperTask.TagId
		currentLap.DiscoveryMinimalUnixTime = currentTimekeeperTask.DiscoveryUnixTime
		currentLap.DiscoveryMaximalUnixTime = currentTimekeeperTask.DiscoveryUnixTime
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
	}


	// Calculate Lap.Id START:
	currentLap.Id = calculateLapId(currentLap, previousLaps)
	// END.

	// 5. Calculate DiscoveryMinimalUnixTime, DiscoveryAverageUnixTime, AverageResultsCount  for each laps with same TagId, RaceId, LapNumber. Remove duplicates BEGIN:.
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

			// Calculate DiscoveryMaximalUnixTime BEGIN:
			if currentLap.DiscoveryMaximalUnixTime < previousLap.DiscoveryMaximalUnixTime {
				currentLap.DiscoveryMaximalUnixTime = previousLap.DiscoveryMinimalUnixTime
			}
			// END.

			// Calculate DiscoveryAverageUnixTime BEGIN:
			currentLap.DiscoveryAverageUnixTime = (currentTimekeeperTask.DiscoveryUnixTime + previousLap.DiscoveryAverageUnixTime)/2
			// END.

			// Calculate AverageResultsCount BEGIN:
			if previousLap.AverageResultsCount == 0 {
				previousLap.AverageResultsCount = 1
			} 
			currentLap.AverageResultsCount = previousLap.AverageResultsCount + 1
			// END.
		}


		// Rewrite my lap data in memory with updated new data (DiscoveryMinimalUnixTime, DiscoveryAverageUnixTime, AverageResultsCount) removing duplicates END.
		// Recreate data slice with data not related to current lap BEGIN:
		if currentLap.TagId != previousLap.TagId || currentLap.RaceId != previousLap.RaceId || currentLap.LapNumber != previousLap.LapNumber {
			// Recreate data slice with old race data not related to current lap :
			otherOldLaps = append(otherOldLaps, previousLap)
		}
		// END.

	}

	// Calculate LapTime BEGIN:
	currentLap.LapTime, err = calculateLapTime(currentLap, otherOldLaps, config)
	if err != nil {
		log.Printf("calculateLapTime error: %s \n", err)
		return
	} 
	// END.


	// Calculate my BestLapTime and BestLapNumber BEGIN:
	if config.VARIABLE_DISTANCE_RACE == false {
		currentLap.BestLapTime, currentLap.BestLapNumber, err = calculateMyBestLapTime(currentLap, otherOldLaps)
		if err != nil {
			log.Printf("calculateMyBestLapTime error: %s \n", err)
			return
		}
	}
	// END.

	// Calculate RaceTotalTime BEGIN:
	// For zero lap RaceTotalTime equals LapTime:
	if currentLap.LapNumber == 0 {
		currentLap.RaceTotalTime = currentLap.LapTime
		currentLap.FasterOrSlowerThanPreviousLapTime = 0
		currentLap.LapIsStrange = false
	} else {
		// RaceTotalTime equals my previous lap RaceTotalTime + my current LapTime:
		for _, otherOldLap := range otherOldLaps {
			if currentLap.TagId == otherOldLap.TagId && currentLap.RaceId == otherOldLap.RaceId && otherOldLap.LapNumber == currentLap.LapNumber-1 {	
				// My previous lap is otherOldLap:
				currentLap.RaceTotalTime = otherOldLap.RaceTotalTime + currentLap.LapTime
				// Calculate if my current lap is faster or slower than previous lap only from 2nd lap result: 
				if currentLap.LapNumber > 1 && config.VARIABLE_DISTANCE_RACE == false {
					currentLap.FasterOrSlowerThanPreviousLapTime = currentLap.LapTime - otherOldLap.LapTime
					// Проверяем если текущий круг длиннее предыдущего в два раза то ставим маркер "LapIsStrange" (подозрительный круг) - возможно спортсмен упал или его данные не считались RFID считывателем - нужно сверить данные в других источниках (фотофиниш или круговые судьи).
					if (currentLap.LapTime - otherOldLap.LapTime) >= otherOldLap.LapTime {
						currentLap.LapIsStrange = true 
					} else  {
						currentLap.LapIsStrange = false
					}
				} else {
					currentLap.FasterOrSlowerThanPreviousLapTime = 0
					currentLap.LapIsStrange = false
				}
			}
		}
	}
	// END.



	// Remove lap duplicates in memory: otherOldLaps + CurrentLap 
	currentLaps = append(otherOldLaps, currentLap)


	// 6. Sort laps by discovery unix time ascending (small -> big) depending on config.AVERAGE_RESULTS global setting:
	if config.AVERAGE_RESULTS {
		sortLapsAscByDiscoveryAverageUnixTime(currentLaps)
	} else {
		sortLapsAscByDiscoveryMinimalUnixTime(currentLaps)
	}

	// 7. Calculate Race Position and Lap Position BEGIN:

	// Create current race slice only with one latest lap data per rider:
	var currentRaceLatestLaps []Data.Lap

	for index, currentInMemoryLap := range currentLaps {
		// Create current race slice:
		if currentInMemoryLap.RaceId == currentLap.RaceId {
			// Add only one latest lap data per rider:
			if !containsTagId(currentRaceLatestLaps, currentInMemoryLap.TagId) {
				// Mark latest lap as latest:
				currentLaps[index].LapIsLatest = true
				// Add latest lap to currentRaceLatestLaps slice:  
				currentRaceLatestLaps = append(currentRaceLatestLaps, currentInMemoryLap)
			} else {
				// Mark not latest laps as not latest:
				currentLaps[index].LapIsLatest = false
			}
		}
	}

	//Сортируе по максимальному LapNumber и минимальному RaceTotalTime:
	currentRaceLatestLaps = sortMaxLapNumberAndMinRaceTotalTime(currentRaceLatestLaps)

	// Slice to store leader of the race current lap:
	var leaderCurrentLap Data.Lap

	// Calculate position in currentRaceLatestLaps:
	for latestPosition, latestLapData := range currentRaceLatestLaps {

		// Calculate the leader of the race current lap: 
		if latestPosition == 0 {
			leaderCurrentLap = latestLapData
		}

		// Set RacePosition, LapPosition and TimeBehindTheLeader in global currentLaps slice:
		for index, lap := range currentLaps {
			if lap.Id == latestLapData.Id {

				if config.RACE_TYPE == "mass-start" {
					currentLaps[index].RacePosition = uint(latestPosition)+1
					currentLaps[index].LapPosition = uint(latestPosition)+1
				} else if config.RACE_TYPE == "delayed-start" {
					if lap.LapNumber == 0 {
						// Zero lap in "delayed-start" race is always zero:
						currentLaps[index].RacePosition = 0
						currentLaps[index].LapPosition = 0
					} else {
						currentLaps[index].RacePosition = uint(latestPosition)+1
						currentLaps[index].LapPosition = uint(latestPosition)+1
					}
				}
				// Calculate TimeBehindTheLeader:
				currentLaps[index].TimeBehindTheLeader = lap.RaceTotalTime - leaderCurrentLap.RaceTotalTime
			}
		}
	}
	// End.


	// 8. Calculate fastest lap in this race BEGIN:
	if config.VARIABLE_DISTANCE_RACE == false {
		setFastestLaps(currentLaps)
	}
	// END.


	// X. Echo results before return:
	i := 0
	for _, lap := range currentRaceLatestLaps {
		if i == 0 {
			 log.Printf("Race: %d, Lap: %d\n", lap.RaceId, lap.LapNumber)
		}
		i++

		//log.Printf("Id=%d, TagId=%s, DiscoveryMinimalUnixTime=%d, DiscoveryMaximalUnixTime=%d, DiscoveryAverageUnixTime=%d, AverageResultsCount=%d, RaceId=%d, LapNumber=%d, LapTime=%d, RaceTotalTime=%d, RacePosition=%d, LapPosition=%d TimeBehindTheLeader=%d, LapIsLatest=%t FasterOrSlowerThanPreviousLapTime=%d LapIsStrange=%t BestLapTime=%d, BestLapNumber=%d, FastestLapInThisRace=%t\n", lap.Id, lap.TagId, lap.DiscoveryMinimalUnixTime, lap.DiscoveryMaximalUnixTime, lap.DiscoveryAverageUnixTime, lap.AverageResultsCount, lap.RaceId, lap.LapNumber, lap.LapTime, lap.RaceTotalTime, lap.RacePosition, lap.LapPosition, lap.TimeBehindTheLeader, lap.LapIsLatest, lap.FasterOrSlowerThanPreviousLapTime, lap.LapIsStrange, lap.BestLapTime, lap.BestLapNumber, lap.FastestLapInThisRace)
			//if lap.LapIsLatest == true {
				log.Printf("%d, %s, %d\n", lap.RacePosition,lap.TagId, lap.TimeBehindTheLeader)
			//}
	}

	// 8. Return currentLaps slice or error.
	return
}
