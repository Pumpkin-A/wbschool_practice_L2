package filedb

import (
	"calendar-server/models"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
)

type FileDB struct {
	filename string
	file     *os.File
}

func New(filename string) (*FileDB, error) {
	file, err := initFile(filename)
	if err != nil {
		return nil, fmt.Errorf("NewFileDB: %w", err)
	}

	return &FileDB{
		filename: filename,
		file:     file,
	}, nil
}

func initFile(filename string) (*os.File, error) {
	ptr, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		return nil, fmt.Errorf("initFile: %w", err)
	}

	return ptr, nil
}

func (fdb *FileDB) GetEvents() ([]models.EventData, error) {
	return getEvents(fdb.file)
}

func getEvents(file *os.File) ([]models.EventData, error) {
	decoder := json.NewDecoder(file)
	if _, err := decoder.Token(); err != nil {
		if errors.Is(err, io.EOF) {
			log.Printf("file %s is empty", file.Name())
			return []models.EventData{}, nil
		}
		return []models.EventData{}, fmt.Errorf("getEvents decoder: %w", err)
	}

	var events []models.EventData
	for decoder.More() {
		var e models.EventData
		if err := decoder.Decode(&e); err != nil {
			return []models.EventData{}, fmt.Errorf("getEvents unmarhsal: %w", err)
		}
		events = append(events, e)
	}

	fmt.Printf("events: %+v\n", events)

	return events, nil
}

func (db *FileDB) SaveEvents(data []models.EventData) error {
	if err := db.file.Truncate(0); err != nil {
		return err
	}
	if _, err := db.file.Seek(0, 0); err != nil {
		return err
	}

	encoder := json.NewEncoder(db.file)
	fmt.Println("Write to file db!!", db.file)

	return encoder.Encode(data)
}

func (db *FileDB) Close() error {
	return db.file.Close()
}
