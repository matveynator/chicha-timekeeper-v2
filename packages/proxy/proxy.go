package Proxy

import (
	"log"
	"fmt"
	"net"
	"time"
	"github.com/eternnoir/gncp"

	"chicha/packages/data"
	"chicha/packages/config"
)

var ProxyConnectionPool gncp.ConnPool

// connectionPoolCreator let connection know how to create new connection.
func connectionPoolCreator() (net.Conn, error) {
  return net.Dial("tcp", Config.PROXY_ADDRESS)
}

func restartConnectionPool() {
   var err error
	//create connection pool to Config.PROXY_ADDRESS with minimum 3 and max 10 connections.
	ProxyConnectionPool, err = gncp.NewPool(3, 10, connectionPoolCreator)
	if err != nil {
		log.Println("Proxy connection pool restart error:", err)
		return
	}
}

func init() {
	if Config.PROXY_ADDRESS != "" {
		log.Println("Started tcp proxy restream to:", Config.PROXY_ADDRESS)
		restartConnectionPool()
	}
}

func ProxyDataToAnotherHost(averageResult Data.AverageResult) {
	//get one connection from pool

	//connection, err := ProxyConnectionPool.Get()
	connection, err := ProxyConnectionPool.GetWithTimeout(time.Duration(1) * time.Second)
	if err != nil {
		log.Println("Connection pool error:", err)
		restartConnectionPool()
		return
	}
	defer connection.Close()
	fmt.Fprintf(connection, "%s, %d, %d, %s\n", averageResult.TagID, averageResult.DiscoveryUnixTime, averageResult.Antenna, averageResult.IP)
}
