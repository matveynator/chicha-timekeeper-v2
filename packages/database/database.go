package Database

import (
	"fmt"
	"log"
	"time"
	"database/sql"


	"chicha/packages/config"
	"chicha/packages/data"
)

var Db *sql.DB

// Buffer for new RFID requests
var laps []Data.Lap

//incoming buffer for laps from readers
var LapsTaskChannel = make(chan Data.Lap)

// channel lockers
var lapsChannelBufferLocker = make(chan int, 1)
var lapsChannelDBLocker = make(chan int, 1)

func init() {
	var err error
	//unlock db operations
	lapsChannelDBLocker <- 0 //Put the initial value into the channel to unlock operations

	if Config.DB_TYPE == "genji" {
		Db, err = sql.Open(Config.DB_TYPE, Config.DB_FILE_PATH+"/"+Config.APP_NAME+".db."+Config.DB_TYPE)
		if err != nil {
			Config.DB_TYPE = "sqlite"
			log.Println("Database error:", err)
			log.Println("Genji is unsupported on this architecture, switching to sqlite db type.")
			Db, err = sql.Open(Config.DB_TYPE, Config.DB_FILE_PATH+"/"+Config.APP_NAME+".db."+Config.DB_TYPE)
			if err != nil {
				log.Println("Database file error:", err)
				log.Println("SQLite is unsupported on this architecture, please use: -dbtype postgres.")
				panic("Exiting due to configuration error.")
			}
		}
	} else if Config.DB_TYPE == "sqlite" {
		Db, err = sql.Open(Config.DB_TYPE, Config.DB_FILE_PATH+"/"+Config.APP_NAME+".db."+Config.DB_TYPE)
		if err != nil {
			Config.DB_TYPE = "genji"
			log.Println("Database file error:", err)
			log.Println("SQLite is unsupported on this architecture, switching to genji db type.")
			Db, err = sql.Open(Config.DB_TYPE, Config.DB_FILE_PATH+"/"+Config.APP_NAME+".db."+Config.DB_TYPE)
			if err != nil {
				log.Println("Database file error:", err)
				log.Println("Genji is unsupported on this architecture, please use: -dbtype postgres.")
				panic("Exiting due to configuration error.")
			}
		}
	} else if Config.DB_TYPE == "postgres" {

		psqlConnectDSN := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s pool_max_conns=10", Config.PG_HOST, Config.PG_PORT, Config.PG_USER, Config.PG_PASS, Config.PG_DB_NAME, Config.PG_SSL)
		Db, err = sql.Open("pgx", psqlConnectDSN)
		if err != nil {
			log.Println("Database config error:", err)
			panic("Exiting due to configuration error.")
		}
		err = Db.Ping()
		if err != nil {
			log.Println("Database connect error:", err)
			panic("Exiting due to configuration error.")
		}
	} else {
		log.Println("Please define valid dbtype (genji / sqlite)")
		panic("Exiting due to configuration error.")
	}
}

func CreateDB() {
	// Create a table. Genji tables are schemaless by default, you don't need to specify a schema.
	_, err := Db.Exec("CREATE TABLE if not exists Lap(ID INT PRIMARY KEY, SportsmanID INT DEFAULT 0, TagID TEXT, DiscoveryMinimalUnixTime INT DEFAULT 0, DiscoveryAverageUnixTime INT DEFAULT 0, UpdatedAt INT DEFAULT 0, RaceID INT DEFAULT 0, PracticeID INT DEFAULT 0, RacePosition INT DEFAULT 0, TimeBehindTheLeader INT DEFAULT 0, LapNumber INT DEFAULT 0, LapTime INT DEFAULT 0, LapPosition INT DEFAULT 0, LapIsCurrent BOOL DEFAULT FALSE, LapIsStrange BOOL DEFAULT FALSE, RaceFinished BOOL DEFAULT FALSE, BestLapTime INT DEFAULT 0, BestLapNumber INT DEFAULT 0, BestLapPosition INT DEFAULT 0, RaceTotalTime INT DEFAULT 0, BetterOrWorseLapTime INT DEFAULT 0, UNIQUE(ID))")
	if err != nil {
		log.Println(err)
	}

	_, err = Db.Exec("CREATE TABLE if not exists AverageResults(ID INT PRIMARY KEY, LapID INT DEFAULT 0, TagID TEXT, DiscoveryUnixTime INT DEFAULT 0, Antenna INT DEFAULT 0, AntennaIP TEXT DEFAULT 0, UNIQUE(ID))")
	if err != nil {
		log.Println(err)
	}
}

