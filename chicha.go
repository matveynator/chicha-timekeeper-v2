package main

import (
	"time"
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
	go Collector.Run(settings)

	for {
		time.Sleep(1 * time.Second)	
	}

}
