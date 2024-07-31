package Timekeeper

import (
	"chicha/pkg/data"
	"chicha/pkg/config"
)

// Функция calculateLapTime вычисляет время круга для текущего круга на основе данных предыдущих кругов и настроек конфигурации.
func calculateLapTime(currentLap Data.Lap, otherOldLaps []Data.Lap, config Config.Settings) (lapTime int64, err error) {

	// Создаем копию массива предыдущих кругов.
	allLaps := make([]Data.Lap, len(otherOldLaps))
	copy(allLaps, otherOldLaps)

	// Добавляем текущий круг к копии массива всех кругов.
	allLaps = append(allLaps, currentLap)

	// Сортируем круги по времени в порядке возрастания.
	if config.AVERAGE_RESULTS {
		sortLapsAscByDiscoveryAverageUnixTime(allLaps)
	} else {
		sortLapsAscByDiscoveryMinimalUnixTime(allLaps)
	}

	// Находим все круги с тем же номером круга и тем же идентификатором гонки, что и текущий круг.
	var sameLapNumberLaps []Data.Lap
	for _, lap := range allLaps {
		if lap.RaceId == currentLap.RaceId && lap.LapNumber == currentLap.LapNumber {
			sameLapNumberLaps = append(sameLapNumberLaps, lap)
		}
	}

	// Сортируем эти круги по времени в порядке возрастания.
	if config.AVERAGE_RESULTS {
		sortLapsAscByDiscoveryAverageUnixTime(sameLapNumberLaps)
	} else {
		sortLapsAscByDiscoveryMinimalUnixTime(sameLapNumberLaps)
	}

	// Обработка времени круга в зависимости от типа старта гонки.
	if config.RACE_TYPE == "delayed-start" {
		// Для гонок с отсроченным стартом первый проход через ворота является стартовым временем,
		// поэтому первый (нулевой) круг имеет нулевое время.
		if	currentLap.LapNumber == 0 {
			lapTime = 0 
		} else {
			// Для не нулевых кругов вычисляем разницу между временем текущего и предыдущего круга.
			for _, allLapData := range allLaps {
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
		// Для гонок с массовым стартом первый проход через ворота считается первым малым кругом,
		// и первый (нулевой) круг является разницей во времени между первым прохождением ворот (лидером) и текущим временем.
		var leaderTime int64 
		if config.AVERAGE_RESULTS {
			leaderTime = sameLapNumberLaps[0].DiscoveryAverageUnixTime
			lapTime = currentLap.DiscoveryAverageUnixTime - leaderTime
		} else {
			leaderTime = sameLapNumberLaps[0].DiscoveryMinimalUnixTime
			lapTime = currentLap.DiscoveryMinimalUnixTime - leaderTime
		}

		for _, allLapData := range allLaps { 
			// Для кругов больше нуля вычисляем время круга как разницу между текущим и предыдущим временем круга.
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

