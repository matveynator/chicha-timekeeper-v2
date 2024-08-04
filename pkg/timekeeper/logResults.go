package Timekeeper

import (
	"chicha/pkg/data"
	"fmt"
	"strings"
)

// Логируем результаты
func logResults(laps []data.Lap) {
	if len(laps) == 0 {
		fmt.Println("Нет данных для отображения.")
	} else {
		// Сортируем последние круги по номеру круга и общему времени гонки
		laps = sortMaxLapNumberAndMinRaceTotalTime(laps)
		fmt.Printf("Номер заезда: %d, Текущий круг: %d, Посчитано результатов: %d\n", laps[0].RaceId, laps[0].LapNumber, len(laps))
		fmt.Printf("%-3s | %-7s | %-30s | %-3s | %-7s | %-10s\n",
			"Pos", "Leader", "TagId", "Lap", "LapTime", "RaceTime")
		fmt.Println(strings.Repeat("-", 75))
		for _, lap := range laps {
			if lap.LapIsLatest {
				fmt.Printf("%-3d | %-7.3f | %-30s | %-3d | %-7.3f | %-10.3f\n", lap.RacePosition, float64(lap.TimeBehindTheLeader)/1000, lap.TagId, lap.LapNumber, float64(lap.LapTime)/1000, float64(lap.RaceTotalTime)/1000)
			}
		}
	}
}
