package Config

import (
	"os"
	"fmt"
//	"testing"
	"flag"
	"time"
	"hash/fnv"
	
)

type Settings struct {
	APP_NAME, VERSION, COLLECTOR_LISTENER_ADDRESS, COLLECTOR_LISTENER_ADDRESS_HASH, WEB_LISTENER_ADDRESS, PROXY_ADDRESS, DB_TYPE, DB_FILE_PATH, DB_FULL_FILE_PATH, PG_HOST, PG_USER, PG_PASS, PG_DB_NAME, PG_SSL, TIME_ZONE string
	PG_PORT int
	AVERAGE_RESULTS bool
	RACE_TIMEOUT_DURATION, MINIMAL_LAP_TIME_DURATION, AVERAGE_DURATION time.Duration
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
	flag.StringVar(&config.COLLECTOR_LISTENER_ADDRESS, "collector", "0.0.0.0:4000", "Provide IP address and port to collect and parse data from RFID and timing readers.")
	flag.StringVar(&config.WEB_LISTENER_ADDRESS, "web", "0.0.0.0:80", "Provide IP address and port to listen for HTTP connections from clients.")
	flag.StringVar(&config.PROXY_ADDRESS, "proxy", "", "Proxy received data to another collector. For example: -proxy '10.9.8.7:4000'.")
	flag.StringVar(&config.TIME_ZONE, "timezone", "UTC", "Set race timezone. Example: Europe/Paris, Africa/Dakar, UTC, https://en.wikipedia.org/wiki/List_of_tz_database_time_zones")
	flag.DurationVar(&config.RACE_TIMEOUT_DURATION, "timeout", 2*time.Minute, "Set race timeout duration. After this time if nobody passes the finish line the race will be stopped. Valid time units are: 's' (second), 'm' (minute), 'h' (hour).")
	flag.DurationVar(&config.MINIMAL_LAP_TIME_DURATION, "laptime", 45*time.Second, "Minimal lap time duration. Results smaller than this duration would be considered wrong." )
	  flag.DurationVar(&config.AVERAGE_DURATION, "average-duration", 1000*time.Millisecond, "Duration to calculate average results. Results passed to reader during this duration will be calculated as average result." )
	flag.BoolVar(&config.AVERAGE_RESULTS, "average", true, "Calculate average results instead of only first results.")

	//db
	flag.StringVar(&config.DB_FILE_PATH, "dbpath", ".", "Provide path to writable directory to store database data.")
	flag.StringVar(&config.DB_TYPE, "dbtype", "sqlite", "Select db type: sqlite / genji / postgres")

	//PostgreSQL related start
	flag.StringVar(&config.PG_HOST, "pghost", "127.0.0.1", "PostgreSQL DB host.")
	flag.IntVar(&config.PG_PORT, "pgport", 5432, "PostgreSQL DB port.")
	flag.StringVar(&config.PG_USER, "pguser", "postgres", "PostgreSQL DB user.")
	flag.StringVar(&config.PG_PASS, "pgpass", "", "PostgreSQL DB password.")
	flag.StringVar(&config.PG_DB_NAME, "pgdbname", "chicha", "PostgreSQL DB name.")
	flag.StringVar(&config.PG_SSL, "pgssl", "prefer", "disable / allow / prefer / require / verify-ca / verify-full - PostgreSQL ssl modes: https://www.postgresql.org/docs/current/libpq-ssl.html")

/*
	//escape test flags:
	var _ = func() bool { 
		testing.Init() 
		return true 
	}()
*/

	//process all flags
	flag.Parse()

	//делаем хеш от порта коллектора чтобы использовать в уникальном названии файла бд
	config.COLLECTOR_LISTENER_ADDRESS_HASH = hash(config.COLLECTOR_LISTENER_ADDRESS)

	//путь к файлу бд
	config.DB_FULL_FILE_PATH = fmt.Sprintf(config.DB_FILE_PATH+"/"+config.APP_NAME+"."+config.COLLECTOR_LISTENER_ADDRESS_HASH+".db."+config.DB_TYPE)

	if *flagVersion  {
		if config.VERSION != "" {
			fmt.Println("Version:", config.VERSION)
		}
		os.Exit(0)
	}
	return
}
