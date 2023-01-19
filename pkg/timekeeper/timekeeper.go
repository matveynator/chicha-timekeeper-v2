package Timekeeper

import (
	"log"
	"runtime"

	"chicha/pkg/data"
)

// To send a task: 
//
// import (
// "chicha/pkg/timekeeper"
// "chicha/pkg/data"
// )
//
// Timekeeper.TimekeeperTask <- Data.RawData

var timekeeperWorkersMaxCount int
var TimekeeperTask chan Data.RawData
var respawnLock chan int

func init() {

	//set max workers equal to number of CPU cores:
	timekeeperWorkersMaxCount = runtime.NumCPU()

	//initialise channel with tasks:
	TimekeeperTask = make(chan Data.RawData)

	//initialize unblocking channel to guard respawn tasks
	respawnLock = make(chan int, timekeeperWorkersMaxCount)

	go func() {
		for {
			// will block if there is timekeeperWorkersMaxCount ints in respawnLock
			respawnLock <- 1 
			go timekeeperWorkerRun(len(respawnLock))
		}
	}()
}

func timekeeperWorkerRun(workerId int) {
	for {
		select {
			//в случае если есть задание в канале TimekeeperTask
		case currentTimekeeperTask := <- TimekeeperTask :
			log.Printf("timekeeperWorker %d received a new job %s\n", workerId, currentTimekeeperTask.TagID)
		}
	}
}

