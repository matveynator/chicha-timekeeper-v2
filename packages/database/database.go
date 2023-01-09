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

var DatabaseTask chan Data.RawData
var respawnLock chan int
//по умолчанию оставляем только один процесс который будет брать задачи и записывать их в базу данных
var databaseWorkersMaxCount int = 1

func Run(config Config.Settings) {
	
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
			go databaseWorkerRun(len(respawnLock), config)
		}
	}()
}

func CreateNewDatabaseTask(taskData Data.RawData) {
	//log.Println("new database task:", taskData.TagID)
	DatabaseTask <- taskData
}

//close dbConnection on programm exit
func deferCleanup(db *sql.DB) {
	<-respawnLock
	err := db.Close() 
	if err != nil {
		log.Println("Error closing database connection:", err)
	}
}

func databaseWorkerRun(workerId int, config Config.Settings ) {
	

	dbConnection, err := connectToDb(config)

	defer deferCleanup(dbConnection)

	if err != nil  {
		MyLog.Printonce(fmt.Sprintf("Database %s is unreachable. Error: %s",  config.DB_TYPE, err))
		return

	} else {
		MyLog.Println(fmt.Sprintf("Database worker #%d connected to %s database", workerId, config.DB_TYPE))
	}

	//initialise dbConnection error channel
	connectionErrorChannel := make(chan error)

	go func() {
		for {
			_, err = dbConnection.Exec("UPDATE DBWatchDog SET UnixTime = ? WHERE ID = 1", time.Now().UnixMilli())
			if err != nil {
				connectionErrorChannel <- err
				return
			} else {
				//log.Println("Database is alive.")
				time.Sleep(1 * time.Second)
			}
		}
	}()

	for {
		select {
			//в случае если есть задание в канале DatabaseTask
		case currentDatabaseTask := <- DatabaseTask :
			//log.Println("Received new database task with TagID:", currentDatabaseTask.TagID)
			_, err := InsertRawDataInDB(dbConnection, currentDatabaseTask)
			if err != nil {
				log.Printf("Database worker %d exited due to processing error: %s\n", workerId, err)
				return
			}
			//do come task:
		case networkError := <-connectionErrorChannel :
			//обнаружена сетевая ошибка - завершаем гоурутину
			log.Printf("Database worker %d exited due to connection error: %s\n", workerId, networkError)
			return
		}
	}
}

