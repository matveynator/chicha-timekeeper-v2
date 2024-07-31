package Timekeeper

import (
	"sort"
	"chicha/pkg/data"
)

// Sort slice by discovery minimal unix time ascending (small -> big):
func sortLapsAscByDiscoveryMinimalUnixTime (lapsToSort []Data.Lap) {
	sort.Slice(lapsToSort, func(i, j int) bool {
		return lapsToSort[i].DiscoveryMinimalUnixTime < lapsToSort[j].DiscoveryMinimalUnixTime
	})
}
