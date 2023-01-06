package Config

import (
	"os"
	"fmt"
	"flag"
	"time"
	"hash/fnv"
)

var APP_NAME = "chicha"
var VERSION, COLLECTOR_LISTENER_ADDRESS, COLLECTOR_LISTENER_ADDRESS_HASH, WEB_LISTENER_ADDRESS, PROXY_ADDRESS, DB_TYPE, DB_FILE_PATH, DB_FULL_FILE_PATH, PG_HOST, PG_USER, PG_PASS, PG_DB_NAME, PG_SSL, TIME_ZONE string
var PG_PORT int
var AVERAGE_RESULTS bool
var RACE_TIMEOUT_DURATION, MINIMAL_LAP_TIME_DURATION time.Duration

//lockers
var ChannelBufferLocker = make(chan int, 1)
var ChannelDBLocker = make(chan int, 1)


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


func ParseFlags()  {

	flagVersion := flag.Bool("version", false, "Output version information")
	flag.StringVar(&COLLECTOR_LISTENER_ADDRESS, "collector", "0.0.0.0:4000", "Provide IP address and port to collect and parse data from RFID and timing readers.")
	flag.StringVar(&WEB_LISTENER_ADDRESS, "web", "0.0.0.0:80", "Provide IP address and port to listen for HTTP connections from clients.")
	flag.StringVar(&PROXY_ADDRESS, "proxy", "", "Proxy received data to another collector. For example: -proxy '10.9.8.7:4000'.")
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
	
	//делаем хеш от порта коллектора чтобы использовать в уникальном названии файла бд
	COLLECTOR_LISTENER_ADDRESS_HASH = hash(COLLECTOR_LISTENER_ADDRESS)

	//путь к файлу бд
	DB_FULL_FILE_PATH = fmt.Sprintf(DB_FILE_PATH+"/"+APP_NAME+"."+COLLECTOR_LISTENER_ADDRESS_HASH+".db."+DB_TYPE)

	if *flagVersion  {
		if VERSION != "" {
			fmt.Println("Version:", VERSION)
		}
		os.Exit(0)
	}

}
