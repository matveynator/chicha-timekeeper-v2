package Timekeeper

import (
	"chicha/pkg/config"
	"chicha/pkg/data"
	"fmt"
	"strings"
)

// Обрабатывает данные гонки в памяти и возвращает обновленные круги и возможную ошибку.
func calculateRaceInMemory(currentTimekeeperTask Data.RawData, previousLaps []Data.Lap, config Config.Settings) (updatedLaps []Data.Lap, err error) {
	var currentLap Data.Lap
	var lapIsValid bool

	// Проверка валидности текущих данных круга.
	lapIsValid, err = checkLapIsValid(currentTimekeeperTask, previousLaps, config)
	if err != nil {
		fmt.Printf("Ошибка валидации данных: %s \n", err)
		return
	} else if !lapIsValid {
		fmt.Println("Получены невалидные данные круга. Пропускаем.")
		updatedLaps = previousLaps
		return
	}

	// Инициализация текущего круга
	currentLap, err = initializeCurrentLap(currentTimekeeperTask, previousLaps, config)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Обновление текущего круга на основе предыдущих данных
	updatedLaps, err = updateCurrentLap(currentLap, previousLaps, currentTimekeeperTask, config)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Рассчитываем позиции и отставания
	calculatePositionsAndDelays(updatedLaps, config)

	// Логируем результаты
	clearScreen()
	logResults(updatedLaps)

	return updatedLaps, nil
}

// Инициализация текущего круга
func initializeCurrentLap(currentTimekeeperTask Data.RawData, previousLaps []Data.Lap, config Config.Settings) (Data.Lap, error) {
	var currentLap Data.Lap
	var err error

	if len(previousLaps) == 0 {
		currentLap = Data.Lap{
			Id:                        1,
			TagId:                     currentTimekeeperTask.TagId,
			DiscoveryMinimalUnixTime:  currentTimekeeperTask.DiscoveryUnixTime,
			DiscoveryAverageUnixTime:  currentTimekeeperTask.DiscoveryUnixTime,
			AverageResultsCount:       1,
			RaceId:                    1,
			RacePosition:              1,
			TimeBehindTheLeader:       0,
			LapNumber:                 0,
			LapPosition:               1,
			LapIsLatest:               true,
			RaceFinished:              false,
			RaceTotalTime:             0,
		}
	} else {
		currentLap = Data.Lap{
			TagId:                     currentTimekeeperTask.TagId,
			DiscoveryMinimalUnixTime:  currentTimekeeperTask.DiscoveryUnixTime,
			DiscoveryMaximalUnixTime:  currentTimekeeperTask.DiscoveryUnixTime,
			DiscoveryAverageUnixTime:  currentTimekeeperTask.DiscoveryUnixTime,
		}

		currentLap.RaceId, err = getCurrentRaceId(currentTimekeeperTask, previousLaps, config)
		if err != nil {
			return currentLap, err
		}

		currentLap.LapNumber, err = getMyCurrentRaceLapNumber(currentTimekeeperTask, previousLaps, config)
		if err != nil {
			return currentLap, err
		}
	}

	currentLap.Id = calculateLapId(currentLap, previousLaps)
	return currentLap, nil
}

// Обновление текущего круга на основе предыдущих данных
func updateCurrentLap(currentLap Data.Lap, previousLaps []Data.Lap, currentTimekeeperTask Data.RawData, config Config.Settings) ([]Data.Lap, error) {
	var nonMatchingLaps []Data.Lap
	var err error

	for _, previousLap := range previousLaps {
		if currentLap.TagId == previousLap.TagId && currentLap.RaceId == previousLap.RaceId && currentLap.LapNumber == previousLap.LapNumber {
			updateLapTimes(&currentLap, previousLap, currentTimekeeperTask)
		} else {
			nonMatchingLaps = append(nonMatchingLaps, previousLap)
		}
	}

	currentLap.LapTime, err = calculateLapTime(currentLap, nonMatchingLaps, config)
	if err != nil {
		return nil, fmt.Errorf("Ошибка в расчете времени круга: %s", err)
	}

	if !config.VARIABLE_DISTANCE_RACE {
		currentLap.BestLapTime, currentLap.BestLapNumber, err = calculateMyBestLapTime(currentLap, nonMatchingLaps)
		if err != nil {
			return nil, fmt.Errorf("Ошибка в расчете лучшего времени круга: %s", err)
		}
	}

	calculateRaceTotalTime(&currentLap, nonMatchingLaps, config)

	updatedLaps := append(nonMatchingLaps, currentLap)
	sortLaps(updatedLaps, config)
	return updatedLaps, nil
}

