package main

import (
	"log"
	"fmt"
	"database/sql"
	_ "modernc.org/sqlite"
)

var databaseURI = "chicha.db.sqlite"

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

	db, err := sql.Open("sqlite", databaseURI)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()


	lap.SportsmanID = "TESTRIDER00061453226008"
  lap.DiscoveryMinimalUnixTime = 1670935623391

	query := `
	CREATE TABLE IF NOT EXISTS results(id INTEGER PRIMARY KEY, rider_id TEXT, date BIGINT);
	INSERT INTO results(rider_id, date) VALUES('TESTRIDER00061453226008', 1670935623391);
	INSERT INTO results(rider_id, date) VALUES('Volkswagen',21600);
	`
	_, err = db.Exec(query)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("table created")


}

