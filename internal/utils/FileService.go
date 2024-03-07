package utils

import (
	"os"
)

func ReadFile(path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return nil, err
	}

	fileSize := fileInfo.Size()
	bytes := make([]byte, fileSize)

	_, err = file.Read(bytes)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}
