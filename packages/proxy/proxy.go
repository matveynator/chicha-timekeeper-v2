package Proxy

import (
	"log"
	"fmt"
	"net"
	"time"

	"chicha/packages/data"
	"chicha/packages/config"
)
//оставляем только один процесс который будет брать задачи и передавать их далее на другой сервер (чтобы сохранялась последовательность)
var proxyWorkersMaxCount = 5
var proxyTask chan Data.AverageResult

func init() {

	if Config.PROXY_ADDRESS != "" {

		//initialise channel with tasks:
		proxyTask = make(chan Data.AverageResult)




		guardChannel := make(chan int, proxyWorkersMaxCount)
		go func() {
			for {
				guardChannel <- 1 // will block if there is proxyWorkersMaxCount ints in guardChannel
				go func() {
					go proxyWorkerRun(555)
					log.Println("Started tcp proxy restream to:", Config.PROXY_ADDRESS)
					//sleep beafore next worker respawn
					time.Sleep(5 * time.Second)
					<-guardChannel // removes an int from guardChannel, allowing another to proceed
				}()

			}
		}()

		//for workerId := 1; workerId <= proxyWorkersMaxCount; workerId++ {
		//	go proxyWorkerRun(workerId)
		//}

		log.Println("Started tcp proxy restream to:", Config.PROXY_ADDRESS)
	}
}




func CreateNewProxyTask(taskData Data.AverageResult) {
	//log.Println("new proxy task:", taskData.TagID)
	proxyTask <- taskData
}

func proxyWorkerRun(workerId int) {
	//connection, err := net.Dial("tcp", Config.PROXY_ADDRESS)
	connection, err := net.DialTimeout("tcp", Config.PROXY_ADDRESS, 15 * time.Second)
	if err != nil {
		log.Printf("Proxy worker %d exited due to net.Dial error: %s\n", workerId, err)
		return
	}
	//close connection on programm exit
	defer connection.Close()


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

