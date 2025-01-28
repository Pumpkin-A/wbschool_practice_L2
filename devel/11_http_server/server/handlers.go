package server

import (
	"calendar-server/models"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

var (
	ErrBadID     error = fmt.Errorf("id is bad")
	ErrBadUserID error = fmt.Errorf("user_id is bad")
	ErrBadName   error = fmt.Errorf("name is bad")
	ErrBadDate   error = fmt.Errorf("date is bad")
	ErrBadJson   error = fmt.Errorf("json is not parsed")
)

type Event struct {
	ID     int    `json:"id"`
	UserID int    `json:"user_id"`
	Name   string `json:"name"`
	Date   string `json:"date"`
}

func convertEvent(data models.EventData) Event {
	return Event{
		ID:     data.ID,
		UserID: data.UserID,
		Name:   data.Name,
		Date:   data.Date,
	}
}

func convertEvents(datas []models.EventData) []Event {
	res := make([]Event, 0, len(datas))
	for _, e := range datas {
		res = append(res, convertEvent(e))
	}
	return res
}

func (s *Server) getEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.NotFound(w, r)
		return
	}

	eventIDs, ok := r.URL.Query()["id"]
	if !ok {
		sendError(w, http.StatusBadRequest, ErrBadID.Error())
		return
	}

	eventID, err := strconv.Atoi(eventIDs[0])
	if err != nil {
		sendError(w, http.StatusBadRequest, ErrBadID.Error())
		return
	}

	event, err := s.events.GetEvent(eventID)
	if err != nil {
		sendError(w, http.StatusServiceUnavailable, err.Error())
		return
	}

	sendResponse(w, http.StatusOK, convertEvent(event))
}

type AddEventRequest struct {
	UserID int    `json:"user_id"`
	Name   string `json:"name"`
	Date   string `json:"date"`
}

func (d AddEventRequest) isValid() error {
	if d.UserID <= 0 {
		return ErrBadID
	}
	if d.Name == "" {
		return ErrBadName
	}
	if _, err := time.Parse("2006-01-02", d.Date); err != nil {
		return ErrBadDate
	}
	return nil
}

func (s *Server) AddEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.NotFound(w, r)
		return
	}

	var req AddEventRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		sendError(w, http.StatusBadRequest, ErrBadJson.Error())
		return
	}
	if err := req.isValid(); err != nil {
		sendError(w, http.StatusBadRequest, err.Error())
		return
	}

	data := convertAddEventRequest(req)

	ID, err := s.events.AddEvent(data)
	if err != nil {
		sendError(w, http.StatusServiceUnavailable, err.Error())
		return
	}

	sendResponse(w, http.StatusOK, struct {
		ID int `json:"id"`
	}{ID: ID})
}

func convertAddEventRequest(req AddEventRequest) models.NewEventData {
	return models.NewEventData{
		UserID: req.UserID,
		Name:   req.Name,
		Date:   req.Date,
	}
}

type UpdateEventRequest struct {
	ID     int    `json:"id"`
	UserID int    `json:"user_id"`
	Name   string `json:"name"`
	Date   string `json:"date"`
}

func (d UpdateEventRequest) isValid() error {
	if d.ID <= 0 {
		return ErrBadID
	}
	if d.UserID < 0 {
		return ErrBadUserID
	}
	if d.Date != "" {
		if _, err := time.Parse("2006-01-02", d.Date); err != nil {
			return ErrBadJson
		}
	}
	// Должно быть хотя бы одно новое поле
	if d.Date == "" && d.ID == 0 && d.UserID == 0 {
		return fmt.Errorf("new field must exist")
	}
	return nil
}

func (s *Server) UpdateEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.NotFound(w, r)
		return
	}

	var req UpdateEventRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		sendError(w, http.StatusBadRequest, ErrBadJson.Error())
		return
	}
	if err := req.isValid(); err != nil {
		sendError(w, http.StatusBadRequest, err.Error())
		return
	}

	data := convertUpdateEventRequest(req)

	updated, err := s.events.UpdateEvent(data)
	if err != nil {
		sendError(w, http.StatusServiceUnavailable, err.Error())
		return
	}
	sendResponse(w, http.StatusOK, convertEvent(updated))
}

