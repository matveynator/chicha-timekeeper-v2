package main

import (
	"time"
	"chicha/pkg/config"
	"chicha/pkg/listener"
	"chicha/pkg/database"
	"chicha/pkg/proxy"
	"chicha/pkg/mylog"
	"chicha/pkg/timekeeper"
)

func main() {

	settings := Config.ParseFlags()

	//run error log daemon
	go MyLog.ErrorLogWorker()
	go Database.Run(settings)
	go Proxy.Run(settings)
	go listener.Run(settings)
	go Timekeeper.Run(settings)

	for {
		time.Sleep(10 * time.Second)	
	}

}
