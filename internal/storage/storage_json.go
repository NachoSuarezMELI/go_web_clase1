package storage

import (
	"encoding/json"
	"io"
	"os"
	"web/clase1/platform/tools"
)

type StorageJSON struct {
	FileName string
}

func NewStorageJSON(fileName string) *StorageJSON {
	return &StorageJSON{
		FileName: fileName,
	}
}

func (s *StorageJSON) Read() ([]byte, error) {
	return tools.ReadFile(s.FileName)
}

func (s *StorageJSON) Write(data []byte) error {
	// Check if data has JSON format
	var jsonData interface{}
	err := json.Unmarshal(data, &jsonData)
	if err != nil {
		return err
	}

	// Write data to file
	file, err := os.OpenFile(s.FileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	w := io.Writer(file)
	_, err = w.Write(data)
	if err != nil {
		return err
	}

	return nil
}
