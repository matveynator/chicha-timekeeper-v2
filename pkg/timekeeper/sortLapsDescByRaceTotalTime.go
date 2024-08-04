package Timekeeper

import (
	"sort"
	"chicha/pkg/data"
)

// Sort slice by race total time descending (big -> small):
func sortLapsDescByRaceTotalTime (lapsToSort []data.Lap) {
	sort.Slice(lapsToSort, func(i, j int) bool {
		return lapsToSort[i].RaceTotalTime > lapsToSort[j].RaceTotalTime
	})
}
