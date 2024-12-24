package mymain

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// TODO rename file to module entity

type Event struct {
	ID     int    `json:"id"`
	UserID int    `json:"user_id"`
	Name   string `json:"name"`
	Date   string `json:"date"`
}

type eventcontextKey int

const (
	name eventcontextKey = iota
	userID
	date
)

type EventStorage struct {
	data   []Event
	lastID int
	rwm    sync.RWMutex
}

func (es *EventStorage) getNewID() int {
	es.rwm.RLock()
	defer func() {
		es.lastID++
		es.rwm.Unlock()
	}()
	return es.lastID
}

func (es *EventStorage) initEventStorage(oldEvents []Event) {
	for i := 0; i < len(oldEvents); i++ {
		if es.lastID < oldEvents[i].ID {
			es.lastID = oldEvents[i].ID
		}
	}
	es.data = oldEvents
	es.lastID = 1
	es.rwm = sync.RWMutex{}
}

// func (es *EventStorage) getEvent(ID int) (Event, error) {
// 	es.rwm.RLock()
// 	defer es.rwm.RUnlock()

// 	index, err := es.findIndex(ID)
// 	if err != nil {
// 		return Event{}, nil
// 	}
// 	return es.data[index], nil
// }

func (es *EventStorage) addEvent(data DataToAddNewEvent) (int, error) {
	newID := es.getNewID()

	es.rwm.Lock()
	es.data = append(es.data, Event{ID: newID, UserID: data.UserID, Name: data.Name, Date: data.Date})
	es.rwm.Unlock()

	return newID, nil
}

func (es *EventStorage) updateEvent(ID int, ctx context.Context) (Event, error) {
	index, err := es.findIndex(ID)
	if err != nil {
		return Event{}, err
	}

	es.rwm.RLock()
	if newName, ok := ctx.Value(name).(string); ok {
		es.data[index].Name = newName
	}
	if newUserID, ok := ctx.Value(userID).(int); ok {
		es.data[index].UserID = newUserID
	}
	if newDate, ok := ctx.Value(date).(string); ok {
		es.data[index].Date = newDate
	}
	es.rwm.RUnlock()

	return es.data[index], nil
}

func (es *EventStorage) DeleteEvent(ID int) (Event, error) {
	index, err := es.findIndex(ID)
	if err != nil {
		return Event{}, err
	}

	es.rwm.Lock()
	deletedEvent := es.data[index]
	if index != len(es.data)-1 {
		es.data[index] = es.data[len(es.data)-1]
	}
	es.data = es.data[:len(es.data)-1]
	es.rwm.Unlock()

	return deletedEvent, nil
}

func (es *EventStorage) findIndex(ID int) (int, error) {
	es.rwm.RLock()
	defer es.rwm.RUnlock()
	for i := 0; i < len(es.data); i++ {
		if es.data[i].ID == ID {
			return i, nil
		}
	}
	return 0, fmt.Errorf("no event with sucn index")
}

func (es *EventStorage) findByDay(day time.Time) ([]Event, error) {
	es.rwm.RLock()
	defer es.rwm.RUnlock()

	var found []Event
	for i := 0; i < len(es.data); i++ {
		curDay, _ := time.Parse("2006-01-02", es.data[i].Date)
		if day.Year() == curDay.Year() && day.YearDay() == curDay.YearDay() {
			found = append(found, es.data[i])
		}
	}
	return found, nil
}

func (es *EventStorage) findByWeek(week time.Time) ([]Event, error) {
	es.rwm.RLock()
	defer es.rwm.RUnlock()

	var found []Event
	yearToFind, weekToFind := week.ISOWeek()
	for i := 0; i < len(es.data); i++ {
		curDay, _ := time.Parse("2006-01-02", es.data[i].Date)
		yearCur, weekCur := curDay.ISOWeek()
		if yearToFind == yearCur && weekToFind == weekCur {
			found = append(found, es.data[i])
		}
	}
	return found, nil
}

func (es *EventStorage) findByMonth(month time.Time) ([]Event, error) {
	es.rwm.RLock()
	defer es.rwm.RUnlock()

	var found []Event
	yearToFind, monthToFind := month.Year(), month.Month()
	for i := 0; i < len(es.data); i++ {
		curData, _ := time.Parse("2006-01-02", es.data[i].Date)
		yearCur, monthCur := curData.Year(), curData.Month()
		if yearToFind == yearCur && monthToFind == monthCur {
			found = append(found, es.data[i])
		}
	}
	return found, nil
}

func (es *EventStorage) findByYear(year int) ([]Event, error) {
	es.rwm.RLock()
	defer es.rwm.RUnlock()

	var found []Event
	for i := 0; i < len(es.data); i++ {
		curDay, _ := time.Parse("2006-01-02", es.data[i].Date)
		if year == curDay.Year() {
			found = append(found, es.data[i])
		}
	}
	return found, nil
}
