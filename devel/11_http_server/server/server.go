package server

import (
	"calendar-server/config"
	"calendar-server/models"
	"context"
	"net/http"
	"time"
)

type EventService interface {
	GetEvent(ID int) (models.EventData, error)
	AddEvent(data models.NewEventData) (int, error)
	UpdateEvent(data models.UpdateEventData) (models.EventData, error)
	DeleteEvent(ID int) (models.EventData, error)
	FindByDay(day time.Time) ([]models.EventData, error)
	FindByWeek(week time.Time) ([]models.EventData, error)
	FindByMonth(month time.Time) ([]models.EventData, error)
	FindByYear(year int) ([]models.EventData, error)
}

type Server struct {
	events EventService
	server *http.Server
}

func New(cfg config.Config, event EventService) *Server {
	httpServer := &http.Server{Addr: cfg.Port, Handler: http.DefaultServeMux}

	s := &Server{
		events: event,
		server: httpServer,
	}

	http.Handle("/event", logMiddleware(http.HandlerFunc(s.getEvent)))
	http.Handle("/create_event", logMiddleware(http.HandlerFunc(s.AddEvent)))
	http.Handle("/update_event", logMiddleware(http.HandlerFunc(s.UpdateEvent)))
	http.Handle("/delete_event", logMiddleware(http.HandlerFunc(s.DeleteEvent)))

	http.Handle("/events_for_day", logMiddleware(http.HandlerFunc(s.getEventsForDay)))
	http.Handle("/events_for_week", logMiddleware(http.HandlerFunc(s.getEventsForWeek)))
	http.Handle("/events_for_month", logMiddleware(http.HandlerFunc(s.getEventsForMonth)))
	http.Handle("/events_for_year", logMiddleware(http.HandlerFunc(s.getEventsForYear)))

	return s
}

func (s *Server) Shutdown() error {
	return s.server.Shutdown(context.Background())
}
