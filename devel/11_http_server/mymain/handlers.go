package mymain

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

func (app *Application) getEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.NotFound(w, r)
		return
	}

	fmt.Println("server:", r.URL.Query())
	eventIDs, ok := r.URL.Query()["id"]
	if !ok {
		sendResponse(w, http.StatusOK, app.events.data)
		return
	}
	eventID, err := strconv.Atoi(eventIDs[0])
	if err != nil {
		sendError(w, http.StatusBadRequest, nil)
		return
	}
	index, err := app.events.findIndex(eventID)
	if err != nil {
		sendError(w, http.StatusOK, nil)
		return
	}

	sendResponse(w, http.StatusOK, app.events.data[index])
}

type DataToAddNewEvent struct {
	UserID int    `json:"user_id"`
	Name   string `json:"name"`
	Date   string `json:"date"`
}

func (d DataToAddNewEvent) isValid() bool {
	if d.UserID <= 0 || d.Name == "" {
		return false
	}
	if _, err := time.Parse("2006-01-02", d.Date); err != nil {
		return false
	}
	return true
}

func (app *Application) AddEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.NotFound(w, r)
		return
	}

	var t DataToAddNewEvent
	err := json.NewDecoder(r.Body).Decode(&t)
	if err != nil {
		sendError(w, http.StatusBadRequest, nil)
		return
	}
	if !t.isValid() {
		sendError(w, http.StatusBadRequest, nil)
		return
	}

	ID, err := app.events.addEvent(t)
	if err != nil {
		sendError(w, http.StatusBadRequest, nil)
	}

	sendResponse(w, http.StatusOK, struct {
		ID int `json:"id"`
	}{ID: ID})
}

type DataToUpdateEvent struct {
	ID     int    `json:"id"`
	UserID int    `json:"user_id"`
	Name   string `json:"name"`
	Date   string `json:"date"`
}

func (d DataToUpdateEvent) isValid() bool {
	if d.ID <= 0 || d.UserID <= 0 || d.Name == "" {
		return false
	}
	if _, err := time.Parse("2006-01-02", d.Date); err != nil {
		return false
	}
	return true
}

func (app *Application) UpdateEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.NotFound(w, r)
		return
	}

	var t DataToUpdateEvent
	err := json.NewDecoder(r.Body).Decode(&t)
	if err != nil {
		sendError(w, http.StatusBadRequest, nil)
		return
	}
	if !t.isValid() {
		sendError(w, http.StatusBadRequest, nil)
		return
	}

	ctx := context.Background()
	if t.Name != "" {
		ctx = context.WithValue(ctx, name, t.Name)
	}
	if t.UserID != 0 {
		ctx = context.WithValue(ctx, userID, t.UserID)
	}
	if t.Date != "" {
		ctx = context.WithValue(ctx, date, t.Date)
	}

	updated, err := app.events.updateEvent(t.ID, ctx)
	if err != nil {
		sendError(w, http.StatusBadRequest, nil)
		return
	}
	sendResponse(w, http.StatusOK, updated)
}

type DataToDeleteEvent struct {
	ID int `json:"id"`
}

func (app *Application) DeleteEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.NotFound(w, r)
		return
	}

	var t DataToDeleteEvent
	err := json.NewDecoder(r.Body).Decode(&t)
	if err != nil {
		sendError(w, http.StatusBadRequest, nil)
		return
	}

	deleted, err := app.events.DeleteEvent(t.ID)
	if err != nil {
		sendError(w, http.StatusBadRequest, nil)
		return
	}
	sendResponse(w, http.StatusOK, deleted)
}

func (app *Application) getEventsForDay(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.NotFound(w, r)
		return
	}

	dayStr, ok := r.URL.Query()["day"]
	if !ok {
		sendError(w, http.StatusBadRequest, nil)
		return
	}
	day, err := time.Parse("2006-01-02", dayStr[0])
	if err != nil {
		sendError(w, http.StatusBadRequest, nil)
		return
	}

	events, err := app.events.findByDay(day)
	if err != nil {
		sendError(w, http.StatusOK, nil)
		return
	}

	sendResponse(w, http.StatusOK, events)
}

func (app *Application) getEventsForWeek(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.NotFound(w, r)
		return
	}

	weekStr, ok := r.URL.Query()["week"]
	if !ok {
		sendError(w, http.StatusBadRequest, nil)
		return
	}
	week, err := time.Parse("2006-01-02", weekStr[0])
	if err != nil {
		sendError(w, http.StatusBadRequest, nil)
		return
	}

	events, err := app.events.findByWeek(week)
	if err != nil {
		sendError(w, http.StatusOK, nil)
		return
	}

	sendResponse(w, http.StatusOK, events)
}

func (app *Application) getEventsForMonth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.NotFound(w, r)
		return
	}

	weekStr, ok := r.URL.Query()["month"]
	if !ok {
		sendError(w, http.StatusBadRequest, nil)
		return
	}
	week, err := time.Parse("2006-01-02", weekStr[0])
	if err != nil {
		sendError(w, http.StatusBadRequest, nil)
		return
	}

	events, err := app.events.findByMonth(week)
	if err != nil {
		sendError(w, http.StatusOK, nil)
		return
	}

	sendResponse(w, http.StatusOK, events)
}

func (app *Application) getEventsForYear(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.NotFound(w, r)
		return
	}

	yearStr, ok := r.URL.Query()["year"]
	if !ok {
		sendError(w, http.StatusBadRequest, nil)
		return
	}
	year, err := strconv.Atoi(yearStr[0])
	if err != nil {
		sendError(w, http.StatusBadRequest, nil)
		return
	}

	events, err := app.events.findByYear(year)
	if err != nil {
		sendError(w, http.StatusOK, nil)
		return
	}

	sendResponse(w, http.StatusOK, events)
}

func sendResponse(w http.ResponseWriter, code int, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(code)

	response := struct {
		Result interface{} `json:"result"`
	}{
		Result: data,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func sendError(w http.ResponseWriter, code int, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(code)

	response := struct {
		Error interface{} `json:"error"`
	}{
		Error: data,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}
