package Database

import (
	"fmt"
	"log"
	"errors"
	"time"

	"database/sql"

	"chicha/pkg/config"
	"chicha/pkg/data"
)


func connectToDb(config Config.Settings)(db *sql.DB, err error) {
	if config.DB_TYPE == "genji" {
		db, err = sql.Open(config.DB_TYPE, config.DB_FULL_FILE_PATH)
		if err != nil {
			config.DB_TYPE = "sqlite"
			log.Println("Database error:", err)
			log.Println("Genji is unsupported on this architecture, switching to sqlite db type.")
			db, err = sql.Open(config.DB_TYPE, config.DB_FULL_FILE_PATH)
			if err != nil {
				err = errors.New(fmt.Sprintf("Database file error: %s", err.Error()))
				log.Println(err)
				log.Println("SQLite is unsupported on this architecture, please use: -dbtype postgres.")
				return
			} else {
				err = createTables(db, config)
				if err != nil {
					err = errors.New(fmt.Sprintf("Database create tables error: %s", err.Error()))
					log.Println(err)
					return
				}
			}
		} else {
			err = createTables(db, config)
			if err != nil {
				err = errors.New(fmt.Sprintf("Database create tables error: %s", err.Error()))
				log.Println(err)
				return
			}
		}
	} else if config.DB_TYPE == "sqlite" {
		db, err = sql.Open(config.DB_TYPE, config.DB_FULL_FILE_PATH)
		if err != nil {
			config.DB_TYPE = "genji"
			log.Println("Database file error:", err)
			log.Println("SQLite is unsupported on this architecture, switching to genji db type.")
			db, err = sql.Open(config.DB_TYPE, config.DB_FULL_FILE_PATH)
			if err != nil {
				err = errors.New(fmt.Sprintf("Database file error: %s", err.Error()))
				log.Println(err)
				log.Println("Genji is unsupported on this architecture, please use: -dbtype postgres.")
				return
			} else {
				err = createTables(db, config)
				if err != nil {
					err = errors.New(fmt.Sprintf("Database create tables error: %s", err.Error()))
					log.Println(err)
					return
				}
			}
		} else {
			err = createTables(db, config)
			if err != nil {
				err = errors.New(fmt.Sprintf("Database create tables error: %s", err.Error()))
				log.Println(err)
				return
			}
		}
	} else if config.DB_TYPE == "postgres" {

		psqlConnectDSN := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s pool_max_conns=10", config.PG_HOST, config.PG_PORT, config.PG_USER, config.PG_PASS, config.PG_DB_NAME, config.PG_SSL)
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
			err = createTables(db, config)
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


func createTables(db *sql.DB, config Config.Settings) (err error) {

	_, err = db.Exec("CREATE TABLE if not exists DBWatchDog(Id INT PRIMARY KEY, UnixTime INT)")
	if err != nil {
		return
	} else {
		//populate DBWatchDog with data (one row with only one Id=1)
		var id int64
		// Create a sql/database DB instance
		err = db.QueryRow("SELECT Id FROM DBWatchDog").Scan(&id)
		if err != nil  {
			_, err = db.Exec("INSERT INTO DBWatchDog (Id,UnixTime) VALUES (?,?)", 1, time.Now().UnixMilli())
			if err != nil {
				return
			} else {
				log.Printf("Created new %s database file: %s \n", config.DB_TYPE, config.DB_FULL_FILE_PATH)
			}
		}
	}

	_, err = db.Exec("CREATE TABLE if not exists Lap(Id INT PRIMARY KEY, SportsmanId INT DEFAULT 0, TagId TEXT, DiscoveryMinimalUnixTime INT DEFAULT 0, DiscoveryAverageUnixTime INT DEFAULT 0, UpdatedAt INT DEFAULT 0, RaceId INT DEFAULT 0, PracticeId INT DEFAULT 0, RacePosition INT DEFAULT 0, TimeBehindTheLeader INT DEFAULT 0, LapNumber INT DEFAULT 0, LapTime INT DEFAULT 0, LapPosition INT DEFAULT 0, LapIsCurrent BOOL DEFAULT FALSE, LapIsStrange BOOL DEFAULT FALSE, RaceFinished BOOL DEFAULT FALSE, BestLapTime INT DEFAULT 0, BestLapNumber INT DEFAULT 0, BestLapPosition INT DEFAULT 0, RaceTotalTime INT DEFAULT 0, BetterOrWorseLapTime INT DEFAULT 0, UNIQUE(Id))")
	if err != nil {
		return
	}

	_, err = db.Exec("CREATE TABLE if not exists RawData(Id INT NOT NULL, TagId TEXT, DiscoveryUnixTime INT,  ReaderIP TEXT, Antenna INT, ProxyIP TEXT, PRIMARY KEY(Id))")
	if err != nil {
		return
	}

	return
}


func InsertRawDataInDB (databaseConnection *sql.DB, rawData Data.RawData) (id int64, err error) {
	err = databaseConnection.QueryRow("SELECT Id FROM RawData order by Id desc limit 1").Scan(&id)
	if err != nil {
		id=1
	} else {
		id++
	}

	_, err = databaseConnection.Exec("INSERT INTO RawData(Id,TagId,DiscoveryUnixTime,ReaderIP,Antenna,ProxyIP) VALUES (?, ?, ?, ?, ?, ?)", id, rawData.TagId, rawData.DiscoveryUnixTime, rawData.ReaderIP, rawData.Antenna, rawData.ProxyIP)
	return
}