func convertUpdateEventRequest(req UpdateEventRequest) models.UpdateEventData {
	var data models.UpdateEventData

	data.ID = req.ID

	if req.Name != "" {
		data.Name = &req.Name
	}

	if req.UserID != 0 {
		data.UserID = &req.UserID
	}

	if req.Date != "" {
		data.Date = &req.Date
	}

	return data
}

type DataToDeleteEvent struct {
	ID int `json:"id"`
}

func (d DataToDeleteEvent) isValid() error {
	if d.ID <= 0 {
		return ErrBadID
	}
	return nil
}

func (s *Server) DeleteEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.NotFound(w, r)
		return
	}

	var req DataToDeleteEvent
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		sendError(w, http.StatusBadRequest, ErrBadJson.Error())
		return
	}
	if err := req.isValid(); err != nil {
		sendError(w, http.StatusBadRequest, err.Error())
		return
	}

	deleted, err := s.events.DeleteEvent(req.ID)
	if err != nil {
		sendError(w, http.StatusServiceUnavailable, err.Error())
		return
	}
	sendResponse(w, http.StatusOK, convertEvent(deleted))
}

func (s *Server) getEventsForDay(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.NotFound(w, r)
		return
	}

	dayStr, ok := r.URL.Query()["day"]
	if !ok {
		sendError(w, http.StatusBadRequest, "day is bad")
		return
	}
	day, err := time.Parse("2006-01-02", dayStr[0])
	if err != nil {
		sendError(w, http.StatusBadRequest, "day not correct")
		return
	}

	events, err := s.events.FindByDay(day)
	if err != nil {
		sendError(w, http.StatusServiceUnavailable, err.Error())
		return
	}

	sendResponse(w, http.StatusOK, convertEvents(events))
}

func (s *Server) getEventsForWeek(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.NotFound(w, r)
		return
	}

	weekStr, ok := r.URL.Query()["week"]
	if !ok {
		sendError(w, http.StatusBadRequest, "week is bad")
		return
	}
	week, err := time.Parse("2006-01-02", weekStr[0])
	if err != nil {
		sendError(w, http.StatusBadRequest, "week is bad")
		return
	}

	events, err := s.events.FindByWeek(week)
	if err != nil {
		sendError(w, http.StatusServiceUnavailable, err.Error())
		return
	}

	sendResponse(w, http.StatusOK, convertEvents(events))
}

func (s *Server) getEventsForMonth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.NotFound(w, r)
		return
	}

	monthStr, ok := r.URL.Query()["month"]
	if !ok {
		sendError(w, http.StatusBadRequest, "month is bad")
		return
	}
	week, err := time.Parse("2006-01-02", monthStr[0])
	if err != nil {
		sendError(w, http.StatusBadRequest, "month is bad")
		return
	}

	events, err := s.events.FindByMonth(week)
	if err != nil {
		sendError(w, http.StatusServiceUnavailable, err.Error())
		return
	}

	sendResponse(w, http.StatusOK, convertEvents(events))
}

func (s *Server) getEventsForYear(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.NotFound(w, r)
		return
	}

	yearStr, ok := r.URL.Query()["year"]
	if !ok {
		sendError(w, http.StatusBadRequest, "year is bad")
		return
	}
	year, err := strconv.Atoi(yearStr[0])
	if err != nil {
		sendError(w, http.StatusBadRequest, "year is bad")
		return
	}

	events, err := s.events.FindByYear(year)
	if err != nil {
		sendError(w, http.StatusServiceUnavailable, err.Error())
		return
	}

	sendResponse(w, http.StatusOK, convertEvents(events))
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

func sendError(w http.ResponseWriter, code int, errText string) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(code)

	response := struct {
		Error string `json:"error"`
	}{
		Error: errText,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}
