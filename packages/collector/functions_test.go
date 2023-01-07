package Collector

import (
	"testing"
	"chicha/packages/data"
)

func TestParseCSVLine(t *testing.T) {

	testData := []struct {
		csvInput []byte
		csvOutput Data.RawData
	}{
		{
			csvInput: []byte("100008020200000100000189, 1622570553397, 3"),
			csvOutput: Data.RawData{ "100008020200000100000189", 1622570553397, "8.8.8.8", 3, "" },
		},
		{
			csvInput: []byte("TONIBOU-001, 1622570553397, 3, 9.9.9.9"),
			csvOutput: Data.RawData{"TONIBOU-001", 1622570553397, "9.9.9.9", 3, "8.8.8.8" },
		},
	}

	for _,testCase := range testData {
		//check1:
		testResult, _ := parseCSVLine(testCase.csvInput, "8.8.8.8")
		if testResult != testCase.csvOutput  {
			t.Errorf("Incorrect result. Expect %v, got %v", testCase.csvOutput, testResult)
		}
	}
}
