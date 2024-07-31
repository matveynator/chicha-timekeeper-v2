package Timekeeper

import (
	"sort"
	"chicha/pkg/data"
)

// Sort slice by lap time ascending (small -> big):
func sortLapsAscLapTime (lapsToSort []Data.Lap) {
	sort.Slice(lapsToSort, func(i, j int) bool {
		return lapsToSort[i].LapTime < lapsToSort[j].LapTime
	})
}
