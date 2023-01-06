package Collector

import (
	"testing"
	"chicha/packages/data"
)

func TestParseCSVLine(t *testing.T) {

	//prepare test data:
	testInput1 := []byte("100008020200000100000189, 1622570553397, 3")
	testInput2 := []byte("100008020200000100000189, 1622570553397, 3, 9.9.9.9")

	var testOutput1 Data.RawData
	testOutput1.TagID="100008020200000100000189"
	testOutput1.DiscoveryUnixTime=1622570553397
	testOutput1.ReaderIP = "8.8.8.8"
	testOutput1.Antenna = 3

	var testOutput2 Data.RawData
	testOutput2.TagID="100008020200000100000189"
	testOutput2.DiscoveryUnixTime=1622570553397
	testOutput2.ReaderIP = "9.9.9.9"
	testOutput2.Antenna = 3
  testOutput2.ProxyIP = "8.8.8.8"


	//check1:
	result1, _ := parseCSVLine(testInput1, "8.8.8.8")
	if testOutput1 != result1 {
		t.Errorf("Incorrect result. Expect %v, got %v", testOutput1, result1)
	}

	//check2:
	result2, _ := parseCSVLine(testInput2, "8.8.8.8")
	if testOutput2 != result2 {
		t.Errorf("Incorrect result. Expect %v, got %v", testOutput2, result2)
	}
}
