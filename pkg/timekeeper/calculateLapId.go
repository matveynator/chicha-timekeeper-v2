package Timekeeper

import (
	"chicha/pkg/data"
)

// calculateLapId рассчитывает идентификатор (Id) текущего круга на основе предыдущих кругов (если они есть).
// currentLap - текущий круг, для которого рассчитывается Id.
// laps - список предыдущих кругов.
func calculateLapId(currentLap Data.Lap, laps []Data.Lap) (id int64) {
	// Если нет предыдущих кругов:
	if len(laps) == 0 {
		// Это первый круг! Устанавливаем id в 1.
		id = 1
	} else {
		// Предыдущие данные о кругах доступны (id будет >= 1)
		var lapFound bool = false
		for _, lap := range laps {
			// Если находим круг с тем же TagId, RaceId и LapNumber, устанавливаем тот же id.
			if currentLap.TagId == lap.TagId && currentLap.RaceId == lap.RaceId && currentLap.LapNumber == lap.LapNumber {
				lapFound = true
				return lap.Id
			}
		}

		// Если такой круг не найден в предыдущих данных, рассчитываем следующий id для круга.
		if !lapFound {
			var lapId int64 = 0
			// Ищем самый высокий id среди предыдущих кругов.
			for _, lap := range laps {
				if lap.Id > lapId {
					lapId = lap.Id
				}
			}
			// Возвращаем увеличенный на 1 самый высокий id.
			id = lapId + 1
		}
	}

	return
}

