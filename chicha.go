package main

import (
	"chicha/packages/collector"
	"chicha/packages/mylog"
)


func main() {
	//run error log daemon
	go MyLog.ErrorLogWorker()


	Collector.StartDataCollector()

}
