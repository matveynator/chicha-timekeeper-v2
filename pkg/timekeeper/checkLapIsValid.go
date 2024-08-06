package Timekeeper

import (
	"log"
	"sort"
	"time"
	"chicha/pkg/data"
	"chicha/pkg/config"
)

// Check time in data is valid:
func checkLapIsValid(lapToCheck data.RawData, previousLaps []data.Lap, config Config.Settings) (valid bool, err error) {

	// Get location from config time zone:
	var location *time.Location
	location, err = time.LoadLocation(config.TIME_ZONE)
	if err != nil {
		// Stop and return error:
		valid = false
		return
	}

	// Calculate my current time:
	myCurrentTime := time.UnixMilli(lapToCheck.DiscoveryUnixTime).In(location)

	// Calculate current system time:
	currentSystemTime := time.Now().In(location)

	// Check if current system time is after lap time + config.MINIMAL_LAP_TIME_DURATION (some very old data received?):
	if currentSystemTime.After(myCurrentTime.Add(config.MINIMAL_LAP_TIME_DURATION)) {
		// Stop and return error:
		valid = false
		return
	}

	// Check: Do we have any previous race data?
	if len(previousLaps) == 0 {
		// This is the first race and lap 
		valid = true
		return
	} else {
		// Sort previous laps by DiscoveryUnixTime:
		previousLapsCopy := previousLaps // Create a copy of the previous laps to avoid altering the original slice.
		sort.Slice(previousLapsCopy, func(i, j int) bool {
			//sort descending by DiscoveryAverageUnixTime or DiscoveryMinimalUnixTime (depending on config.AVERAGE_RESULTS setting:
			if config.AVERAGE_RESULTS {
				return previousLapsCopy[i].DiscoveryAverageUnixTime > previousLapsCopy[j].DiscoveryAverageUnixTime
			} else {
				return previousLapsCopy[i].DiscoveryMinimalUnixTime > previousLapsCopy[j].DiscoveryMinimalUnixTime
			}
		})
		// Get anyones very last time added to memory: 
		var lastLapTime time.Time
		if config.AVERAGE_RESULTS {
			lastLapTime = time.UnixMilli(previousLapsCopy[0].DiscoveryAverageUnixTime).In(location)
		} else {
			lastLapTime = time.UnixMilli(previousLapsCopy[0].DiscoveryMinimalUnixTime).In(location)
		}
		// Check if my current time is after config.RACE_TIMEOUT_DURATION:
		if myCurrentTime.After(lastLapTime.Add(config.RACE_TIMEOUT_DURATION)) {
			// If previous race time expired (timed out), data is valid:
			log.Println("Time is valid as race expired and we will be starting new race.")
			valid = true
			return
		} else {
			// If race not expired, calculate lap time (valid if > minimal lap time before previous results) and not more than precision time.
			// Create slice with only my data in it:
			var myPreviousLaps []data.Lap
			for _, myPreviousLap := range previousLapsCopy {
				if lapToCheck.TagId == myPreviousLap.TagId {
					// Add my previous data to slice:	
					myPreviousLaps = append(myPreviousLaps, myPreviousLap)
				}
			}
			// Check if I have some previous laps in this race:
			if len(myPreviousLaps) > 0 {
				// I have some previous laps in this race.

				// Sort my laps by discovery unix time descending (big -> small) depending on config.AVERAGE_RESULTS global setting:
				if config.AVERAGE_RESULTS {
					sortLapsDescByDiscoveryAverageUnixTime(myPreviousLaps)
				} else {
					sortLapsDescByDiscoveryMinimalUnixTime(myPreviousLaps)
				}
				// Get my last lap time:
				var myLastLapTime time.Time
				if config.AVERAGE_RESULTS {
					myLastLapTime = time.UnixMilli(myPreviousLaps[0].DiscoveryAverageUnixTime).In(location)
				} else {
					myLastLapTime = time.UnixMilli(myPreviousLaps[0].DiscoveryMinimalUnixTime).In(location)
				}
				// Data is invalid if my time is after config.AVERAGE_DURATION and before than config.MINIMAL_LAP_TIME_DURATION:
				if myCurrentTime.After(myLastLapTime.Add(config.AVERAGE_DURATION)) && myCurrentTime.Before(myLastLapTime.Add(config.MINIMAL_LAP_TIME_DURATION)) {
					// Time is more than config.AVERAGE_DURATION and less than config.MINIMAL_LAP_TIME_DURATION - is invalid:
					valid = false
					return
				} else {
					valid = true
					return
				}	
			} else {
				// Race is valid and I have no laps in it yet - set lap is valid:
				valid = true
				return
			}
		}
	}
	return
}
