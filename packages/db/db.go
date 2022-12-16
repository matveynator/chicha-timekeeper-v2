package Db

import (
	"fmt"
	"log"
	"chicha/packages/config"
	"database/sql"
)

type XMLData  struct {
	TagID                       string
	DiscoveryTime								string
	Antenna                     uint8
}

type AverageResult	struct {
	ID													int64
	RaceID											uint
	LapNumber										int
	TagID												string
	DiscoveryUnixTime   				int64
	Antenna             				uint8
	AntennaIP           				string
}

type Lap struct {
	ID                          int64
	SportsmanID                 string
	TagID                       string
	DiscoveryMinimalUnixTime    int64
	DiscoveryAverageUnixTime    int64
	UpdatedAt                   int64
	RaceID                      uint
	PracticeID                  uint
	RacePosition                uint
	TimeBehindTheLeader         int64
	LapNumber                   int
	LapTime                     int64
	LapPosition                 uint
	LapIsCurrent                bool
	LapIsStrange                bool
	RaceFinished                bool
	BestLapTime                 int64
	BestLapNumber               int
	BestLapPosition             uint
	RaceTotalTime               int64
	BetterOrWorseLapTime        int64
}


var Db *sql.DB

func init() {
	var err error
	if Config.DB_TYPE == "genji" {
		Db, err = sql.Open(Config.DB_TYPE, Config.DB_FILE_PATH+"/"+Config.APP_NAME+".db."+Config.DB_TYPE)
		if err != nil {
			Config.DB_TYPE = "sqlite"
			log.Println("Genji is unsupported on this architecture, switching to sqlite db type.")
			Db, err = sql.Open(Config.DB_TYPE, Config.DB_FILE_PATH+"/"+Config.APP_NAME+".db."+Config.DB_TYPE)
			if err != nil {
				log.Println("Database file error:", err)
				panic("Exiting due to configuration error.")
			}
		}
	} else if Config.DB_TYPE == "sqlite" {
		Db, err = sql.Open(Config.DB_TYPE, Config.DB_FILE_PATH+"/"+Config.APP_NAME+".db."+Config.DB_TYPE)
		if err != nil {
			Config.DB_TYPE = "genji"
			log.Println("SQLite is unsupported on this architecture, switching to genji db type.")
			Db, err = sql.Open(Config.DB_TYPE, Config.DB_FILE_PATH+"/"+Config.APP_NAME+".db."+Config.DB_TYPE)
			if err != nil {
				log.Println("Database file error:", err)
				panic("Exiting due to configuration error.")
			}
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

func UpdateLapInDB (lap Lap) (err error) {
	_, err = Db.Exec("UPDATE Lap SET(SportsmanID = ?, TagID = ?, DiscoveryMinimalUnixTime = ?, DiscoveryAverageUnixTime = ?, UpdatedAt = ?, RaceID = ?, PracticeID = ?, RacePosition = ?, TimeBehindTheLeader = ?, LapNumber = ?, LapTime = ?, LapPosition = ?, LapIsCurrent = ?, LapIsStrange = ?, RaceFinished = ?, BestLapTime = ?, BestLapNumber = ?, BestLapPosition = ?, RaceTotalTime = ?, BetterOrWorseLapTime = ?) WHERE ID = ?", lap.SportsmanID, lap.TagID, lap.DiscoveryMinimalUnixTime, lap.DiscoveryAverageUnixTime, lap.UpdatedAt, lap.RaceID, lap.PracticeID, lap.RacePosition, lap.TimeBehindTheLeader, lap.LapNumber, lap.LapTime, lap.LapPosition, lap.LapIsCurrent, lap.LapIsStrange, lap.RaceFinished, lap.BestLapTime, lap.BestLapNumber, lap.BestLapPosition, lap.RaceTotalTime, lap.BetterOrWorseLapTime, lap.ID)
	return
}

func InsertLapInDB (lap Lap) (id int64, err error) {
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

func SelectLapFromDB(oldLap Lap) (lap Lap, err error) {
	err = Db.QueryRow("SELECT ID,SportsmanID,TagID,DiscoveryMinimalUnixTime,DiscoveryAverageUnixTime,UpdatedAt,RaceID,PracticeID,RacePosition,TimeBehindTheLeader,LapNumber,LapTime,LapPosition,LapIsCurrent,LapIsStrange,RaceFinished,BestLapTime,BestLapNumber,BestLapPosition,RaceTotalTime,BetterOrWorseLapTime FROM Lap WHERE ID = ?", oldLap.ID).Scan(&lap.ID, &lap.SportsmanID, &lap.TagID, &lap.DiscoveryMinimalUnixTime,  &lap.DiscoveryAverageUnixTime, &lap.UpdatedAt,  &lap.RaceID,  &lap.PracticeID,  &lap.RacePosition,  &lap.TimeBehindTheLeader,  &lap.LapNumber,  &lap.LapTime,  &lap.LapPosition,  &lap.LapIsCurrent,  &lap.LapIsStrange,  &lap.RaceFinished,  &lap.BestLapTime,  &lap.BestLapNumber,  &lap.BestLapPosition,  &lap.RaceTotalTime,  &lap.BetterOrWorseLapTime)
	return
}


func GetCurrentRaceDataFromDB() (laps []Lap, err error) {
	var raceID int64
	err = Db.QueryRow("SELECT RaceID FROM Lap ORDER BY DiscoveryMinimalUnixTime").Scan(&raceID)
	if err == nil {
		rows, err := Db.Query("SELECT * FROM Lap WHERE RaceID = ?", raceID)
		defer rows.Close()
		if err == nil {
			// Loop through rows, using Scan to assign column data to struct fields.
			for rows.Next() {
				var lap Lap
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
	var lap Lap

	err := Db.QueryRow("SELECT ID,SportsmanID,TagID,DiscoveryMinimalUnixTime,DiscoveryAverageUnixTime,UpdatedAt,RaceID,PracticeID,RacePosition,TimeBehindTheLeader,LapNumber,LapTime,LapPosition,LapIsCurrent,LapIsStrange,RaceFinished,BestLapTime,BestLapNumber,BestLapPosition,RaceTotalTime,BetterOrWorseLapTime FROM Lap order by ID desc limit 1").Scan(&lap.ID, &lap.SportsmanID, &lap.TagID, &lap.DiscoveryMinimalUnixTime,  &lap.DiscoveryAverageUnixTime, &lap.UpdatedAt,  &lap.RaceID,  &lap.PracticeID,  &lap.RacePosition,  &lap.TimeBehindTheLeader,  &lap.LapNumber,  &lap.LapTime,  &lap.LapPosition,  &lap.LapIsCurrent,  &lap.LapIsStrange,  &lap.RaceFinished,  &lap.BestLapTime,  &lap.BestLapNumber,  &lap.BestLapPosition,  &lap.RaceTotalTime,  &lap.BetterOrWorseLapTime )
	if err != nil {
		log.Println(err)
	}
	fmt.Println(lap.TagID, lap.DiscoveryMinimalUnixTime)
}