func UpdateDB() { 
	var id int64
	// Create a sql/database DB instance
	err := Db.QueryRow("SELECT ID FROM Lap order by ID desc limit 1").Scan(&id)
	if err != nil {
		log.Println(err)
		id=1
		fmt.Println("id=",id)
	} else {
		id=id+1
		fmt.Println("id=",id)
	}
	_, err = Db.Exec("INSERT INTO Lap (ID, TagID, DiscoveryMinimalUnixTime) VALUES (?, ?, ?)", id, "100008020200000100000425", 1636112997241)
	if err != nil {
		log.Println(err)
	}
}

func UpdateLapInDB (lap Data.Lap) (err error) {
	_, err = Db.Exec("UPDATE Lap SET(SportsmanID = ?, TagID = ?, DiscoveryMinimalUnixTime = ?, DiscoveryAverageUnixTime = ?, UpdatedAt = ?, RaceID = ?, PracticeID = ?, RacePosition = ?, TimeBehindTheLeader = ?, LapNumber = ?, LapTime = ?, LapPosition = ?, LapIsCurrent = ?, LapIsStrange = ?, RaceFinished = ?, BestLapTime = ?, BestLapNumber = ?, BestLapPosition = ?, RaceTotalTime = ?, BetterOrWorseLapTime = ?) WHERE ID = ?", lap.SportsmanID, lap.TagID, lap.DiscoveryMinimalUnixTime, lap.DiscoveryAverageUnixTime, lap.UpdatedAt, lap.RaceID, lap.PracticeID, lap.RacePosition, lap.TimeBehindTheLeader, lap.LapNumber, lap.LapTime, lap.LapPosition, lap.LapIsCurrent, lap.LapIsStrange, lap.RaceFinished, lap.BestLapTime, lap.BestLapNumber, lap.BestLapPosition, lap.RaceTotalTime, lap.BetterOrWorseLapTime, lap.ID)
	return
}

func InsertLapInDB (lap Data.Lap) (id int64, err error) {
	err = Db.QueryRow("SELECT ID FROM Lap order by ID desc limit 1").Scan(&id)
	if err == nil {
		//auto increment ID:
		id = id + 1
	} else {
		//not found? create first ID:
		id = 1
	}
	lap.ID=id;
	_, err = Db.Exec("INSERT INTO Lap(ID,SportsmanID,TagID,DiscoveryMinimalUnixTime,DiscoveryAverageUnixTime,UpdatedAt,RaceID,PracticeID,RacePosition,TimeBehindTheLeader,LapNumber,LapTime,LapPosition,LapIsCurrent,LapIsStrange,RaceFinished,BestLapTime,BestLapNumber,BestLapPosition,RaceTotalTime,BetterOrWorseLapTime) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", lap.ID, lap.SportsmanID, lap.TagID, lap.DiscoveryMinimalUnixTime, lap.DiscoveryAverageUnixTime, lap.UpdatedAt, lap.RaceID, lap.PracticeID, lap.RacePosition, lap.TimeBehindTheLeader, lap.LapNumber, lap.LapTime, lap.LapPosition, lap.LapIsCurrent, lap.LapIsStrange, lap.RaceFinished, lap.BestLapTime, lap.BestLapNumber, lap.BestLapPosition, lap.RaceTotalTime, lap.BetterOrWorseLapTime)
	return
}

