package Collector

import (
	"io"
	"log"
	"net"
	"time"

	"chicha/packages/config"
	"chicha/packages/data"
	"chicha/packages/database"
	"chicha/packages/proxy"
	"chicha/packages/timekeeper"
)
//initial function
func init() {
	//set microsecond resolution for logging:
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
}

func processConnection(connection net.Conn) {
	defer connection.Close()
	var tempDelay time.Duration // how long to sleep on accept failure
	var rawData Data.RawData

	// Read connection in lap
	for {
		buf := make([]byte, 1024)
		size, err := connection.Read(buf)
		if err != nil {
			if err == io.EOF {
				//log.Println("conn.Read(buf) error:", err)
				//log.Println("Message EOF detected - closing LAN connection.")
				break
			}

			networkError, ok := err.(*net.OpError) 
			if ok && networkError.Temporary() {
				if tempDelay == 0 {
					tempDelay = 5 * time.Millisecond
				} else {
					tempDelay = 2 * tempDelay
				}
				if max := 1 * time.Second; tempDelay > max {
					tempDelay = max
				}
				log.Printf("Accept error: %v; retrying in %v", err, tempDelay)
				time.Sleep(tempDelay)
				continue
			}

			break
		}
		tempDelay = 0

		data := buf[:size]

		remoteIPAddress := connection.RemoteAddr().(*net.TCPAddr).IP.String()

		//various data formats processing (text csv, xml) start:
		if !IsValidXML(data) {
			// CSV data processing
			rawData, err = parseCSVLine(data, remoteIPAddress)
			if err != nil {
				log.Println(err)
			}
		} else {
			// XML data processing
			rawData, err = parseXMLPacket(data, remoteIPAddress)
			if err != nil {
				log.Println(err)
			}
		}

		//create a proxy task if needed (via Proxy.ProxyTask channel):
		if Config.PROXY_ADDRESS != "" {
			//send rawData to Proxy.ProxyTask channel
			Proxy.ProxyTask <- rawData
		}

		//create timekeeper task:
		Timekeeper.TimekeeperTask <- rawData

		//create a database task:
		Database.DatabaseTask <- rawData

	}
}

// Start data collector from RFID readers.
func StartDataCollector() {
	//parse configuration:
	Config.ParseFlags()

	// Start listener
	collector, err := net.Listen("tcp", Config.COLLECTOR_LISTENER_ADDRESS)
	if err != nil {
		log.Panicln("Error: collector can't start. ", err)
	}
	defer collector.Close()


	// Listen new connections
	for {
		connection, err := collector.Accept()
		if err != nil {
			if err != io.EOF {
				log.Panicln("tcp connection error:", err)
			}
		}

		go processConnection(connection)
	}
}

