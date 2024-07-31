package Timekeeper

import (
	"chicha/pkg/data"
	"fmt"
	"strings"
)

// Логируем результаты
func logResults(laps []Data.Lap) {
	fmt.Printf("Номер заезда: %d, Текущий круг: %d\n", laps[0].RaceId, laps[0].LapNumber)
	fmt.Printf("%-8s | %-6s | %-30s | %-8s | %-10s | %-10s | %-10s\n", 
		"RacePos", "Behind", "TagId", "LapNum", "LapTime", "RaceTime", "Latest?")
	fmt.Println(strings.Repeat("-", 113))
	for _, lap := range laps {
		if lap.LapIsLatest {
		fmt.Printf("%-8d | %-6d | %-30s | %-8d | %-10d | %-10d | %-10t \n", 
			lap.RacePosition, lap.TimeBehindTheLeader, lap.TagId, lap.LapNumber, lap.LapTime, lap.RaceTotalTime, lap.LapIsLatest)
		}
	}
}
