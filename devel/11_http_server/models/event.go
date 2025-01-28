package models

type eventcontextKey int

const (
	Name eventcontextKey = iota
	UserID
	Date
)

// Теги для хранения в файле
type EventData struct {
	ID     int    `json:"id"`
	UserID int    `json:"user_id"`
	Name   string `json:"name"`
	Date   string `json:"date"`
}

type NewEventData struct {
	UserID int
	Name   string
	Date   string
}

type UpdateEventData struct {
	ID     int
	UserID *int
	Name   *string
	Date   *string
}
