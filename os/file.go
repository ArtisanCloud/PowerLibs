package os

import (
	"io"
	"os"
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
