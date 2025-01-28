package server

import (
	"calendar-server/config"
	eventstorage "calendar-server/eventStorage"
	"calendar-server/filedb"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

func Test_application_getEventByID(t *testing.T) {
	cfg := config.NewDefaultConfig()
	db, err := filedb.New("test_db.txt")
	if err != nil {
		log.Fatalf("NewFileDB: %s", err.Error())
	}

	es, err := eventstorage.New(db)
	if err != nil {
		log.Fatalf("eventstorage: %s", err.Error())
	}
	server := New(*cfg, es)
	// init ended

	request, _ := http.NewRequest(http.MethodGet, "/event", nil)
	q := request.URL.Query()
	q.Add("id", "3")
	request.URL.RawQuery = q.Encode()
	response := httptest.NewRecorder()

	server.getEvent(response, request)
	checkResponseCode(t, http.StatusOK, response.Code)

	b := response.Body.Bytes()
	fmt.Println(string(b))
	eventsGot := struct {
		Result struct {
			ID     int    `json:"id"`
			UserID int    `json:"user_id"`
			Name   string `json:"name"`
			Date   string `json:"date"`
		} `json:"result"`
	}{}
	if err := json.Unmarshal(b, &eventsGot); err != nil {
		t.Errorf("JSON invalid: %s", err.Error())
	}

	eventsExpected := Event{ID: 3, UserID: 200, Name: "three", Date: "2024-11-30"}

	if eventsExpected == eventsGot.Result {
		t.Errorf("Expected %v array. Got %v", eventsExpected, eventsGot.Result)
	}
}

func Test_application_getEventByYear(t *testing.T) {
	cfg := config.NewDefaultConfig()
	db, err := filedb.New("test_db.txt")
	if err != nil {
		log.Fatalf("NewFileDB: %s", err.Error())
	}

	es, err := eventstorage.New(db)
	if err != nil {
		log.Fatalf("eventstorage: %s", err.Error())
	}
	server := New(*cfg, es)
	// init ended

	request, _ := http.NewRequest(http.MethodGet, "/events_for_year", nil)
	q := request.URL.Query()
	q.Add("year", "2024")
	request.URL.RawQuery = q.Encode()
	response := httptest.NewRecorder()

	server.getEventsForYear(response, request)
	checkResponseCode(t, http.StatusOK, response.Code)

	b := response.Body.Bytes()
	fmt.Println(string(b))

	type Result struct {
		Result []Event `json:"result"`
	}
	var eventsGot Result
	if err := json.Unmarshal(b, &eventsGot); err != nil {
		t.Errorf("JSON invalid: %s", err.Error())
	}

	eventsExpected := Result{Result: []Event{
		{ID: 1, UserID: 100, Name: "first", Date: "2024-12-30"},
		{ID: 2, UserID: 100, Name: "second", Date: "2024-12-20"},
		{ID: 3, UserID: 100, Name: "three", Date: "2024-11-30"},
	},
	}

	if !reflect.DeepEqual(eventsExpected, eventsGot) {
		t.Errorf("Expected %v array. Got %v", eventsExpected, eventsGot.Result)
	}
}
