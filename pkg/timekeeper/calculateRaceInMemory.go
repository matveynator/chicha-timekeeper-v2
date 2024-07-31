package Timekeeper

import (
	"chicha/pkg/config"
	"chicha/pkg/data"
	"log"
	"fmt"
)

// Функция обрабатывает данные гонки в памяти и возвращает обновленные круги и возможную ошибку.
// currentTimekeeperTask - текущие данные круга, полученные от RFID.
// previousLaps - данные о предыдущих кругах.
// config - настройки гонки.
func calculateRaceInMemory(currentTimekeeperTask Data.RawData, previousLaps []Data.Lap, config Config.Settings) (currentLaps []Data.Lap, err error) {

	var currentLap Data.Lap
	var lapIsValid bool

	// Проверка валидности текущих показаний гонщика, чтобы убедиться, что он не устарел и соответствует временным ограничениям заданным в конфиге гонки чтобы решить добавлять их в гонку или нет (защита от подтасовки результатов - например невозможно пройти круг быстрее чем определенное время указанное в конфиге или нельзя накидать результатов улучшающих результат спустя определенное время или с неавторизованного айпи адреса / считывателя).
	lapIsValid, err = checkLapIsValid(currentTimekeeperTask, previousLaps, config)
	if err != nil {
		log.Printf("Ошибка валидации данных: %s \n", err)
		return
	} else {
		if !lapIsValid {
			// Присланные данные мы считаем что невалидны, пропускаем их и возвращаем предыдущие круги.
			log.Println("Получены невалидные данные круга. Пропускаем.")
			currentLaps = previousLaps
			return
		}
	}

	// Проверка, есть ли уже данные о кругах в памяти.
	if len(previousLaps) == 0 {
		// Нет данных о гонке и кругах в памяти - создаем первую гонку и круг.
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
		// Данные о гонке уже есть в памяти, обновляем текущий круг новыми данными.
		currentLap.TagId = currentTimekeeperTask.TagId
		currentLap.DiscoveryMinimalUnixTime = currentTimekeeperTask.DiscoveryUnixTime
		currentLap.DiscoveryMaximalUnixTime = currentTimekeeperTask.DiscoveryUnixTime
		currentLap.DiscoveryAverageUnixTime = currentTimekeeperTask.DiscoveryUnixTime

		// Рассчитываем номер заезда в котором едет гонщик.
		previousLapsCopy := previousLaps
		currentLap.RaceId, err = getCurrentRaceId(currentTimekeeperTask, previousLapsCopy, config)
		if err != nil {
			log.Println(err)
			return
		}

		// Рассчитываем номер круга который едет гонщик в данном заезде.
		currentLap.LapNumber, err = getMyCurrentRaceLapNumber(currentTimekeeperTask, previousLapsCopy, config)
		if err != nil {
			log.Println(err)
			return
		}
	}

	// Рассчитываем ID текущего круга. calculateLapId рассчитывает идентификатор (Id) текущего круга на основе предыдущих кругов (если они есть)
	currentLap.Id = calculateLapId(currentLap, previousLaps)

	// Объявляем срез для хранения данных о предыдущих кругах, которые не связаны с текущим кругом.
	var otherOldLaps []Data.Lap

	// Проходим по каждому кругу из списка предыдущих кругов.
	for _, previousLap := range previousLaps {
		// Проверяем, совпадают ли идентификаторы текущего и предыдущего кругов:
		// - TagId (идентификатор метки RFID)
		// - RaceId (идентификатор гонки)
		// - LapNumber (номер круга)
		if currentLap.TagId == previousLap.TagId && currentLap.RaceId == previousLap.RaceId && currentLap.LapNumber == previousLap.LapNumber {
			// Если идентификаторы совпадают, обновляем данные текущего круга на основе предыдущего круга.

			// Обновляем минимальное время обнаружения, если у предыдущего круга оно меньше.
			if currentLap.DiscoveryMinimalUnixTime > previousLap.DiscoveryMinimalUnixTime {
				currentLap.DiscoveryMinimalUnixTime = previousLap.DiscoveryMinimalUnixTime
			}

			// Обновляем максимальное время обнаружения, если у предыдущего круга оно больше.
			if currentLap.DiscoveryMaximalUnixTime < previousLap.DiscoveryMaximalUnixTime {
				currentLap.DiscoveryMaximalUnixTime = previousLap.DiscoveryMaximalUnixTime
			}

			// Пересчитываем среднее время обнаружения как среднее арифметическое между текущим и предыдущим значениями.
			currentLap.DiscoveryAverageUnixTime = (currentTimekeeperTask.DiscoveryUnixTime + previousLap.DiscoveryAverageUnixTime) / 2

			// Если у предыдущего круга количество результатов полученных от транспондера не задано, устанавливаем его в 1.
			if previousLap.AverageResultsCount < 1 {
				previousLap.AverageResultsCount = 1
			}

			// Увеличиваем общее среднее количество результатов текущего круга полученных от транспондера на единицу.
			currentLap.AverageResultsCount = previousLap.AverageResultsCount + 1

		} else {
			// Если идентификаторы текущего и предыдущего кругов НЕ совпадают,
			// добавляем предыдущий круг в срез других кругов.
			otherOldLaps = append(otherOldLaps, previousLap)
		}
	}

	// Рассчитываем время круга.
	currentLap.LapTime, err = calculateLapTime(currentLap, otherOldLaps, config)
	if err != nil {
		log.Printf("Ошибка в расчете времени круга: %s \n", err)
		return
	}

	// Рассчитываем лучшее время круга и номер лучшего круга, если гонка с фиксированной дистанцией.
	if !config.VARIABLE_DISTANCE_RACE {
		currentLap.BestLapTime, currentLap.BestLapNumber, err = calculateMyBestLapTime(currentLap, otherOldLaps)
		if err != nil {
			log.Printf("Ошибка в расчете лучшего времени круга: %s \n", err)
			return
		}
	}

	// Рассчитываем общее время гонки для текущего круга.
	if currentLap.LapNumber == 0 {
		// Если это первый круг (LapNumber == 0), устанавливаем общее время гонки равным времени текущего круга.
		currentLap.RaceTotalTime = currentLap.LapTime

		// Поскольку это первый круг, то нет предыдущего круга для сравнения, поэтому время на круг быстрее или медленнее предыдущего круга устанавливаем в 0.
		currentLap.FasterOrSlowerThanPreviousLapTime = 0

		// Первый круг не может быть странным, поэтому устанавливаем флаг LapIsStrange в false.
		currentLap.LapIsStrange = false
	} else {
		// Если это не первый круг, проходим по всем другим кругам в памяти (otherOldLaps).
		for _, otherOldLap := range otherOldLaps {
			// Находим предыдущий круг того же участника (TagId) в той же гонке (RaceId) и предыдущим номером круга (LapNumber - 1).
			if currentLap.TagId == otherOldLap.TagId && currentLap.RaceId == otherOldLap.RaceId && otherOldLap.LapNumber == currentLap.LapNumber-1 {
				// Общее время гонки равно времени гонки до предыдущего круга плюс время текущего круга.
				currentLap.RaceTotalTime = otherOldLap.RaceTotalTime + currentLap.LapTime

				// Если номер текущего круга больше 1 и гонка не имеет переменной дистанции.
				if currentLap.LapNumber > 1 && !config.VARIABLE_DISTANCE_RACE {
					// Рассчитываем разницу между временем текущего круга и временем предыдущего круга.
					currentLap.FasterOrSlowerThanPreviousLapTime = currentLap.LapTime - otherOldLap.LapTime

					// Если время текущего круга больше или равно времени предыдущего круга более чем на 100%, считаем текущий круг подозрительным.
					if (currentLap.LapTime - otherOldLap.LapTime) >= otherOldLap.LapTime {
						currentLap.LapIsStrange = true
					} else {
						currentLap.LapIsStrange = false
					}
				} else {
					// Для первого и второго кругов в гонке разница во времени и подозрительность круга устанавливаются в 0 и false соответственно.
					currentLap.FasterOrSlowerThanPreviousLapTime = 0
					currentLap.LapIsStrange = false
				}
			}
		}
	}

	// Объединяем старые и обновленные данные кругов.
	currentLaps = append(otherOldLaps, currentLap)

	// Сортируем круги по времени обнаружения в зависимости от глобальной настройки AVERAGE_RESULTS.
	if config.AVERAGE_RESULTS {
		sortLapsAscByDiscoveryAverageUnixTime(currentLaps)
	} else {
		sortLapsAscByDiscoveryMinimalUnixTime(currentLaps)
	}

	// Рассчитываем позицию гонки и круга.
	var currentRaceLatestLaps []Data.Lap

	// Проходим по всем текущим кругам.
	for index, currentInMemoryLap := range currentLaps {
		// Проверяем, совпадает ли идентификатор гонки текущего круга с идентификатором гонки обрабатываемого круга.
		if currentInMemoryLap.RaceId == currentLap.RaceId {
			// Проверяем, содержится ли TagId текущего круга в списке последних кругов текущей гонки.
			if !containsTagId(currentRaceLatestLaps, currentInMemoryLap.TagId) {
				// Если нет, помечаем этот круг как последний.
				currentLaps[index].LapIsLatest = true
				// Добавляем этот круг в список последних кругов текущей гонки.
				currentRaceLatestLaps = append(currentRaceLatestLaps, currentInMemoryLap)
			} else {
				// Если да, помечаем этот круг как не последний.
				currentLaps[index].LapIsLatest = false
			}
		}
	}

	// Сортируем список последних кругов текущей гонки по максимальному номеру круга и минимальному общему времени гонки.
	currentRaceLatestLaps = sortMaxLapNumberAndMinRaceTotalTime(currentRaceLatestLaps)

	// Объявляем переменную для хранения данных о лидирующем круге текущей гонки.
	var leaderCurrentLap Data.Lap

	// Проходим по отсортированным последним кругам текущей гонки.
	for latestPosition, latestLapData := range currentRaceLatestLaps {
		// Если это первый круг (лидирующий), сохраняем его данные.
		if latestPosition == 0 {
			leaderCurrentLap = latestLapData
		}

		// Обновляем данные о позиции круга и времени отставания от лидера для всех текущих кругов.
		for index, lap := range currentLaps {
			// Проверяем, совпадает ли идентификатор текущего круга с идентификатором последнего круга.
			if lap.Id == latestLapData.Id {
				// Обновляем данные в зависимости от типа гонки.
				if config.RACE_TYPE == "mass-start" {
					// Для гонки с массовым стартом устанавливаем позицию в гонке и кругу.
					currentLaps[index].RacePosition = uint(latestPosition) + 1
					currentLaps[index].LapPosition = uint(latestPosition) + 1
				} else if config.RACE_TYPE == "delayed-start" {
					// Для гонки с раздельным стартом:
					if lap.LapNumber == 0 {
						// Если это нулевой круг, позиция равна нулю.
						currentLaps[index].RacePosition = 0
						currentLaps[index].LapPosition = 0
					} else {
						// Для всех остальных кругов устанавливаем позицию в гонке и кругу.
						currentLaps[index].RacePosition = uint(latestPosition) + 1
						currentLaps[index].LapPosition = uint(latestPosition) + 1
					}
				}
				// Рассчитываем время отставания от лидера.
				currentLaps[index].TimeBehindTheLeader = lap.RaceTotalTime - leaderCurrentLap.RaceTotalTime
			}
		}
	}

	// Рассчитываем самый быстрый круг в гонке.
	if !config.VARIABLE_DISTANCE_RACE {
		setFastestLaps(currentLaps)
	}

	// Логируем результаты перед возвратом.
	i := 0
	for _, lap := range currentLaps {
		if i == 0 {
			fmt.Printf("Номер заезда: %d, Текущий круг: %d\n", lap.RaceId, lap.LapNumber)
			fmt.Printf("Позиция | Имя гонщика                | Отставание от лидера |\n")
		}
		  i++
			fmt.Printf("  %d    | %s | %d                   |\n", lap.RacePosition, lap.TagId, lap.TimeBehindTheLeader)
	}

	// Возвращаем обновленные данные кругов.
	return
}
