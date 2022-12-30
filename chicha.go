package main

import (
	"chicha/packages/database"
	"chicha/packages/collector"
	"chicha/packages/mylog"
)


func main() {
	//run error log daemon
	go MyLog.ErrorLogWorker()

	//spin forever go routine to save in db with some interval:
	go Database.SaveLapsBufferSimplyToDB()

	Database.CreateDB()
	Database.UpdateDB()
	Collector.StartDataCollector()
}
