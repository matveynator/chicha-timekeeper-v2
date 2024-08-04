package Timekeeper

import (
	"sort"
	"chicha/pkg/data"
)

// Sort slice by discovery maximal unix time descending (big -> small):
func sortLapsDescByDiscoveryMaximalUnixTime (lapsToSort []data.Lap) {
	sort.Slice(lapsToSort, func(i, j int) bool {
		return lapsToSort[i].DiscoveryMaximalUnixTime > lapsToSort[j].DiscoveryMaximalUnixTime
	})
}
