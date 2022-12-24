package Proxy

import (
	"log"
	"fmt"
	"net"
	"time"

	"chicha/packages/mylog"
	"chicha/packages/data"
	"chicha/packages/config"
)
//оставляем только один процесс который будет брать задачи и передавать их далее на другой сервер (чтобы сохранялась последовательность)
var proxyWorkersMaxCount int = 1
var proxyTask chan Data.AverageResult
var respawnLock chan int

func init() {

	if Config.PROXY_ADDRESS != "" {

		//initialise channel with tasks:
		proxyTask = make(chan Data.AverageResult)

		//initialize unblocking channel to guard respawn tasks
		respawnLock = make(chan int, proxyWorkersMaxCount)

		go func() {
			for {
				// will block if there is proxyWorkersMaxCount ints in respawnLock
				respawnLock <- 1 
				//sleep 1 second
				time.Sleep(1 * time.Second)
				go proxyWorkerRun(len(respawnLock))
			}
		}()

		//for workerId := 1; workerId <= proxyWorkersMaxCount; workerId++ {
		//	go proxyWorkerRun(workerId)
		//}

	}
}

func CreateNewProxyTask(taskData Data.AverageResult) {
	//log.Println("new proxy task:", taskData.TagID)
	proxyTask <- taskData
}

//close connection on programm exit
func deferCleanup(connection net.Conn) {
	<-respawnLock
	if connection != nil {
		err := connection.Close() 
		if err != nil {
			log.Println("Error closing proxy connection:", err)
		}
	}
}

func proxyWorkerRun(workerId int) {

	connection, err := net.DialTimeout("tcp", Config.PROXY_ADDRESS, 15 * time.Second)

	defer deferCleanup(connection)

	if err != nil  {
		MyLog.Printonce(fmt.Sprintf("Proxy destination %s unreachable. Error: %s",  Config.PROXY_ADDRESS, err))
		return

	} else {
		MyLog.Println(fmt.Sprintf("Proxy worker #%d connected to destination %s", workerId, Config.PROXY_ADDRESS))
	}

	//initialise connection error channel
	connectionErrorChannel := make(chan error)

	go func() {
		buffer := make([]byte, 1024)
		for {
			numberOfLines, err := connection.Read(buffer)
			if err != nil {
				connectionErrorChannel <- err
				return
			}
			if numberOfLines > 0 {
				log.Println("Proxy worker received unexpected data back: %s", buffer[:numberOfLines])
			}
		}
	}()

	for {
		select {
			//в случае если есть задание в канале proxyTask
		case currentProxyTask := <- proxyTask :
			//fmt.Println("proxyWorker", workerId, "processing new job...")
			_, networkSendingError := fmt.Fprintf(connection, "%s, %d, %d, %s\n", currentProxyTask.TagID, currentProxyTask.DiscoveryUnixTime, currentProxyTask.Antenna, currentProxyTask.IP)
			if err != nil {
				//в случае потери связи во время отправки мы возвращаем задачу обратно в канал proxyTask
				proxyTask <- currentProxyTask
				log.Printf("Proxy worker %d exited due to sending error: %s\n", workerId, networkSendingError)
				//и завершаем работу гоурутины
				return
			} else {
				//fmt.Println("proxyWorker", workerId, "finished job.")
			}
		case networkError := <-connectionErrorChannel :
			//обнаружена сетевая ошибка - завершаем гоурутину
			log.Printf("Proxy worker %d exited due to connection error: %s\n", workerId, networkError)
			return
		}
	}
}

