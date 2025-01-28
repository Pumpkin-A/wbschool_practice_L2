package eventstorage

import (
	"calendar-server/models"
	"fmt"
	"log"
	"sync"
	"time"
)

type DB interface {
	GetEvents() ([]models.EventData, error)
	SaveEvents([]models.EventData) error
}

type EventStorage struct {
	db     DB
	events []models.EventData
	lastID int
	rwm    sync.RWMutex
}

func New(db DB) (*EventStorage, error) {
	oldEvents, err := db.GetEvents()
	if err != nil {
		log.Fatalf("File db error open %v", err)
		return nil, fmt.Errorf("New: %w", err)
	}

	var lastID int
	for i := 0; i < len(oldEvents); i++ {
		if lastID < oldEvents[i].ID {
			lastID = oldEvents[i].ID
		}
	}

	return &EventStorage{
		db:     db,
		events: oldEvents,
		lastID: lastID,
		rwm:    sync.RWMutex{},
	}, nil
}

func (es *EventStorage) Close() error {
	return es.db.SaveEvents(es.events)
}

func (es *EventStorage) getNewID() int {
	defer func() {
		es.lastID++
	}()
	return es.lastID
}

func (es *EventStorage) GetEvent(ID int) (models.EventData, error) {
	es.rwm.RLock()
	defer es.rwm.RUnlock()

	index, err := es.findIndexByID(ID)
	if err != nil {
		return models.EventData{}, fmt.Errorf("GetEvent: %w", err)
	}

	return es.events[index], nil
}

func (es *EventStorage) AddEvent(data models.NewEventData) (int, error) {
	es.rwm.Lock()

	newID := es.getNewID()
	event := models.EventData{
		ID:     newID,
		UserID: data.UserID,
		Name:   data.Name,
		Date:   data.Date,
	}
	es.addEvent(event)

	es.rwm.Unlock()

	return newID, nil
}

func (es *EventStorage) addEvent(event models.EventData) {
	es.events = append(es.events, event)
}

func (es *EventStorage) UpdateEvent(data models.UpdateEventData) (models.EventData, error) {
	es.rwm.Lock()
	defer es.rwm.Unlock()

	index, err := es.findIndexByID(data.ID)
	if err != nil {
		return models.EventData{}, fmt.Errorf("UpdateEvent: %w", err)
	}

	if data.Name != nil {
		es.events[index].Name = *data.Name
	}
	if data.UserID != nil {
		es.events[index].UserID = *data.UserID
	}
	if data.Date != nil {
		es.events[index].Date = *data.Date
	}

	return es.events[index], nil
}

func (es *EventStorage) DeleteEvent(ID int) (models.EventData, error) {
	es.rwm.Lock()
	defer es.rwm.Unlock()

	index, err := es.findIndexByID(ID)
	if err != nil {
		return models.EventData{}, fmt.Errorf("DeleteEvent: %w", err)
	}

	deleted, err := es.deleteEventByIndex(index)
	if err != nil {
		return models.EventData{}, fmt.Errorf("DeleteEvent: %w", err)
	}

	return deleted, nil
}

func (es *EventStorage) deleteEventByIndex(index int) (models.EventData, error) {
	if index < 0 || index >= len(es.events) {
		return models.EventData{}, fmt.Errorf("deleteEventByIndex: incorrect index of event: %d, events count: %d", index, len(es.events))
	}
	deletedEvent := es.events[index]
	if index != len(es.events)-1 {
		es.events[index] = es.events[len(es.events)-1]
	}
	es.events = es.events[:len(es.events)-1]

	return deletedEvent, nil
}

func (es *EventStorage) findIndexByID(ID int) (int, error) {
	for i := 0; i < len(es.events); i++ {
		if es.events[i].ID == ID {
			return i, nil
		}
	}
	return 0, fmt.Errorf("no event with index: %d", ID)
}

func (es *EventStorage) FindByDay(day time.Time) ([]models.EventData, error) {
	es.rwm.RLock()
	defer es.rwm.RUnlock()

	var found []models.EventData
	for i := 0; i < len(es.events); i++ {
		curDay, _ := time.Parse("2006-01-02", es.events[i].Date)
		if day.Year() == curDay.Year() && day.YearDay() == curDay.YearDay() {
			found = append(found, es.events[i])
		}
	}

	return found, nil
}

func (es *EventStorage) FindByWeek(week time.Time) ([]models.EventData, error) {
	es.rwm.RLock()
	defer es.rwm.RUnlock()

	var found []models.EventData
	yearToFind, weekToFind := week.ISOWeek()
	for i := 0; i < len(es.events); i++ {
		curDay, _ := time.Parse("2006-01-02", es.events[i].Date)
		yearCur, weekCur := curDay.ISOWeek()
		if yearToFind == yearCur && weekToFind == weekCur {
			found = append(found, es.events[i])
		}
	}

	return found, nil
}

func (es *EventStorage) FindByMonth(month time.Time) ([]models.EventData, error) {
	es.rwm.RLock()
	defer es.rwm.RUnlock()

	var found []models.EventData
	yearToFind, monthToFind := month.Year(), month.Month()
	for i := 0; i < len(es.events); i++ {
		curData, _ := time.Parse("2006-01-02", es.events[i].Date)
		yearCur, monthCur := curData.Year(), curData.Month()
		if yearToFind == yearCur && monthToFind == monthCur {
			found = append(found, es.events[i])
		}
	}

	return found, nil
}

func (es *EventStorage) FindByYear(year int) ([]models.EventData, error) {
	es.rwm.RLock()
	defer es.rwm.RUnlock()

	var found []models.EventData
	for i := 0; i < len(es.events); i++ {
		curDay, err := time.Parse("2006-01-02", es.events[i].Date)
		if err != nil {
			return nil, fmt.Errorf("cant parse internal data date")
		}
		if year == curDay.Year() {
			found = append(found, es.events[i])
		}
	}

	return found, nil
}
