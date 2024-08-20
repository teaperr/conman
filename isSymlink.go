package main

import (
	"os"
)

func isSymlink(file string) (bool, error) {
	fileInfo, err := os.Lstat(file)
	if err != nil {
		return false, err
	}
	return fileInfo.Mode()&os.ModeSymlink != 0, nil
}
