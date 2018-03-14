package utils

import (
	"os"
	"strings"
)

const (
	BakExt = ".bak"
)

func BackupFile(path string) error {
	if !IsExist(path) {
		return os.ErrNotExist
	}

	return os.Rename(path, path+BakExt)
}

func RollBack(path string) error {
	backFilename := path + BakExt
	if !IsExist(backFilename) {
		return os.Remove(path)
	}

	return os.Rename(backFilename, path)
}

func IsExist(path string) bool {
	exist := true
	if _, err := os.Stat(path); os.IsNotExist(err) {
		exist = false
	}

	return exist
}

func GetExtention(filename string) string {
	splitByDot := strings.Split(filename, ".")
	return splitByDot[len(splitByDot)-1]
}

func RemoveFileIfExist(path string) error {
	if !IsExist(path) {
		return nil
	}
	return os.Remove(path)
}
