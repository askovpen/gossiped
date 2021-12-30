package utils

import (
	"os"
)

// FileExists Check file exists
func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func CreateFlagFile(filename string) error {
	if filename != "" {
		_, err := os.Create(filename)
		return err
	}
	return nil
}
