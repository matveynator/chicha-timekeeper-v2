package Timekeeper

import (
	"chicha/pkg/config"
	"chicha/pkg/data"
)

// calculatePositionsAndDelays рассчитывает позиции и отставания участников гонки
func calculatePositionsAndDelays(laps []data.Lap, config Config.Settings) {
	var latestLapSlice []data.Lap // Массив для хранения последних кругов в гонке
	var leaderLap data.Lap        // Переменная для хранения круга лидера гонки

	// Создаем карту для хранения последних кругов каждого участника
	latestLapMap := make(map[string]data.Lap)
	for _, lap := range laps {
		// Проверяем, принадлежит ли круг текущей гонке
		if lap.RaceId == laps[0].RaceId {
			currentLap, exists := latestLapMap[lap.TagId]

			if !exists {
				// Если круг участника еще не записан в карту, добавляем его
				latestLapMap[lap.TagId] = lap
			} else if lap.LapNumber > currentLap.LapNumber {
				// Если номер текущего круга больше, обновляем запись
				latestLapMap[lap.TagId] = lap
			} else if lap.LapNumber == currentLap.LapNumber {
				if lap.RaceTotalTime < currentLap.RaceTotalTime {
					// Если номера кругов равны, но время текущего круга меньше, обновляем запись
					latestLapMap[lap.TagId] = lap
				}
			}
		}
	}

	// Маркируем последние круги в общем списке кругов
	for i, lap := range laps {
		latestLapInMap, exists := latestLapMap[lap.TagId]

		if exists && lap.Id == latestLapInMap.Id {
			laps[i].LapIsLatest = true
			latestLapSlice = append(latestLapSlice, lap)
		} else {
			laps[i].LapIsLatest = false
		}
	}

	// Сортируем последние круги по номеру круга и общему времени гонки
	latestLapSlice = sortMaxLapNumberAndMinRaceTotalTime(latestLapSlice)

	// Присваиваем позиции и рассчитываем отставания
	for position, latestLapInSlice := range latestLapSlice {
		// Первый круг в списке - это круг лидера
		if position == 0 {
			leaderLap = latestLapInSlice
		}

		for i, lap := range laps {
			// Проходим только по кругам из слайса свежих кругов гонки:
			if lap.Id == latestLapInSlice.Id {

				// Присваиваем позицию в гонке и позицию круга
				// Для гонок с отложенным стартом, если номер круга равен 0, позиция устанавливается в 0
				if config.RACE_TYPE == "delayed-start" && lap.LapNumber == 0 {
					laps[i].RacePosition = 0
					laps[i].LapPosition = 0
				} else {
					laps[i].RacePosition = uint(position) + 1
					laps[i].LapPosition = uint(position) + 1
				}

				// Рассчитываем время отставания от лидера только если номера кругов совпадают
				if lap.LapNumber == leaderLap.LapNumber {
					laps[i].TimeBehindTheLeader = leaderLap.RaceTotalTime - lap.RaceTotalTime
				} else {
					/*
														// Отставание от лидера в случае отставания на 1 или более кругов равно отставанию от лидера в последнем круге пройденном этим гонщиком плюс время прощедшее с момента прохождения гонщиком этого круга до настоящего момента

														// Пробегаем по старым кругам и находим время лидера
														for _, lapFromMemory := range laps {

														// Ищем данные и время лидера старого круга в котором ехал гонщик
														if lapFromMemory.LapNumber == lap.LapNumber {
														if lapFromMemory.RacePosition == 1 {
					oldLapLeader := lapFromMemory

														//Время лидера того круга в котором ехал гонщик:
														if config.AVERAGE_RESULTS {
														oldLapLeader.DiscoveryAverageUnixTime
														} else {
														oldLapLeader.DiscoveryMinimalUnixTime
														}
														}
														}

														}
														//laps[i].TimeBehindTheLeader = 0
					*/

				}
			}
		}
	}

	// Если гонка не является переменной дистанцией, устанавливаем самые быстрые круги
	if !config.VARIABLE_DISTANCE_RACE {
		setFastestLaps(laps)
	}
}
