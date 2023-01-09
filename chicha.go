package main

import (
	"chicha/packages/config"
	"chicha/packages/collector"
	"chicha/packages/database"
	"chicha/packages/proxy"
	"chicha/packages/mylog"
)

func main() {

	settings := Config.ParseFlags()

	//run error log daemon
	go MyLog.ErrorLogWorker()
	go Database.Run(settings)
	go Proxy.Run(settings)
	Collector.StartDataCollector(settings)
}
