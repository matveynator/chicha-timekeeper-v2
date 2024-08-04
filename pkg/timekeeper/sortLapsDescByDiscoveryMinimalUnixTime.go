package Timekeeper

import (
	"sort"
	"chicha/pkg/data"
)

// Sort slice by discovery minimal unix time descending (big -> small):
func sortLapsDescByDiscoveryMinimalUnixTime (lapsToSort []data.Lap) {
	sort.Slice(lapsToSort, func(i, j int) bool {
		return lapsToSort[i].DiscoveryMinimalUnixTime > lapsToSort[j].DiscoveryMinimalUnixTime
	})
}
