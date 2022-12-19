package Config

import (
	"os"
	"fmt"
	"flag"
	"time"
)

var APP_NAME = "chicha"
var VERSION, RFID_LISTENER_ADDRESS, WEB_LISTENER_ADDRESS, REDIRECT_ADDRESS, DB_TYPE, DB_FILE_PATH, PG_HOST, PG_USER, PG_PASS, PG_DB_NAME, PG_SSL, TIME_ZONE string
var PG_PORT int
var AVERAGE_RESULTS bool
var RACE_TIMEOUT_DURATION, MINIMAL_LAP_TIME_DURATION time.Duration


func isFlagPassed(name string) bool {
	found := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
}


func init()  {

	flagVersion := flag.Bool("version", false, "Output version information")
	flag.StringVar(&RFID_LISTENER_ADDRESS, "rfid", "0.0.0.0:4000", "Provide IP address and port to listen for RFID data.")
	flag.StringVar(&WEB_LISTENER_ADDRESS, "web", "0.0.0.0:80", "Provide IP address and port to listen for HTTP connections.")
	flag.StringVar(&REDIRECT_ADDRESS, "redirect", "", "Redirect rfid data stream to another instance. For example: -redirect '10.9.8.7:4000'.")
	flag.StringVar(&TIME_ZONE, "timezone", "Europe/London", "Set race timezone.")
	flag.DurationVar(&RACE_TIMEOUT_DURATION, "timeout", 2*time.Minute, "Set race timeout duration. After this time if nobody passes the finish line the race will be stopped. Valid time units are: 's' (second), 'm' (minute), 'h' (hour).")
	flag.DurationVar(&MINIMAL_LAP_TIME_DURATION, "laptime", 45*time.Second, "Minimal lap time duration. Results smaller than this duration would be considered wrong." )
	flag.BoolVar(&AVERAGE_RESULTS, "average", true, "Calculate average results instead of only first results.")

	//db
	flag.StringVar(&DB_FILE_PATH, "dbpath", ".", "Provide path to writable directory to store database data.")
	flag.StringVar(&DB_TYPE, "dbtype", "sqlite", "Select db type: sqlite / genji / postgres")

	//PostgreSQL related start
	flag.StringVar(&PG_HOST, "pghost", "127.0.0.1", "PostgreSQL DB host.")
	flag.IntVar(&PG_PORT, "pgport", 5432, "PostgreSQL DB port.")
	flag.StringVar(&PG_USER, "pguser", "postgres", "PostgreSQL DB user.")
	flag.StringVar(&PG_PASS, "pgpass", "", "PostgreSQL DB password.")
	flag.StringVar(&PG_DB_NAME, "pgdbname", "chicha", "PostgreSQL DB name.")
	flag.StringVar(&PG_SSL, "pgssl", "prefer", "disable / allow / prefer / require / verify-ca / verify-full - PostgreSQL ssl modes: https://www.postgresql.org/docs/current/libpq-ssl.html")

	//process all flags
	flag.Parse()


	if *flagVersion  {
		if VERSION != "" {
			fmt.Println("Version:", VERSION)
		}
		os.Exit(0)
	}

}
