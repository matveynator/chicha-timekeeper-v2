package Config

import (
	"os"
	"fmt"
	//	"testing"
	"flag"
	"time"
	"hash/fnv"

)
var CompileVersion string

type Settings struct {
	APP_NAME, VERSION, DAEMON_LISTENER_ADDRESS, DAEMON_LISTENER_ADDRESS_HASH, WEB_LISTENER_ADDRESS, PROXY_ADDRESS, DB_TYPE, DB_FILE_PATH, DB_FULL_FILE_PATH, PG_HOST, PG_USER, PG_PASS, PG_DB_NAME, PG_SSL, TIME_ZONE, RACE_TYPE string
	PG_PORT int
	AVERAGE_RESULTS, VARIABLE_DISTANCE_RACE bool
	RACE_TIMEOUT_DURATION, MINIMAL_LAP_TIME_DURATION, AVERAGE_DURATION, DB_SAVE_INTERVAL_DURATION time.Duration
}

func isFlagPassed(name string) bool {
	found := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
}

func hash(s string) string {
	hash := fnv.New32a()
	hash.Write([]byte(s))
	return fmt.Sprint(hash.Sum32())
}

func ParseFlags() (config Settings)  { 
	config.APP_NAME = "chicha"
	flagVersion := flag.Bool("version", false, "Output version information")


	flag.StringVar(&config.RACE_TYPE, "race-type", "mass-start", "Valid race calculation variants are: 'delayed-start' or 'mass-start'. 1. 'mass-start': start time is not taken into account as everybody starts at the same time, the first gate passage is equal to the short lap, positions are counted based on the minimum time to complete maximum number of laps/stages/gates including the short lap. 2. 'delayed-start': start time is taken into account as everyone starts with some time delay, the first gate passage (short lap) is equal to the start time, positions are counted based on the minimum time to complete maximum number of laps/stages/gates excluding short lap.")

	flag.StringVar(&config.DAEMON_LISTENER_ADDRESS, "listener", "0.0.0.0:4000", "Please specify the IP address and port number on which the incoming RFID handler should operate.")
	flag.StringVar(&config.WEB_LISTENER_ADDRESS, "web", "0.0.0.0:80", "Please specify the IP address and port number on which HTTP web interface will be running.")
	flag.StringVar(&config.PROXY_ADDRESS, "proxy", "", "Proxy incoming data to another chicha listener. For example: -proxy '10.9.8.7:4000'.")
	flag.StringVar(&config.TIME_ZONE, "timezone", "UTC", "Set race timezone. Example: Europe/Paris, Africa/Dakar, UTC, https://en.wikipedia.org/wiki/List_of_tz_database_time_zones")
	flag.DurationVar(&config.RACE_TIMEOUT_DURATION, "timeout", 2*time.Minute, "Set race timeout duration. After this time if nobody passes the finish line the race will be stopped. Valid time units are: 's' (second), 'm' (minute), 'h' (hour).")
	flag.DurationVar(&config.MINIMAL_LAP_TIME_DURATION, "lap-time", 45*time.Second, "Minimal lap time duration refers to the shortest acceptable time for a lap, and any lap time result that falls below this duration would be deemed inaccurate." )
	flag.DurationVar(&config.AVERAGE_DURATION, "average-duration", 1000*time.Millisecond, "Duration to calculate average results. Any results provided to the reader within this duration will be computed as the average result." )
	flag.BoolVar(&config.AVERAGE_RESULTS, "average", false, "Calculate average results instead of minimal results.")
	flag.BoolVar(&config.VARIABLE_DISTANCE_RACE, "variable-distance-race", false, "A race in which each stage features differing distances.")

	//db
	flag.StringVar(&config.DB_FILE_PATH, "db-path", ".", "Provide path to writable directory to store database data.")
	flag.StringVar(&config.DB_TYPE, "db-type", "genji", "Select db type: sqlite / genji / postgres")
	flag.DurationVar(&config.DB_SAVE_INTERVAL_DURATION, "db-save-interval", 30000*time.Millisecond, "Duration to save data from memory to database (disk). Setting duration too low may cause unpredictable performance results." )

	//PostgreSQL related start
	flag.StringVar(&config.PG_HOST, "pg-host", "127.0.0.1", "PostgreSQL DB host.")
	flag.IntVar(&config.PG_PORT, "pg-port", 5432, "PostgreSQL DB port.")
	flag.StringVar(&config.PG_USER, "pg-user", "postgres", "PostgreSQL DB user.")
	flag.StringVar(&config.PG_PASS, "pg-pass", "", "PostgreSQL DB password.")
	flag.StringVar(&config.PG_DB_NAME, "pg-db-name", "chicha", "PostgreSQL DB name.")
	flag.StringVar(&config.PG_SSL, "pg-ssl", "prefer", "disable / allow / prefer / require / verify-ca / verify-full - PostgreSQL ssl modes: https://www.postgresql.org/docs/current/libpq-ssl.html")

	/*
	//escape test flags:
	var _ = func() bool { 
		testing.Init() 
		return true 
	}()
	*/

	//process all flags
	flag.Parse()

	//config.RACE_TYPE
	if config.RACE_TYPE != "mass-start" && config.RACE_TYPE != "delayed-start" {
		fmt.Printf("Error: race-type must be 'mass-start' or 'delayed-start', but you defined: %s. \n", config.RACE_TYPE)
		os.Exit(1)
	}

	//делаем хеш от порта коллектора чтобы использовать в уникальном названии файла бд
	config.DAEMON_LISTENER_ADDRESS_HASH = hash(config.DAEMON_LISTENER_ADDRESS)

	//путь к файлу бд
	config.DB_FULL_FILE_PATH = fmt.Sprintf(config.DB_FILE_PATH+"/"+config.APP_NAME+"."+config.DAEMON_LISTENER_ADDRESS_HASH+".db."+config.DB_TYPE)

	//set version from CompileVersion variable at build time
	config.VERSION = CompileVersion 

	if *flagVersion  {
		if config.VERSION != "" {
			fmt.Println("Version:", config.VERSION)
		}
		os.Exit(0)
	}

	// Startup banner START:
	fmt.Printf("Starting %s ", config.APP_NAME)
	if config.VERSION != "" {
		fmt.Printf("version %s ", config.VERSION)
	}
	fmt.Printf("at %s and web at %s, race type is \"%s\" and timezone is %s, minimal lap/stage duration is %s. ", config.DAEMON_LISTENER_ADDRESS, config.WEB_LISTENER_ADDRESS, config.RACE_TYPE, config.TIME_ZONE, config.MINIMAL_LAP_TIME_DURATION)

	if config.AVERAGE_RESULTS  {
		fmt.Printf("Calculating average time results. ")
	} else {
		fmt.Printf("Calculating minimal time results. ")
	}

	if config.VARIABLE_DISTANCE_RACE {
		fmt.Printf("Performing a race in which each stage features differing distances.\n")
	} else {
		fmt.Printf("Performing a race in which each stage features same distances (laps).\n")
	}
	// END.


	return
}
