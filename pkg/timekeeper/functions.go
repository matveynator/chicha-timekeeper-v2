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




// getCurrentRaceId returns the race ID of the current lap based on previous laps.
// If there are no previous laps, the race ID will be set to 1.
func getCurrentRaceId(currentTimekeeperTask Data.RawData, previousLaps []Data.Lap, config Config.Settings) (raceId uint, err error) {
	if len(previousLaps) == 0 {
		raceId = 1 // If there are no previous laps, assign race ID 1.
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
				// If previous race time expired (timed out), increment race ID:
				log.Println("race expired, starting new race.")
				raceId = previousLaps[0].RaceId + 1
			} else {
				// Return current race ID:
				raceId = previousLaps[0].RaceId
			}
		}
	}
	return
}

// getMyCurrentRaceLapNumber returns the current lap number of the current racer based on previous laps and current timekeeper task.
func getMyCurrentRaceLapNumber(currentTimekeeperTask Data.RawData, previousLaps []Data.Lap, config Config.Settings) (lapNumber uint, err error) {

	var myPreviousLaps []Data.Lap

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


func calculateLapTime(currentLap Data.Lap, otherOldLaps []Data.Lap, config Config.Settings) (lapTime int64, err error) {

	// Create a copy of otherOldLaps
	allLaps := make([]Data.Lap, len(otherOldLaps))
	copy(allLaps, otherOldLaps)

	// Add current lap to allLaps copy:
	allLaps = append(allLaps, currentLap)

	// Sort by time ascending small->big :
	if config.AVERAGE_RESULTS {
		sortLapsAscByDiscoveryAverageUnixTime(allLaps)
	} else {
		sortLapsAscByDiscoveryMinimalUnixTime(allLaps)
	}

	// Get all laps in the same lap number of the same race:
	var sameLapNumberLaps []Data.Lap
	for _, lap := range allLaps {
		if lap.RaceId == currentLap.RaceId && lap.LapNumber == currentLap.LapNumber {
			sameLapNumberLaps = append(sameLapNumberLaps, lap)
		}
	}

	// Sort by time ascending small->big :
	if config.AVERAGE_RESULTS {
		sortLapsAscByDiscoveryAverageUnixTime(sameLapNumberLaps)
	} else {
		sortLapsAscByDiscoveryMinimalUnixTime(sameLapNumberLaps)
	}


	if config.RACE_TYPE == "delayed-start" {
		// First gate pass is a start time, so the first (zero) lap is a zero time for all:
		// Everyone in zero lap number has zero lap time:
		if	currentLap.LapNumber == 0 {
			lapTime = 0 
		} else {
			// Non zero lap number = difference between previous lap time and current lap time:
			for _, allLapData := range allLaps {
				// Get my previous lap data and calculate difference with my current time:
				if allLapData.TagId == currentLap.TagId && allLapData.RaceId == currentLap.RaceId && allLapData.LapNumber == (currentLap.LapNumber - 1) {
					if config.AVERAGE_RESULTS {
						lapTime = currentLap.DiscoveryAverageUnixTime - allLapData.DiscoveryAverageUnixTime 
					} else {
						lapTime = currentLap.DiscoveryMinimalUnixTime - allLapData.DiscoveryMinimalUnixTime
					}
				}
			}
		}
	} else if config.RACE_TYPE == "mass-start" {
		// First gate pass is a first small lap, so the first (zero) lap is the difference time between the first gate pass (the leader) and the current pass time. Only leader has zero lap time in zero lap number.

		// Get zero lap leader time:
		// If lap number is zero - calculcate my lap time as difference between my current time and leader time:
		var leaderTime int64 
		if config.AVERAGE_RESULTS {
			leaderTime = sameLapNumberLaps[0].DiscoveryAverageUnixTime
			lapTime = currentLap.DiscoveryAverageUnixTime - leaderTime
		} else {
			leaderTime = sameLapNumberLaps[0].DiscoveryMinimalUnixTime
			lapTime = currentLap.DiscoveryMinimalUnixTime - leaderTime
		}

		for _, allLapData := range allLaps { 
			// For laps greater than zero calculate lap time as difference between my current time and my previous time: 
			// Get my previous lap data and calculate difference with my current time:
			if currentLap.LapNumber > 0 && allLapData.TagId == currentLap.TagId && allLapData.RaceId == currentLap.RaceId && allLapData.LapNumber == (currentLap.LapNumber - 1) {
				if config.AVERAGE_RESULTS {
					lapTime = currentLap.DiscoveryAverageUnixTime - allLapData.DiscoveryAverageUnixTime
				} else {
					lapTime = currentLap.DiscoveryMinimalUnixTime - allLapData.DiscoveryMinimalUnixTime
				}
			}
		}
	}

	return
}


// Функция которая отсортирует слайс и для каждой уникальной группы данных с определенным RaceId выберет все LapNumber больше нуля и выберет среди них минимальное  значение поле LapTime и поставит для этой структуры FastestLapInThisRace true, а всем остальным  из тойже группы одинаковых RaceId поставит FastestLapInThisRace false
func setFastestLapInEachRace(laps []Data.Lap) {
	// Сортировка по RaceId и LapNumber
	sort.Slice(laps, func(i, j int) bool {
		if laps[i].RaceId == laps[j].RaceId {
			return laps[i].LapNumber < laps[j].LapNumber
		}
		return laps[i].RaceId < laps[j].RaceId
	})

	var currentRaceId uint
	var currentGroup []Data.Lap
	for _, lap := range laps {
		// Если это новая группа, обрабатываем предыдущую
		if lap.RaceId != currentRaceId {
			setFastestLapInGroup(currentGroup)
			currentRaceId = lap.RaceId
			currentGroup = []Data.Lap{}
		}
		currentGroup = append(currentGroup, lap)
	}
	// Обработка последней группы
	setFastestLapInGroup(currentGroup)
}

// setFastestLapInGroup устанавливает поле FastestLapInThisRace
// для группы laps в соответствии с условиями
func setFastestLapInGroup(laps []Data.Lap) {
	// Создаем карту для хранения минимального LapTime для каждого LapNumber
	minLapTimes := make(map[uint]int64)
	for _, lap := range laps {
		if lap.LapNumber > 0 {
			if lapTime, ok := minLapTimes[lap.LapNumber]; !ok || lap.LapTime < lapTime {
				minLapTimes[lap.LapNumber] = lap.LapTime
			}
		}
	}
	// Устанавливаем поле FastestLapInThisRace в соответствии с найденным минимальным LapTime
	for _, lap := range laps {
		if lap.LapNumber > 0 {
			if lap.LapTime == minLapTimes[lap.LapNumber] {
				lap.FastestLapInThisRace = true
			} else {
				lap.FastestLapInThisRace = false
			}
		}
	}
}


//  Подсчитываем мое лучшее время круга:
func calculateMyBestLapTime(currentLap Data.Lap, otherOldLaps []Data.Lap) (bestLapTime int64, bestLapNumber uint, err error) {
	if	currentLap.LapNumber == 0 {
		bestLapTime = 0
		bestLapNumber = 0
	} else if currentLap.LapNumber == 1 {
		bestLapTime = currentLap.LapTime
		bestLapNumber = 1
	} else {
		// Make a selection that includes only the complete laps that belong to me:
		var mySameRaceFullLaps []Data.Lap
		for _, otherOldLap := range otherOldLaps {
			if otherOldLap.RaceId == currentLap.RaceId &&  otherOldLap.TagId == currentLap.TagId && otherOldLap.LapNumber > 0  {
				mySameRaceFullLaps = append(mySameRaceFullLaps, otherOldLap)
			}
		}
		mySameRaceFullLaps = append(mySameRaceFullLaps, currentLap)

		// Sort by lap time ascending small->big :
		sortLapsAscLapTime(mySameRaceFullLaps)

		// Assign fastest lap to bestLapTime and bestLapNumber:
		bestLapTime = mySameRaceFullLaps[0].LapTime
		bestLapNumber = mySameRaceFullLaps[0].LapNumber
	}
	return
}

// Сортировка по максимальному LapNumber и минимальному RaceTotalTime:
func sortMaxLapNumberAndMinRaceTotalTime(laps []Data.Lap) []Data.Lap {
	sort.Slice(laps, func(i, j int) bool {
		if laps[i].LapNumber > laps[j].LapNumber {
			return true
		} else if laps[i].LapNumber < laps[j].LapNumber {
			return false
		} else {
			return laps[i].RaceTotalTime < laps[j].RaceTotalTime
		}
	})
	return laps
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

// Check slice contains TagID:
func containsTagId(laps []Data.Lap, tagId string) bool {
	for _, lap := range laps {
		if lap.TagId == tagId {
			return true
		}
	}
	return false
}


// Sort slice by lap time ascending (small -> big):
func sortLapsAscLapTime (lapsToSort []Data.Lap) {
	sort.Slice(lapsToSort, func(i, j int) bool {
		return lapsToSort[i].LapTime < lapsToSort[j].LapTime
	})
}

// Sort slice by race total time descending (big -> small):
func sortLapsDescByRaceTotalTime (lapsToSort []Data.Lap) {
	sort.Slice(lapsToSort, func(i, j int) bool {
		return lapsToSort[i].RaceTotalTime > lapsToSort[j].RaceTotalTime
	})
}

// Sort slice by race total time ascending (small -> big):
func sortLapsAscByRaceTotalTime (lapsToSort []Data.Lap) {
	sort.Slice(lapsToSort, func(i, j int) bool {
		return lapsToSort[i].RaceTotalTime < lapsToSort[j].RaceTotalTime
	})
}

// Sort slice by discovery average unix time descending (big -> small):
func sortLapsDescByDiscoveryAverageUnixTime (lapsToSort []Data.Lap) {
	sort.Slice(lapsToSort, func(i, j int) bool {
		return lapsToSort[i].DiscoveryAverageUnixTime > lapsToSort[j].DiscoveryAverageUnixTime
	})
}

// Sort slice by discovery average unix time ascending (small -> big):
func sortLapsAscByDiscoveryAverageUnixTime (lapsToSort []Data.Lap) {
	sort.Slice(lapsToSort, func(i, j int) bool {
		return lapsToSort[i].DiscoveryAverageUnixTime < lapsToSort[j].DiscoveryAverageUnixTime
	})
}


// Sort slice by discovery minimal unix time descending (big -> small):
func sortLapsDescByDiscoveryMinimalUnixTime (lapsToSort []Data.Lap) {
	sort.Slice(lapsToSort, func(i, j int) bool {
		return lapsToSort[i].DiscoveryMinimalUnixTime > lapsToSort[j].DiscoveryMinimalUnixTime
	})
}

// Sort slice by discovery minimal unix time ascending (small -> big):
func sortLapsAscByDiscoveryMinimalUnixTime (lapsToSort []Data.Lap) {
	sort.Slice(lapsToSort, func(i, j int) bool {
		return lapsToSort[i].DiscoveryMinimalUnixTime < lapsToSort[j].DiscoveryMinimalUnixTime
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
		currentLap.LapIsLatest = true
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


	// 6. Sort laps by discovery unix time descending (big -> small) depending on config.AVERAGE_RESULTS global setting:
	if config.AVERAGE_RESULTS {
		sortLapsDescByDiscoveryAverageUnixTime(currentLaps)
	} else {
		sortLapsDescByDiscoveryMinimalUnixTime(currentLaps)
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
		setFastestLapInEachRace(currentLaps)
	}
	// END.


	// X. Echo results before return:
	for _, lap := range currentLaps {
		log.Printf("Id=%d, TagId=%s, DiscoveryMinimalUnixTime=%d, DiscoveryAverageUnixTime=%d, AverageResultsCount=%d, RaceId=%d, LapNumber=%d, LapTime=%d, RaceTotalTime=%d, RacePosition=%d, LapPosition=%d TimeBehindTheLeader=%d, LapIsLatest=%t FasterOrSlowerThanPreviousLapTime=%d LapIsStrange=%t BestLapTime=%d, BestLapNumber=%d, FastestLapInThisRace=%t\n", lap.Id, lap.TagId, lap.DiscoveryMinimalUnixTime, lap.DiscoveryAverageUnixTime, lap.AverageResultsCount, lap.RaceId, lap.LapNumber, lap.LapTime, lap.RaceTotalTime, lap.RacePosition, lap.LapPosition, lap.TimeBehindTheLeader, lap.LapIsLatest, lap.FasterOrSlowerThanPreviousLapTime, lap.LapIsStrange, lap.BestLapTime, lap.BestLapNumber, lap.FastestLapInThisRace)
	}

	// 8. Return currentLaps slice or error.
	return
}
