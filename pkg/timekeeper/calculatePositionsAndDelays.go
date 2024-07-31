package Timekeeper

import (
	"chicha/pkg/config"
	"chicha/pkg/data"
)

// Рассчитываем позиции и отставания
func calculatePositionsAndDelays(laps []Data.Lap, config Config.Settings) {
    var latestLapsInRace []Data.Lap
    var leaderLap Data.Lap

    // Mark the latest laps
    lapMap := make(map[string]Data.Lap)
    for _, lap := range laps {
        if lap.RaceId == laps[0].RaceId {
            if currentLap, exists := lapMap[lap.TagId]; !exists || lap.LapNumber > currentLap.LapNumber || 
               (lap.LapNumber == currentLap.LapNumber && lap.RaceTotalTime < currentLap.RaceTotalTime) {
                lapMap[lap.TagId] = lap
            }
        }
    }

    for i, lap := range laps {
        if latestLap, exists := lapMap[lap.TagId]; exists && lap.Id == latestLap.Id {
            laps[i].LapIsLatest = true
            latestLapsInRace = append(latestLapsInRace, lap)
        } else {
            laps[i].LapIsLatest = false
        }
    }

    // Sort latest laps by lap number and total race time
    latestLapsInRace = sortMaxLapNumberAndMinRaceTotalTime(latestLapsInRace)

    // Assign positions and calculate delays
    for position, lapData := range latestLapsInRace {
        if position == 0 {
            leaderLap = lapData
        }

        for i, lap := range laps {
            if lap.Id == lapData.Id {
                laps[i].RacePosition = uint(position) + 1
                laps[i].LapPosition = uint(position) + 1

                if config.RACE_TYPE == "delayed-start" && lap.LapNumber == 0 {
                    laps[i].RacePosition = 0
                    laps[i].LapPosition = 0
                }

                laps[i].TimeBehindTheLeader = lap.RaceTotalTime - leaderLap.RaceTotalTime
            }
        }
    }

    if !config.VARIABLE_DISTANCE_RACE {
        setFastestLaps(laps)
    }
}

