package Proxy

import (
	"log"
	"fmt"
	"net"

	"chicha/packages/config"
)

func init() {
	if Config.PROXY_ADDRESS != "" {
		log.Println("Started tcp proxy restream to: ", Config.PROXY_ADDRESS)
	}
}

func ProxyDataToAnotherHost(tagID string, unixTime int64, antennaNumber uint8, ip string) {
	connection, err := net.Dial("tcp", Config.PROXY_ADDRESS)
	if err != nil {
		log.Println("dial error:", err)
		return
	}
	defer connection.Close()
	fmt.Fprintf(connection, "%s, %d, %d, %s\n", tagID, unixTime, antennaNumber, ip)
}
