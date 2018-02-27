package services

import (
	"cobb.com/server/database"
	"cobb.com/server/events/models"
	"database/sql"
	"errors"
	"time"
)

type EventsServiceClient struct {
	DB database.DBClient
}

func NewClient(dbClient database.Client) *EventsServiceClient {
	return &EventsServiceClient{
		DB: dbClient.DbClient,
	}
}

func (es *EventsServiceClient) AddCannedRowToEventsDB() (interface{}, sql.Result, error) {
	now := time.Now()
	result, err := es.DB.Exec("INSERT INTO event (name, description, date_added) VALUES ($1, $2, $3);",
		"Coheed & Cambria with Taking Back Sunday",
		"Prog Rock favorites Coheed & Cambria are back with a new national tour!",
		now)
	if err != nil {
		return nil, nil, err
	}
	return now, result, err
}
func (es *EventsServiceClient) AddNewEvent(name string, description string) (*models.EventRow, sql.Result, error) {
	now := time.Now()
	result, err := es.DB.Exec("INSERT INTO event (name, description, date_added) VALUES ($1, $2, $3);",
		name, description, now)
	if err != nil {
		return nil, nil, err
	}
	var event models.EventRow
	row := es.DB.QueryRow("SELECT id, name, description, date_added FROM event WHERE date_added = $1", now)
	if row != nil {
		if err = row.Scan(&event.ID, &event.Name, &event.Description, &event.DateAdded); err != nil {
			return nil, nil, err
		}
	} else {
		var rowsChanged int64
		rowsChanged, err = result.RowsAffected()
		if err != nil {
			return nil, nil, err
		}
		if rowsChanged > 0 {
			return nil, nil, errors.New("500: Something is broken in AddNewEvent fetching")
		} else {
			return nil, nil, errors.New("404: Event not found")
		}
	}
	return &event, result, err
}
func (es *EventsServiceClient) GetAllEvents() (*[]models.GetAllEventsResponseStruct, error) {
	rows, err := es.DB.Query("SELECT id, name, description FROM event;")
	if err != nil {
		return nil, err
	}
	eventRows := make([]models.GetAllEventsResponseStruct, 0)
	for rows.Next() {
		var eventRow models.GetAllEventsResponseStruct
		if err := rows.Scan(&eventRow.ID, &eventRow.Name, &eventRow.Description); err != nil {
			return nil, err
		}
		eventRows = append(eventRows, eventRow)
	}
	return &eventRows, nil
}
