package Collector

import (
	"testing"
	"chicha/packages/data"
)

func TestParseCSVLine(t *testing.T) {
	var dataSlice [][]byte
	dataSlice = append(dataSlice, []byte("100008020200000100000189, 1622570553397, 3"))
	dataSlice = append(dataSlice, []byte("100008020200000100000189, 1622570553397, 3, 9.9.9.9"))

	var rawData Data.RawData
	var err error

	for data := range []dataSlice {
		rawData, err = parseCSVLine(data, "8.8.8.8")
		if err != nil {
			t.Errorln(err)
		} else {
			t.Logf("TAG=%s, TIME=%d, Reader-IP=%s, Reader-Antenna=%d, Proxy-IP=%s\n", rawData.TagID, rawData.DiscoveryUnixTime, rawData.ReaderIP, rawData.Antenna, rawData.ProxyIP)
		}
	}
}
