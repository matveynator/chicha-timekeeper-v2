package main

import (
	"log"
	"fmt"
	"database/sql"
	_ "modernc.org/sqlite"
	"chicha/packages/config"
)


type Lap struct {
	ID                          int64
	SportsmanID                 string
	TagID                       string
	DiscoveryMinimalUnixTime    int64
	DiscoveryAverageUnixTime    int64
	UpdatedAt                   int64
	RaceID                      uint
	RacePosition         				uint
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


func main() {
  var lap Lap

	db, err := sql.Open("sqlite", Config.DB_FILE_PATH)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()


	lap.SportsmanID = "TESTRIDER00061453226008"
  lap.DiscoveryMinimalUnixTime = 1670935623391

	query := `
	CREATE TABLE if not exists Lap(ID INT PRIMARY KEY, SportsmanID INT DEFAULT 0, TagID TEXT, DiscoveryMinimalUnixTime INT DEFAULT 0, DiscoveryAverageUnixTime INT DEFAULT 0, UpdatedAt INT DEFAULT 0, RaceID INT DEFAULT 0, PracticeID INT DEFAULT 0, CurrentRacePosition INT DEFAULT 0, TimeBehindTheLeader INT DEFAULT 0, LapNumber INT DEFAULT 0, LapTime INT DEFAULT 0, LapPosition INT DEFAULT 0, LapIsCurrent BOOL DEFAULT FALSE, LapIsStrange BOOL DEFAULT FALSE, RaceFinished BOOL DEFAULT FALSE, BestLapTime INT DEFAULT 0, BestLapNumber INT DEFAULT 0, BestLapPosition INT DEFAULT 0, RaceTotalTime INT DEFAULT 0, BetterOrWorseLapTime INT DEFAULT 0, UNIQUE(ID));
	CREATE TABLE if not exists AverageResults(ID INT PRIMARY KEY, LapID INT DEFAULT 0, TagID TEXT, DiscoveryUnixTime INT DEFAULT 0, Antenna INT DEFAULT 0, AntennaIP TEXT DEFAULT 0, UNIQUE(ID));

	INSERT INTO Lap(SportsmanID, DiscoveryAverageUnixTime) VALUES('TESTRIDER00061453226008', 1670935623391);
	`
	_, err = db.Exec(query)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("table created")


}

