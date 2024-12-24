package mymain

import (
	"config"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

type Application struct {
	events EventStorage
	db     DataManipulator
	server *http.Server
}

func (app *Application) initDateStorage(cfg *config.Config) {
	app.db = newFileDB(cfg.DbFilename)
	evs, err := app.db.openConnection()
	if err != nil {
		log.Fatalf("File db error open %v", err)
	}
	app.events.initEventStorage(evs)
}

func (app *Application) initHTTPServer(cfg *config.Config) {
	http.Handle("/event", logMiddleware(http.HandlerFunc(app.getEvent)))
	http.Handle("/create_event", logMiddleware(http.HandlerFunc(app.AddEvent)))
	http.Handle("/update_event", logMiddleware(http.HandlerFunc(app.UpdateEvent)))
	http.Handle("/delete_event", logMiddleware(http.HandlerFunc(app.DeleteEvent)))

	http.Handle("/events_for_day", logMiddleware(http.HandlerFunc(app.getEventsForDay)))
	http.Handle("/events_for_week", logMiddleware(http.HandlerFunc(app.getEventsForWeek)))
	http.Handle("/events_for_month", logMiddleware(http.HandlerFunc(app.getEventsForMonth)))
	http.Handle("/events_for_year", logMiddleware(http.HandlerFunc(app.getEventsForYear)))
	app.server = &http.Server{Addr: cfg.Port, Handler: http.DefaultServeMux}

}

func (app *Application) Shutdown() {
	err := app.db.closeConnction(app.events.data)
	if err != nil {
		log.Printf("Error while closing DB; %\nv", err)
	} else {
		log.Println("DB closed!")
	}

	err = app.server.Shutdown(context.Background())
	if err != nil {
		log.Printf("Error while closing server; %\nv", err)
	} else {
		log.Println("Server closed!")
	}
}

func StartWork() {
	cfg := config.NewDefaultConfig()

	var app Application
	app.initDateStorage(cfg)
	app.initHTTPServer(cfg)
	defer app.Shutdown()

	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, syscall.SIGINT, syscall.SIGTERM)
	serverErrors := make(chan error, 1)
	go func() {
		log.Printf("server: started at address %v", cfg.Port)
		serverErrors <- http.ListenAndServe(":8080", nil)
	}()

	for {
		select {
		case err := <-serverErrors:
			fmt.Printf("server error %v\n", err)
		case <-osSignals:
			fmt.Println("service shutdown...")
			return
		}
	}
}
