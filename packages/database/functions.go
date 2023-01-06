package Database

import (
	"fmt"
	"log"
	"errors"
	"time"

	"database/sql"

	"chicha/packages/config"
	"chicha/packages/data"
)


func connectToDb()(db *sql.DB, err error) {
	if Config.DB_TYPE == "genji" {
		db, err = sql.Open(Config.DB_TYPE, Config.DB_FULL_FILE_PATH)
		if err != nil {
			Config.DB_TYPE = "sqlite"
			log.Println("Database error:", err)
			log.Println("Genji is unsupported on this architecture, switching to sqlite db type.")
			db, err = sql.Open(Config.DB_TYPE, Config.DB_FULL_FILE_PATH)
			if err != nil {
				err = errors.New(fmt.Sprintf("Database file error: %s", err.Error()))
				log.Println(err)
				log.Println("SQLite is unsupported on this architecture, please use: -dbtype postgres.")
				return
			} else {
				err = createTables(db)
				if err != nil {
					err = errors.New(fmt.Sprintf("Database create tables error: %s", err.Error()))
					log.Println(err)
					return
				}
			}
		} else {
			err = createTables(db)
			if err != nil {
				err = errors.New(fmt.Sprintf("Database create tables error: %s", err.Error()))
				log.Println(err)
				return
			}
		}
	} else if Config.DB_TYPE == "sqlite" {
		db, err = sql.Open(Config.DB_TYPE, Config.DB_FULL_FILE_PATH)
		if err != nil {
			Config.DB_TYPE = "genji"
			log.Println("Database file error:", err)
			log.Println("SQLite is unsupported on this architecture, switching to genji db type.")
			db, err = sql.Open(Config.DB_TYPE, Config.DB_FULL_FILE_PATH)
			if err != nil {
				err = errors.New(fmt.Sprintf("Database file error: %s", err.Error()))
				log.Println(err)
				log.Println("Genji is unsupported on this architecture, please use: -dbtype postgres.")
				return
			} else {
				err = createTables(db)
				if err != nil {
					err = errors.New(fmt.Sprintf("Database create tables error: %s", err.Error()))
					log.Println(err)
					return
				}
			}
		} else {
			err = createTables(db)
			if err != nil {
				err = errors.New(fmt.Sprintf("Database create tables error: %s", err.Error()))
				log.Println(err)
				return
			}
		}
	} else if Config.DB_TYPE == "postgres" {

		psqlConnectDSN := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s pool_max_conns=10", Config.PG_HOST, Config.PG_PORT, Config.PG_USER, Config.PG_PASS, Config.PG_DB_NAME, Config.PG_SSL)
		db, err = sql.Open("pgx", psqlConnectDSN)
		if err != nil {
			err = errors.New(fmt.Sprintf("Database config error: %s", err.Error()))
			log.Println(err)
			return
		}
		err = db.Ping()
		if err != nil {
			err = errors.New(fmt.Sprintf("Database connect error: %s", err.Error()))
			log.Println(err)
			return
		} else {
			err = createTables(db)
			if err != nil {
				err = errors.New(fmt.Sprintf("Database create tables error: %s", err.Error()))
				log.Println(err)
				return
			}
		}
	} else {
		err = errors.New("Please define valid dbtype (genji / sqlite)")
		log.Println(err)
		return
	}
	return
}


func createTables(db *sql.DB) (err error) {

	_, err = db.Exec("CREATE TABLE if not exists DBWatchDog(ID INT PRIMARY KEY, UnixTime INT)")
	if err != nil {
		return
	} else {
		//populate DBWatchDog with data (one row with only one ID=1)
		var id int64
		// Create a sql/database DB instance
		err = db.QueryRow("SELECT ID FROM DBWatchDog").Scan(&id)
		if err != nil  {
			_, err = db.Exec("INSERT INTO DBWatchDog (ID,UnixTime) VALUES (?,?)", 1, time.Now().UnixMilli())
			if err != nil {
				return
			} else {
				log.Printf("Created new %s database file: %s \n", Config.DB_TYPE, Config.DB_FULL_FILE_PATH)
			}
		}
	}

	_, err = db.Exec("CREATE TABLE if not exists Lap(ID INT PRIMARY KEY, SportsmanID INT DEFAULT 0, TagID TEXT, DiscoveryMinimalUnixTime INT DEFAULT 0, DiscoveryAverageUnixTime INT DEFAULT 0, UpdatedAt INT DEFAULT 0, RaceID INT DEFAULT 0, PracticeID INT DEFAULT 0, RacePosition INT DEFAULT 0, TimeBehindTheLeader INT DEFAULT 0, LapNumber INT DEFAULT 0, LapTime INT DEFAULT 0, LapPosition INT DEFAULT 0, LapIsCurrent BOOL DEFAULT FALSE, LapIsStrange BOOL DEFAULT FALSE, RaceFinished BOOL DEFAULT FALSE, BestLapTime INT DEFAULT 0, BestLapNumber INT DEFAULT 0, BestLapPosition INT DEFAULT 0, RaceTotalTime INT DEFAULT 0, BetterOrWorseLapTime INT DEFAULT 0, UNIQUE(ID))")
	if err != nil {
		return
	}

	_, err = db.Exec("CREATE TABLE if not exists RawData(ID INT NOT NULL, TagID TEXT, DiscoveryUnixTime INT,  ReaderIP TEXT, Antenna INT, ProxyIP TEXT, PRIMARY KEY(ID))")
	if err != nil {
		return
	}

	return
}


func InsertRawDataInDB (databaseConnection *sql.DB, rawData Data.RawData) (id int64, err error) {
	err = databaseConnection.QueryRow("SELECT ID FROM RawData order by ID desc limit 1").Scan(&id)
	if err != nil {
		id=1
	} else {
		id++
	}

	_, err = databaseConnection.Exec("INSERT INTO RawData(ID,TagID,DiscoveryUnixTime,ReaderIP,Antenna,ProxyIP) VALUES (?, ?, ?, ?, ?, ?)", id, rawData.TagID, rawData.DiscoveryUnixTime, rawData.ReaderIP, rawData.Antenna, rawData.ProxyIP)
	return
}
