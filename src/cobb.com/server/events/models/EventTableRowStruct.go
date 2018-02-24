package models

import "time"

type EventRow struct {
	ID          int       `pq:"id" json:"id"`
	Name        string    `pq:"name" json:"name"`
	Description string    `pq:"description" json:"description"`
	DateAdded   time.Time `pq:"date_added" json:"date"`
}

type GetAllEventsResponseStruct struct {
	ID          int       `pq:"id" json:"id"`
	Name        string    `pq:"name" json:"name"`
	Description string    `pq:"description" json:"description"`
}

func NewEventRow(id int, name string, description string, dateAdded time.Time) *EventRow {
	return &EventRow{
		ID:          id,
		Name:        name,
		Description: description,
		DateAdded:   dateAdded,
	}
}
