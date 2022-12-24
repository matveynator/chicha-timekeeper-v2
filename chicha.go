package main

import (
	"chicha/packages/database"
	"chicha/packages/collector"
	"chicha/packages/mylog"
)


func main() {
  //run error log daemon
	go MyLog.ErrorLogWorker()

  Database.CreateDB()
	Database.UpdateDB()
  Collector.StartDataCollector()


}

