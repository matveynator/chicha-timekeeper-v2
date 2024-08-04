package Timekeeper

import (
	"math"
	"chicha/pkg/data"
)

// Функция которая отсортирует слайс и для каждой уникальной группы данных с определенным RaceId выберет все LapNumber больше нуля и выберет среди них минимальное  значение поле LapTime и поставит для этой структуры FastestLapInThisRace true, а всем остальным  из тойже группы одинаковых RaceId поставит FastestLapInThisRace false
func setFastestLaps(laps []data.Lap) {

	type UniqueRaceId struct {
		RaceId 	uint
		LapId  	int64
		LapTime int64
	}

  uniqueRaceIds := make([]UniqueRaceId, 0)

	// Мапа для хранения уникальных RaceId
	raceIdsMap := make(map[uint]bool)

	// Итерируемся по всем структурам Lap и выбираем уникальные RaceId
	for _, lap := range laps {
		if _, ok := raceIdsMap[lap.RaceId]; !ok {
			raceIdsMap[lap.RaceId] = true
			uniqueRaceIds = append(uniqueRaceIds, UniqueRaceId{
				RaceId:  lap.RaceId,
				LapId:   0,
				LapTime: math.MaxInt64,
			})
		}
	}
	// Update uniqueRaceIds slice with LapTime to maximum int64 value available:
	//for uniqueRaceIdIndex, _ := range uniqueRaceIds{
	//	uniqueRaceIds[uniqueRaceIdIndex].LapTime = math.MaxInt64
	//}

	// Update uniqueRaceIds slice with fastest LapId and LapTime:
	for lapIndex, lap := range laps {

		// Unset all laps FastestLapInThisRace value globally:
		laps[lapIndex].FastestLapInThisRace = false

		// Update uniqueRaceIds slice with fastest laps in it:
		for uniqueRaceIdIndex, uniqueRaceId := range uniqueRaceIds{
			if lap.LapNumber > 0 && lap.RaceId == uniqueRaceId.RaceId {
				if lap.LapTime < uniqueRaceId.LapTime {
						uniqueRaceIds[uniqueRaceIdIndex].LapId = lap.Id
						uniqueRaceIds[uniqueRaceIdIndex].LapTime = lap.LapTime
				}
			}
		}
	}

	// Update fastest laps in each race and set FastestLapInThisRace = true:
	for lapIndex, lap := range laps {
		if lap.LapNumber > 0 {
			for _, uniqueRaceId := range uniqueRaceIds{
				if lap.Id == uniqueRaceId.LapId {
					laps[lapIndex].FastestLapInThisRace = true
				}
			}
		}
	}
}