func SelectLapFromDB(oldLap Data.Lap) (lap Data.Lap, err error) {
	err = Db.QueryRow("SELECT ID,SportsmanID,TagID,DiscoveryMinimalUnixTime,DiscoveryAverageUnixTime,UpdatedAt,RaceID,PracticeID,RacePosition,TimeBehindTheLeader,LapNumber,LapTime,LapPosition,LapIsCurrent,LapIsStrange,RaceFinished,BestLapTime,BestLapNumber,BestLapPosition,RaceTotalTime,BetterOrWorseLapTime FROM Lap WHERE ID = ?", oldLap.ID).Scan(&lap.ID, &lap.SportsmanID, &lap.TagID, &lap.DiscoveryMinimalUnixTime,  &lap.DiscoveryAverageUnixTime, &lap.UpdatedAt,  &lap.RaceID,  &lap.PracticeID,  &lap.RacePosition,  &lap.TimeBehindTheLeader,  &lap.LapNumber,  &lap.LapTime,  &lap.LapPosition,  &lap.LapIsCurrent,  &lap.LapIsStrange,  &lap.RaceFinished,  &lap.BestLapTime,  &lap.BestLapNumber,  &lap.BestLapPosition,  &lap.RaceTotalTime,  &lap.BetterOrWorseLapTime)
	return
}


func GetCurrentRaceDataFromDB() (laps []Data.Lap, err error) {
	var raceID int64
	err = Db.QueryRow("SELECT RaceID FROM Lap ORDER BY DiscoveryMinimalUnixTime").Scan(&raceID)
	if err == nil {
		rows, err := Db.Query("SELECT * FROM Lap WHERE RaceID = ?", raceID)
		defer rows.Close()
		if err == nil {
			// Loop through rows, using Scan to assign column data to struct fields.
			for rows.Next() {
				var lap Data.Lap
				err = rows.Scan(&lap.ID, &lap.SportsmanID, &lap.TagID, &lap.DiscoveryMinimalUnixTime,  &lap.DiscoveryAverageUnixTime, &lap.UpdatedAt,  &lap.RaceID,  &lap.PracticeID,  &lap.RacePosition,  &lap.TimeBehindTheLeader,  &lap.LapNumber,  &lap.LapTime,  &lap.LapPosition,  &lap.LapIsCurrent,  &lap.LapIsStrange,  &lap.RaceFinished,  &lap.BestLapTime,  &lap.BestLapNumber,  &lap.BestLapPosition,  &lap.RaceTotalTime,  &lap.BetterOrWorseLapTime)
				if err == nil {
					laps = append(laps, lap)
				}
			}
		}
	}
	return
}


func SelectFromDB() {
	var lap Data.Lap

	err := Db.QueryRow("SELECT ID,SportsmanID,TagID,DiscoveryMinimalUnixTime,DiscoveryAverageUnixTime,UpdatedAt,RaceID,PracticeID,RacePosition,TimeBehindTheLeader,LapNumber,LapTime,LapPosition,LapIsCurrent,LapIsStrange,RaceFinished,BestLapTime,BestLapNumber,BestLapPosition,RaceTotalTime,BetterOrWorseLapTime FROM Lap order by ID desc limit 1").Scan(&lap.ID, &lap.SportsmanID, &lap.TagID, &lap.DiscoveryMinimalUnixTime,  &lap.DiscoveryAverageUnixTime, &lap.UpdatedAt,  &lap.RaceID,  &lap.PracticeID,  &lap.RacePosition,  &lap.TimeBehindTheLeader,  &lap.LapNumber,  &lap.LapTime,  &lap.LapPosition,  &lap.LapIsCurrent,  &lap.LapIsStrange,  &lap.RaceFinished,  &lap.BestLapTime,  &lap.BestLapNumber,  &lap.BestLapPosition,  &lap.RaceTotalTime,  &lap.BetterOrWorseLapTime )
	if err != nil {
		log.Println(err)
	}
	fmt.Println(lap.TagID, lap.DiscoveryMinimalUnixTime)
}


