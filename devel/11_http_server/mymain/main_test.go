package mymain

import (
	"config"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"
)

var app Application

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

func TestMain(m *testing.M) {
	cfg := config.NewTestConfig()
	app.initDateStorage(cfg)
	app.initHTTPServer(cfg)
	defer app.Shutdown()

	exitVal := m.Run()

	os.Exit(exitVal)
}

func Test_application_getEvent(t *testing.T) {
	request, _ := http.NewRequest(http.MethodGet, "/event", nil)
	response := httptest.NewRecorder()

	app.getEvent(response, request)
	checkResponseCode(t, http.StatusOK, response.Code)

	eventsGot := struct {
		Result []Event `json:"result"`
	}{}
	err := json.Unmarshal(response.Body.Bytes(), &eventsGot)
	if err != nil {
		t.Errorf("JSON invalid")
	}

	eventsExpected := []Event{
		{ID: 1, UserID: 100, Name: "first", Date: "2020-04-30"},
		{ID: 2, UserID: 100, Name: "second", Date: "2021-05-20"},
		{ID: 3, UserID: 200, Name: "thrid", Date: "2021-07-10"},
	}

	if !reflect.DeepEqual(eventsExpected, eventsGot.Result) {
		t.Errorf("Expected %v array. Got %v", eventsExpected, eventsGot)
	}
}

func Test_application_getEventByID(t *testing.T) {
	request, _ := http.NewRequest(http.MethodGet, "/event", nil)
	q := request.URL.Query()
	q.Add("id", "3")
	request.URL.RawQuery = q.Encode()
	response := httptest.NewRecorder()

	app.getEvent(response, request)
	checkResponseCode(t, http.StatusOK, response.Code)

	eventsGot := struct {
		Result Event `json:"result"`
	}{}
	err := json.Unmarshal(response.Body.Bytes(), &eventsGot)
	if err != nil {
		t.Errorf("JSON invalid")
	}
	fmt.Println(eventsGot)

	eventsExpected := Event{ID: 3, UserID: 200, Name: "thrid", Date: "2021-07-10"}

	if !reflect.DeepEqual(eventsExpected, eventsGot.Result) {
		t.Errorf("Expected %v array. Got %v", eventsExpected, eventsGot)
	}
}
