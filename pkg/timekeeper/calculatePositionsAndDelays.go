package Timekeeper

import (
	"chicha/pkg/config"
	"chicha/pkg/data"
)

// calculatePositionsAndDelays рассчитывает позиции и отставания участников гонки
func calculatePositionsAndDelays(laps []Data.Lap, config Config.Settings) {
    var latestLapsInRace []Data.Lap // Массив для хранения последних кругов в гонке
    var leaderLap Data.Lap // Переменная для хранения круга лидера гонки

    // Создаем карту для хранения последних кругов каждого участника
    lapMap := make(map[string]Data.Lap)
    for _, lap := range laps {
        // Проверяем, принадлежит ли круг текущей гонке
        if lap.RaceId == laps[0].RaceId {
            currentLap, exists := lapMap[lap.TagId]

            if !exists {
                // Если круг участника еще не записан в карту, добавляем его
                lapMap[lap.TagId] = lap
            } else if lap.LapNumber > currentLap.LapNumber {
                // Если номер текущего круга больше, обновляем запись
                lapMap[lap.TagId] = lap
            } else if lap.LapNumber == currentLap.LapNumber {
                if lap.RaceTotalTime < currentLap.RaceTotalTime {
                    // Если номера кругов равны, но время текущего круга меньше, обновляем запись
                    lapMap[lap.TagId] = lap
                }
            }
        }
    }

    // Маркируем последние круги в общем списке кругов
    for i, lap := range laps {
        latestLap, exists := lapMap[lap.TagId]
        if exists && lap.Id == latestLap.Id {
            laps[i].LapIsLatest = true
            latestLapsInRace = append(latestLapsInRace, lap)
        } else {
            laps[i].LapIsLatest = false
        }
    }

    // Сортируем последние круги по номеру круга и общему времени гонки
    latestLapsInRace = sortMaxLapNumberAndMinRaceTotalTime(latestLapsInRace)

    // Присваиваем позиции и рассчитываем отставания
    for position, lapData := range latestLapsInRace {
        // Первый круг в списке - это круг лидера
        if position == 0 {
            leaderLap = lapData
        }

        for i, lap := range laps {
            if lap.Id == lapData.Id {
                laps[i].RacePosition = uint(position) + 1
                laps[i].LapPosition = uint(position) + 1

                // Для гонок с отложенным стартом, если номер круга равен 0, позиция устанавливается в 0
                if config.RACE_TYPE == "delayed-start" && lap.LapNumber == 0 {
                    laps[i].RacePosition = 0
                    laps[i].LapPosition = 0
                }

                // Рассчитываем время отставания от лидера только если номера кругов совпадают
                if lap.LapNumber == leaderLap.LapNumber {
                    laps[i].TimeBehindTheLeader = leaderLap.RaceTotalTime - lap.RaceTotalTime
                } else {
                    laps[i].TimeBehindTheLeader = 0 // или какое-то другое значение, чтобы указать, что отставание не рассчитывается
                }
            }
        }
    }

    // Если гонка не является переменной дистанцией, устанавливаем самые быстрые круги
    if !config.VARIABLE_DISTANCE_RACE {
        setFastestLaps(laps)
    }
}

