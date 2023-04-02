<p align="left">
    <img property="og:image" src="https://repository-images.githubusercontent.com/577755312/57f67b11-437b-448f-b53e-cf47165612c2" width="25%">
</p>

# Chicha - the competition timekeeper (chronograph). Version 2.0

Free chronograf for runners, bycicles, motorcycles, carts, cars, trucks, atv and other types of race competitions. 
UHF-RFID compatible.

<p align="left">
    <img property="og:image" src="https://repository-images.githubusercontent.com/368199185/e26c553e-b23e-4bae-b4d2-c2df502e9f04" width="75%">
</p>


## [Demo: http://chicha.zabiyaka.net](http://chicha.zabiyaka.net/)


## Usage of chicha:
```
chicha -h
Usage of chicha:
  -average
    	Calculate average results instead of only first results. (default true)
  -average-duration duration
    	Duration to calculate average results. Results passed to reader during this duration will be calculated as average result. (default 1s)
  -collector string
    	Provide IP address and port to collect and parse data from RFID and timing readers. (default "0.0.0.0:4000")
  -db-path string
    	Provide path to writable directory to store database data. (default ".")
  -db-save-interval duration
    	Duration to save data from memory to database (disk). Setting duration too low may cause unpredictable performance results. (default 30s)
  -db-type string
    	Select db type: sqlite / genji / postgres (default "genji")
  -lap-time duration
    	Minimal lap time duration. Results smaller than this duration would be considered wrong. (default 45s)
  -pg-db-name string
    	PostgreSQL DB name. (default "chicha")
  -pg-host string
    	PostgreSQL DB host. (default "127.0.0.1")
  -pg-pass string
    	PostgreSQL DB password.
  -pg-port int
    	PostgreSQL DB port. (default 5432)
  -pg-ssl string
    	disable / allow / prefer / require / verify-ca / verify-full - PostgreSQL ssl modes: https://www.postgresql.org/docs/current/libpq-ssl.html (default "prefer")
  -pg-user string
    	PostgreSQL DB user. (default "postgres")
  -proxy string
    	Proxy incoming data to another chicha collector. For example: -proxy '10.9.8.7:4000'.
  -timeout duration
    	Set race timeout duration. After this time if nobody passes the finish line the race will be stopped. Valid time units are: 's' (second), 'm' (minute), 'h' (hour). (default 2m0s)
  -timezone string
    	Set race timezone. Example: Europe/Paris, Africa/Dakar, UTC, https://en.wikipedia.org/wiki/List_of_tz_database_time_zones (default "UTC")
  -version
    	Output version information
  -web string
    	Provide IP address and port to listen for HTTP connections from clients. (default "0.0.0.0:80")
```
