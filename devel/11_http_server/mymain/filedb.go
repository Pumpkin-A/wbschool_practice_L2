package mymain

import (
	"encoding/json"
	"fmt"
	"os"
)

// Concrete file DB realiztion

type fileDB struct {
	filename string
	filePtr  *os.File
}

func newFileDB(filename string) *fileDB {
	return &fileDB{filename: filename, filePtr: nil}
}

func (db *fileDB) openConnection() ([]Event, error) {
	ptr, err := os.OpenFile(db.filename, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		return []Event{}, err
	}

	decoder := json.NewDecoder(ptr)
	if _, err := decoder.Token(); err != nil {
		return []Event{}, nil
	}
	var events []Event

	for decoder.More() {
		var e Event
		if err := decoder.Decode(&e); err != nil {
			return []Event{}, err
		}
		events = append(events, e)
	}

	db.filePtr = ptr
	return events, nil
}

func (db *fileDB) closeConnction(data []Event) error {
	defer db.filePtr.Close()
	if err := db.filePtr.Truncate(0); err != nil {
		return err
	}
	if _, err := db.filePtr.Seek(0, 0); err != nil {
		return err
	}

	encoder := json.NewEncoder(db.filePtr)
	fmt.Println("Write to file db!!", db.filePtr)
	return encoder.Encode(data)
}
