package os

import (
	"bytes"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// https://www.golangprograms.com/files-directories-examples.html

func CopyFile(src string, dst string) (err error) {

	fin, err := os.Open(src)
	if err != nil {
		return err
	}
	defer fin.Close()

	fOut, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer fOut.Close()

	_, err = io.Copy(fOut, fin)

	if err != nil {
		return err
	}

	return err
}

func MoveFile(src string, dst string) (err error) {
	err = os.Rename(src, dst)
	return err
}

func ConvertFileHandleToReader(obj interface{}) (data io.Reader, err error) {
	switch obj.(type) {
	case string:
		return strings.NewReader(obj.(string)), nil

	case []byte:
		return bytes.NewReader(obj.([]byte)), nil

	case io.Reader:
		return obj.(io.Reader), nil

	default:
		return nil, errors.New("not support file handle data")
	}

}

func CreateDirectoriesForFiles(outputFile string) error {
	dir := filepath.Dir(outputFile)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	return nil
}
