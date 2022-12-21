package Proxy

import (
	"log"
	"fmt"
	"net"

	"chicha/packages/data"
	"chicha/packages/config"
)

var proxyWorkersCount = 5
var proxyTask chan Data.AverageResult

func init() {

	if Config.PROXY_ADDRESS != "" {

		//initialise channel with tasks:
		proxyTask = make(chan Data.AverageResult)

		for workerId := 1; workerId <= proxyWorkersCount; workerId++ {
			go proxyWorkerRun(workerId)
		}

		log.Println("Started tcp proxy restream to:", Config.PROXY_ADDRESS)
	}
}

func CreateNewProxyTask(taskData Data.AverageResult) {
	log.Println("new proxy task:", taskData.TagID)
	proxyTask <- taskData
}

func proxyWorkerRun(workerId int) {
	connection, err := net.Dial("tcp", Config.PROXY_ADDRESS)
	if err != nil {
		log.Printf("worker %d exited due to net.Dial error: %s", workerId, err)
		return
	}
	//close connection on programm exit
	defer connection.Close()

	for {
		select {
		case currentProxyTask:= <- proxyTask:
			fmt.Println("proxyWorker", workerId, "processing new job...")
			fmt.Fprintf(connection, "%s, %d, %d, %s\n", currentProxyTask.TagID, currentProxyTask.DiscoveryUnixTime, currentProxyTask.Antenna, currentProxyTask.IP)
			fmt.Println("proxyWorker", workerId, "finished job.")
		}
	}
}

