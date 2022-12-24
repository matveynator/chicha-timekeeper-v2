package Collector

import (
	"bytes"
	"io"
	"log"
	"net"
	"time"
	"fmt"
	"strconv"
	"strings"
	"encoding/csv"
	"encoding/xml"

	"chicha/packages/config"
	"chicha/packages/data"
	//"chicha/packages/database"
	"chicha/packages/proxy"
)

func IsValidXML(data []byte) bool {
	return xml.Unmarshal(data, new(interface{})) == nil
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

			if ne, ok := err.(*net.OpError); ok && ne.Temporary() {
				if tempDelay == 0 {
					tempDelay = 5 * time.Millisecond
				} else {
					tempDelay *= 2
				}
				if max := 1 * time.Second; tempDelay > max {
					tempDelay = max
				}
				log.Printf("http: Accept error: %v; retrying in %v", err, tempDelay)
				time.Sleep(tempDelay)
				continue
			}

			break
		}
		tempDelay = 0

		data := buf[:size]

		rawData.IP = fmt.Sprintf("%s", connection.RemoteAddr().(*net.TCPAddr).IP)

		//various data formats processing (text csv, xml) start:
		if !IsValidXML(data) {
			// CSV data processing
			//log.Println("Received data is not XML, trying CSV text...")
			//received data of type TEXT (parse TEXT).
			r := csv.NewReader(bytes.NewReader(data))
			r.Comma = ','
			//r.FieldsPerRecord = 3
			//make variable number of fields "-1"
			r.FieldsPerRecord = -1
			CSV, err := r.Read()
			if err != nil {
				log.Println("Recived incorrect CSV data", err)
				continue
			}

			// Prepare antenna position
			antennaPosition, err := strconv.ParseInt(strings.TrimSpace(CSV[2]), 10, 64)
			if err != nil {
				log.Println("Recived incorrect Antenna position CSV value:", err)
				continue
			}
			rawData.DiscoveryUnixTime, err = strconv.ParseInt(strings.TrimSpace(CSV[1]), 10, 64)
			if err != nil {
				log.Println("Recived incorrect discovery unix time CSV value:", err)
				continue
			}
			rawData.TagID = strings.TrimSpace(CSV[0])
			rawData.Antenna = uint8(antennaPosition)

			// XML data processing
		} else {
			// XML data processing
			// Prepare date
			//log.Println("Received data is valid XML")
			err := xml.Unmarshal(data, &rawData)
			if err != nil {
				log.Println("xml.Unmarshal ERROR:", err)
				continue
			}
			//log.Println("TIME_ZONE=", Config.TIME_ZONE)
			loc, err := time.LoadLocation(Config.TIME_ZONE)
			if err != nil {
				log.Println(err)
				continue
			}
			xmlTimeFormat := `2006/01/02 15:04:05.000`
			discoveryTime, err := time.ParseInLocation(xmlTimeFormat, string(rawData.DiscoveryUnixTime), loc)
			if err != nil {
				log.Println("time.ParseInLocation ERROR:", err)
				continue
			}
			rawData.DiscoveryUnixTime = discoveryTime.UnixNano()/int64(time.Millisecond)
			// Additional preparing for TagID
			rawData.TagID = strings.ReplaceAll(rawData.TagID, " ", "")

			// Prepare antenna position
			rawData.Antenna = uint8(rawData.Antenna)
		}
		//various data formats processing (text csv, xml) end.

		//set microsecond resolution for logging:
		log.SetFlags(log.LstdFlags | log.Lmicroseconds)
		//Debug all received data from RFID reader
		log.Printf("TAG=%s, TIME=%d, IP=%s, ANT=%d\n", rawData.TagID, rawData.DiscoveryUnixTime, rawData.IP, rawData.Antenna)

		if Config.PROXY_ADDRESS != "" {
			go Proxy.CreateNewProxyTask(rawData)
		}

/*
		if len(laps) == 0 {
			//laps buffer empty - recreate last race from db:
			log.Println("laps buffer empty - recreate last race from db")
			laps, err = Database.GetCurrentRaceDataFromDB()
			if err == nil {
				log.Printf("laps buffer recreated with %d records from db.\n", len(laps))
				go Database.addNewLapToLapsBuffer(rawData)
			} else {
				log.Println("laps buffer recreation failed with:", err)
				go Database.addNewLapToLapsBuffer(rawData)
			}
		}

		if len(laps) > 0 {
			// Add current Lap to Laps buffer
			go Database.addNewLapToLapsBuffer(rawData)
		}

*/
	}
}

// Start data collector from RFID readers.
func StartDataCollector() {

	//unlock buffer operations
	Config.ChannelBufferLocker <- 0 //Put the initial value into the channel to unlock operations

	//unlock db operations
	Config.ChannelDBLocker <- 0 //Put the initial value into the channel to unlock operations

	//spin forever go routine to save in db with some interval:
	//go saveLapsBufferToDB()

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

