package Timekeeper

import (
	"sort"
	"chicha/pkg/data"
)

// Сортировка по максимальному LapNumber и минимальному RaceTotalTime:
func sortMaxLapNumberAndMinRaceTotalTime(laps []data.Lap) []data.Lap {
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

