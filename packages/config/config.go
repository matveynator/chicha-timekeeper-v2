package Config

import (
	"log"
	"os"
	"flag"
	"time"
)

var VERSION, RFID_LISTENER_ADDRESS, WEB_LISTENER_ADDRESS, REDIRECT_ADDRESS, DB_FILE_PATH, TIME_ZONE string
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
	RFID_LISTENER_ADDRESS = *flag.String("rfid", "0.0.0.0:4000", "Provide IP address and port to listen for RFID data.")
	WEB_LISTENER_ADDRESS = *flag.String("web", "0.0.0.0:80", "Provide IP address and port to listen for HTTP connections.")
	REDIRECT_ADDRESS = *flag.String("redirect", "", "Redirect rfid data stream to another instance. For example: -redirect '10.9.8.7:4000'.")
	DB_FILE_PATH = *flag.String("db", "chicha.db.sqlite", "Provide /path/to/chicha.db.sqlite database.") 
	TIME_ZONE = *flag.String("timezone", "Europe/London", "Set race timezone.")
	RACE_TIMEOUT_DURATION = *flag.Duration("timeout", 2*time.Minute, "Set race timeout duration. After this time if nobody passes the finish line the race will be stopped. Valid time units are: 's' (second), 'm' (minute), 'h' (hour).")
	MINIMAL_LAP_TIME_DURATION = *flag.Duration("laptime", 45*time.Second, "Minimal lap time duration. Results smaller than this duration would be considered wrong." )
	AVERAGE_RESULTS = *flag.Bool("average", true, "Calculate average results instead of only first results.")

	//process all flags
	flag.Parse()


	if *flagVersion  {
		if VERSION != "" {
			log.Println("Version:", VERSION)
		}
		os.Exit(0)
	}
	//startup:
	log.Println("Welcome to CHICHA, the competition timekeeper (chronograph).")
	log.Println("github.com/matveynator/chicha")
	if VERSION != "" {
		log.Println("Version:", VERSION)
	}
}
