package Timekeeper

import (
	"sort"
	"chicha/pkg/data"
)

// Sort slice by discovery average unix time ascending (small -> big):
func sortLapsAscByDiscoveryAverageUnixTime (lapsToSort []data.Lap) {
	sort.Slice(lapsToSort, func(i, j int) bool {
		return lapsToSort[i].DiscoveryAverageUnixTime < lapsToSort[j].DiscoveryAverageUnixTime
	})
}
