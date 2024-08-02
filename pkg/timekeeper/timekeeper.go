package Timekeeper

import (
	"log"
	"fmt"
	"chicha/pkg/data"
	"chicha/pkg/config"
	"chicha/pkg/database"

)



//оставляем только один процесс который будет калькулировать время
var timekeeperWorkersMaxCount int = 1
var TimekeeperRawData chan Data.RawData
var respawnLock chan int

func Run(config Config.Settings) {

	//initialize channel with tasks:
	TimekeeperRawData = make(chan Data.RawData)

	//initialize unblocking channel to guard respawn tasks
	respawnLock = make(chan int, timekeeperWorkersMaxCount)

	go func() {
		for {
			// will block if there is timekeeperWorkersMaxCount ints in respawnLock
			respawnLock <- 1
			//sleep 1 second
			//time.Sleep(1 * time.Second)
			go timekeeperWorkerRun(config)
		}
	}()
}



func timekeeperWorkerRun(config Config.Settings) (err error) {
	var CurrentLaps []Data.Lap

	for {
		select {
			//в случае если есть задание в канале TimekeeperRawData
		case currentTimekeeperRawData := <- TimekeeperRawData :

			// If this is the first race data - check database if any previous race data availabe?
			if len(CurrentLaps) == 0 {

				// Request to get recent race laps if awailable in database (blocking!):
				Database.RequestRecentRaceLapsChan <- currentTimekeeperRawData
				// Wait for an answer (blocking!)
				CurrentLaps = <- Database.ReplyWithRecentRaceLapsChan

			} 

			CurrentLaps, err = calculateRaceInMemory(currentTimekeeperRawData, CurrentLaps, config)
			if err != nil {
				log.Println("Timekeeper fatal error: ", err)
				return
			}	else {
				// показываем результаты в терминале
				// clear screen:
				fmt.Print("\033[H\033[2J")
				// выводим в терминал форматированно:
				logResults(CurrentLaps)
			}
		}
	}
}


