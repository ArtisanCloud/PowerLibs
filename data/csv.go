package data

import (
	"bytes"
	"encoding/csv"
	"os"
)

func CSVEncode(arrayCSV [][]string) ([]byte, error) {

	buffer := new(bytes.Buffer)
	writer := csv.NewWriter(buffer)

	err := writer.WriteAll(arrayCSV)

	return buffer.Bytes(), err

}

func CSVEncodeToFile(arrayCSV [][]string, file *os.File) (err error) {

	writer := csv.NewWriter(file)

	err = writer.WriteAll(arrayCSV)

	return err

}
