package Timekeeper

import (
	"sort"
	"chicha/pkg/data"
)

// Sort slice by race total time ascending (small -> big):
func sortLapsAscByRaceTotalTime (lapsToSort []data.Lap) {
	sort.Slice(lapsToSort, func(i, j int) bool {
		return lapsToSort[i].RaceTotalTime < lapsToSort[j].RaceTotalTime
	})
}
