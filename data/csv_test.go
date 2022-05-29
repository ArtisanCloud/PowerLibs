package data

import (
	"os"
	"testing"
)

func Test_CSVEncodeToFile(t *testing.T) {
	var csvStruct [][]string
	csvStruct = [][]string{
		{"name", "address", "phone"},
		{"Ram", "Tokyo", "1236524"},
		{"Shaym", "Beijing", "8575675.484"},
	}

	csvFile, err := os.Create("csv_test.csv")
	if err != nil {
		t.Error(err)
	}

	defer csvFile.Close()

	err = CSVEncodeToFile(csvStruct, csvFile)
	if err != nil {
		t.Error(err)
	}

}