func SaveLapsBufferSimplyToDB() {

	//loop forever:
	for range time.Tick(1 * time.Second) {
		<-lapsChannelDBLocker //grab the ticket via channel (lock)
		log.Println("Saving buffer to database started.")
		for _, lap := range laps {
			var newLap Data.Lap

			err := Db.QueryRow("SELECT ID,SportsmanID,TagID,DiscoveryMinimalUnixTime,DiscoveryAverageUnixTime,UpdatedAt,RaceID,PracticeID,RacePosition,TimeBehindTheLeader,LapNumber,LapTime,LapPosition,LapIsCurrent,LapIsStrange,RaceFinished,BestLapTime,BestLapNumber,BestLapPosition,RaceTotalTime,BetterOrWorseLapTime FROM Lap WHERE TagID = ? and RaceID = ? and LapNumber = ? limit 1;", newLap.TagID, newLap.RaceID, newLap.LapNumber).Scan(&newLap.ID, &newLap.SportsmanID, &newLap.TagID, &newLap.DiscoveryMinimalUnixTime,  &newLap.DiscoveryAverageUnixTime, &newLap.UpdatedAt,  &newLap.RaceID,  &newLap.PracticeID,  &newLap.RacePosition,  &newLap.TimeBehindTheLeader,  &newLap.LapNumber,  &newLap.LapTime,  &newLap.LapPosition,  &newLap.LapIsCurrent,  &newLap.LapIsStrange,  &newLap.RaceFinished,  &newLap.BestLapTime,  &newLap.BestLapNumber,  &newLap.BestLapPosition,  &newLap.RaceTotalTime,  &newLap.BetterOrWorseLapTime)

			log.Printf("race_id = %d, lap_number = %d, tag_id = %s, discovery_unix_time = %d \n", newLap.RaceID, newLap.LapNumber, newLap.TagID, newLap.DiscoveryMinimalUnixTime);
			if err == nil {
				//found old data - just update it:

				newLap.SportsmanID = lap.SportsmanID
				newLap.TagID = lap.TagID
				newLap.DiscoveryMinimalUnixTime =  lap.DiscoveryMinimalUnixTime
				newLap.DiscoveryAverageUnixTime = lap.DiscoveryAverageUnixTime
				newLap.AverageResultsCount = lap.AverageResultsCount
				newLap.Antenna  = lap.Antenna
				newLap.AntennaIP  = lap.AntennaIP
				newLap.UpdatedAt  = lap.UpdatedAt
				newLap.RaceID = lap.RaceID
				newLap.RacePosition = lap.RacePosition
				newLap.TimeBehindTheLeader = lap.TimeBehindTheLeader
				newLap.LapNumber = lap.LapNumber
				newLap.LapTime = lap.LapTime
				newLap.LapPosition = lap.LapPosition
				newLap.LapIsCurrent = lap.LapIsCurrent
				newLap.LapIsStrange = lap.LapIsStrange
				newLap.StageFinished = lap.StageFinished
				newLap.BestLapTime = lap.BestLapTime
				newLap.BestLapNumber = lap.BestLapNumber
				newLap.BestLapPosition = lap.BestLapPosition
				newLap.RaceTotalTime = lap.RaceTotalTime
				newLap.BetterOrWorseLapTime = lap.BetterOrWorseLapTime

				err := UpdateLapInDB(newLap)
				if err != nil {
					log.Println("Error. Not updated in database:", err)
				}
			} else {
				log.Println("Data not found in database:", err)
				//not found - create new
				_, err := InsertLapInDB(lap)
				if err != nil {
					log.Println("Error. Not created new data in database:", err)
				}
			}
		}

		lapsChannelDBLocker <- 1 //give ticket back via channel (unlock)
	}
}

