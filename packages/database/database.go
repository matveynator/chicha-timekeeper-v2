package Database

import (
	"log"
	"fmt"
	"time"

	"database/sql"

	"chicha/packages/mylog"
	"chicha/packages/data"
	"chicha/packages/config"
)
//оставляем только один процесс который будет брать задачи и записывать их в базу данных

var databaseWorkersMaxCount int = 1
var DatabaseTask chan Data.RawData
var respawnLock chan int

func init() {

	if Config.DB_TYPE != "" {

		//initialise channel with tasks:
		DatabaseTask = make(chan Data.RawData)

		//initialize unblocking channel to guard respawn tasks
		respawnLock = make(chan int, databaseWorkersMaxCount)

		go func() {
			for {
				// will block if there is databaseWorkersMaxCount in respawnLock
				respawnLock <- 1 
				//sleep 1 second
				time.Sleep(1 * time.Second)
				go databaseWorkerRun(len(respawnLock))
			}
		}()
	}
}

func CreateNewDatabaseTask(taskData Data.RawData) {
	//log.Println("new database task:", taskData.TagID)
	DatabaseTask <- taskData
}

//close connection on programm exit
func deferCleanup(db *sql.DB) {
	<-respawnLock
	err := db.Close() 
	if err != nil {
		log.Println("Error closing database connection:", err)
	}
}

func databaseWorkerRun(workerId int) {

	connection, err := connectToDb()

	defer deferCleanup(connection)

	if err != nil  {
		MyLog.Printonce(fmt.Sprintf("Database %s is unreachable. Error: %s",  Config.DB_TYPE, err))
		return

	} else {
		MyLog.Println(fmt.Sprintf("Database worker #%d connected to %s database", workerId, Config.DB_TYPE))
	}

	//initialise connection error channel
	connectionErrorChannel := make(chan error)

	go func() {
		var id int64
		for {
			err := connection.QueryRow("SELECT ID FROM DBWatchDog").Scan(&id)
			if err != nil {
				connectionErrorChannel <- err
				return
			} else {
				//do nothing (sleep):
				if id == 1 {
					time.Sleep(1 * time.Second)
				}
			}
		}
	}()

	for {
		select {
			//в случае если есть задание в канале DatabaseTask
		case currentDatabaseTask := <- DatabaseTask :
			log.Println(currentDatabaseTask.TagID)
			//do come task:
		case networkError := <-connectionErrorChannel :
			//обнаружена сетевая ошибка - завершаем гоурутину
			log.Printf("Database worker %d exited due to connection error: %s\n", workerId, networkError)
			return
		}
	}
}

