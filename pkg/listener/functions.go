package listener
import (
	"bytes"
	"fmt"
	"log"
	"net"
	"time"
	"strconv"
	"strings"
	"errors"
	"encoding/csv"
	"encoding/xml"

	"chicha/pkg/config"
	"chicha/pkg/data"
)
func IsValidXML(byteData []byte) bool {
	return xml.Unmarshal(byteData, new(interface{})) == nil
}

func parseCSVLine(byteData []byte, remoteIPAddress string) (rawData data.RawData, err error){
	if !IsValidXML(byteData) {
		// CSV data processing
		var csvField []string
		csvField, err = csv.NewReader(bytes.NewReader(byteData)).Read()
		if err != nil {
			log.Println(err)
			return
		} else {
			//if one packet received (valid):
			numberOfCSVColumns := len(csvField)
			if numberOfCSVColumns == 3 || numberOfCSVColumns == 4 {
				// Prepare antenna position
				var antennaPosition int64
				antennaPosition, err = strconv.ParseInt(string(strings.TrimSpace(csvField[2])), 10, 64)
				if err != nil {
					log.Println("Recived incorrect Antenna position CSV value:", err)
					return
				} else { 
					rawData.DiscoveryUnixTime, err = strconv.ParseInt(string(strings.TrimSpace(csvField[1])), 10, 64)
					if err != nil {
						log.Println("Recived incorrect discovery unix time CSV value:", err)
						return
					} else {
						rawData.TagId = string(strings.TrimSpace(csvField[0]))

						rawData.Antenna = uint(antennaPosition)

						if numberOfCSVColumns == 3 {
							//reader connection without proxy
							rawData.ReaderIP = remoteIPAddress

							//Debug all received data from RFId reader
							log.Printf("TAG=%s, TIME=%d, Reader-IP=%s, Reader-ANT=%d\n", rawData.TagId, rawData.DiscoveryUnixTime, rawData.ReaderIP, rawData.Antenna)

						} else if numberOfCSVColumns == 4 {
							//proxy connection
							if net.ParseIP(string(strings.TrimSpace(csvField[3]))) != nil {
								//sending data with proxy
								rawData.ReaderIP = string(strings.TrimSpace(csvField[3]))
								rawData.ProxyIP = remoteIPAddress
							} else {
								//sending csvField[3] is not an IP address
								rawData.ReaderIP = remoteIPAddress
							}
							//Debug all received data from PROXY
							log.Printf("TAG=%s, TIME=%d, Reader-IP=%s, Reader-Antenna=%d, Proxy-IP=%s\n", rawData.TagId, rawData.DiscoveryUnixTime, rawData.ReaderIP, rawData.Antenna, rawData.ProxyIP)

						}
					}
				}
			} else {
				err = errors.New(fmt.Sprintf("Error: Parser expected 3 or 4 fields in CSV, but received %d fields.", numberOfCSVColumns))
			}
		}
	} else {
		err = errors.New("Error: received data is not a valid CSV - it is XML!")
	}
	return
}

func parseXMLPacket(byteData []byte, remoteIPAddress string, config Config.Settings)(rawData data.RawData, err error) {
	if IsValidXML(byteData) {
		// XML data processing
		var rawXMLData data.RawXMLData
		err = xml.Unmarshal(byteData, &rawXMLData)
		if err != nil {
			log.Println("xml.Unmarshal ERROR:", err)
			return
		}

		var loc *time.Location
		loc, err = time.LoadLocation(config.TIME_ZONE)
		if err != nil {
			log.Println(err)
			return
		}
		xmlTimeFormat := `2006/01/02 15:04:05.000`
		var discoveryTime time.Time
		discoveryTime, err = time.ParseInLocation(xmlTimeFormat, rawXMLData.DiscoveryTime, loc)
		if err != nil {
			log.Println("time.ParseInLocation ERROR:", err)
			return
		}
		rawData.DiscoveryUnixTime = discoveryTime.UnixNano()/int64(time.Millisecond)
		rawData.TagId = strings.ReplaceAll(rawXMLData.TagId, " ", "")
		rawData.Antenna = uint(rawXMLData.Antenna)

		if net.ParseIP(rawXMLData.ReaderIP) != nil {
			//connection received from proxy (not from reader).
			rawData.ReaderIP = rawXMLData.ReaderIP
			rawData.ProxyIP = remoteIPAddress

			//Debug all received data from PROXY
			log.Printf("TAG=%s, TIME=%d, Reader-IP=%s, Reader-Antenna=%d, Proxy-IP=%s\n", rawData.TagId, rawData.DiscoveryUnixTime, rawData.ReaderIP, rawData.Antenna, rawData.ProxyIP)
		} else {
			//connection received from reader (not from proxy)
			rawData.ReaderIP = remoteIPAddress
			//Debug all received data from RFID reader
			log.Printf("TAG=%s, TIME=%d, Reader-IP=%s, Reader-ANT=%d\n", rawData.TagId, rawData.DiscoveryUnixTime, rawData.ReaderIP, rawData.Antenna)
		}
	} else {
		err = errors.New("Error: received data is not a valid XML.")
	}
	return
}