// Обновление времен круга
func updateLapTimes(currentLap *Data.Lap, previousLap Data.Lap, currentTimekeeperTask Data.RawData) {
	if currentLap.DiscoveryMinimalUnixTime > previousLap.DiscoveryMinimalUnixTime {
		currentLap.DiscoveryMinimalUnixTime = previousLap.DiscoveryMinimalUnixTime
	}

	if currentLap.DiscoveryMaximalUnixTime < previousLap.DiscoveryMaximalUnixTime {
		currentLap.DiscoveryMaximalUnixTime = previousLap.DiscoveryMaximalUnixTime
	}

	currentLap.DiscoveryAverageUnixTime = (currentTimekeeperTask.DiscoveryUnixTime + previousLap.DiscoveryAverageUnixTime) / 2

	if previousLap.AverageResultsCount < 1 {
		previousLap.AverageResultsCount = 1
	}

	currentLap.AverageResultsCount = previousLap.AverageResultsCount + 1
}

// Рассчитываем общее время гонки для текущего круга
func calculateRaceTotalTime(currentLap *Data.Lap, nonMatchingLaps []Data.Lap, config Config.Settings) {
	if currentLap.LapNumber == 0 {
		currentLap.RaceTotalTime = currentLap.LapTime
		currentLap.FasterOrSlowerThanPreviousLapTime = 0
		currentLap.LapIsStrange = false
	} else {
		for _, lap := range nonMatchingLaps {
			if currentLap.TagId == lap.TagId && currentLap.RaceId == lap.RaceId && lap.LapNumber == currentLap.LapNumber-1 {
				currentLap.RaceTotalTime = lap.RaceTotalTime + currentLap.LapTime

				if currentLap.LapNumber > 1 && !config.VARIABLE_DISTANCE_RACE {
					currentLap.FasterOrSlowerThanPreviousLapTime = currentLap.LapTime - lap.LapTime

					if (currentLap.LapTime - lap.LapTime) >= lap.LapTime {
						currentLap.LapIsStrange = true
					} else {
						currentLap.LapIsStrange = false
					}
				} else {
					currentLap.FasterOrSlowerThanPreviousLapTime = 0
					currentLap.LapIsStrange = false
				}
			}
		}
	}
}

// Сортируем круги по времени обнаружения
func sortLaps(laps []Data.Lap, config Config.Settings) {
	if config.AVERAGE_RESULTS {
		sortLapsAscByDiscoveryAverageUnixTime(laps)
	} else {
		sortLapsAscByDiscoveryMinimalUnixTime(laps)
	}
}

// Рассчитываем позиции и отставания
func calculatePositionsAndDelays(laps []Data.Lap, config Config.Settings) {
	var latestLapsInRace []Data.Lap
	var leaderLap Data.Lap

	for index, lap := range laps {
		if lap.RaceId == laps[0].RaceId {
			if !containsTagId(latestLapsInRace, lap.TagId) {
				laps[index].LapIsLatest = true
				latestLapsInRace = append(latestLapsInRace, lap)
			} else {
				laps[index].LapIsLatest = false
			}
		}
	}

	latestLapsInRace = sortMaxLapNumberAndMinRaceTotalTime(latestLapsInRace)

	for position, lapData := range latestLapsInRace {
		if position == 0 {
			leaderLap = lapData
		}

		for index, lap := range laps {
			if lap.Id == lapData.Id {
				if config.RACE_TYPE == "mass-start" {
					laps[index].RacePosition = uint(position) + 1
					laps[index].LapPosition = uint(position) + 1
				} else if config.RACE_TYPE == "delayed-start" {
					if lap.LapNumber == 0 {
						laps[index].RacePosition = 0
						laps[index].LapPosition = 0
					} else {
						laps[index].RacePosition = uint(position) + 1
						laps[index].LapPosition = uint(position) + 1
					}
				}
				laps[index].TimeBehindTheLeader = lap.RaceTotalTime - leaderLap.RaceTotalTime
			}
		}
	}

	if !config.VARIABLE_DISTANCE_RACE {
		setFastestLaps(laps)
	}
}

// Функция для очистки экрана
func clearScreen() {
	fmt.Print("\033[H\033[2J")
}


// Логируем результаты
func logResults(laps []Data.Lap) {
	fmt.Printf("Номер заезда: %d, Текущий круг: %d\n", laps[0].RaceId, laps[0].LapNumber)
	fmt.Printf("%-10s | %-20s | %-30s | %-8s | %-10s | %-10s | %-10s\n", 
		"RacePos", "Behind", "TagId", "LapNum", "LapTime", "RaceTime", "Latest?")
	fmt.Println(strings.Repeat("-", 120))
	for _, lap := range laps {
		fmt.Printf("%-10d | %-20d | %-30s | %-8d | %-10d | %-10d | %-10t \n", 
			lap.RacePosition, lap.TimeBehindTheLeader, lap.TagId, lap.LapNumber, lap.LapTime, lap.RaceTotalTime, lap.LapIsLatest)
	}
}
