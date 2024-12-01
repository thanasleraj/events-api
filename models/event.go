package models

import (
	"time"

	"example.com/events-api/db"
)

type Event struct {
	ID          int64
	Name        string    `binding:"required"`
	Description string    `binding:"required"`
	Location    string    `binding:"required"`
	DateTime    time.Time `binding:"required"`
	UserID      int
}

func (event *Event) Save() error {
	query := `
	INSERT INTO events (name, description, location, dateTime, userId)
	VALUES (?, ?, ?, ?, ?)
	`
	statement, err := db.DB.Prepare(query)

	if err != nil {
		return err
	}

	defer statement.Close()

	result, err := statement.Exec(
		event.Name, event.Description, event.Location, event.DateTime, event.UserID,
	)

	if err != nil {
		return err
	}

	id, err := result.LastInsertId()

	event.ID = id

	return err
}

func GetAllEvents() ([]Event, error) {
	query := "SELECT * FROM events"

	rows, err := db.DB.Query(query)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	events := []Event{}

	for rows.Next() {
		var event Event
		err = rows.Scan(
			&event.ID, &event.Name, &event.Description, &event.Location, &event.DateTime, &event.UserID,
		)

		if err != nil {
			return nil, err
		}

		events = append(events, event)
	}

	return events, nil
}

func GetEventById(id int64) (*Event, error) {
	query := "SELECT * FROM events WHERE id = ?"
	row := db.DB.QueryRow(query, id)

	var event Event
	err := row.Scan(
		&event.ID, &event.Name, &event.Description, &event.Location, &event.DateTime, &event.UserID,
	)

	if err != nil {
		return nil, err
	}

	return &event, nil
}

func (event *Event) Update() error {
	query := `
	UPDATE events 
	SET name = ?, description = ?, location = ?, dateTime = ?
	WHERE id = ?
	`
	statement, err := db.DB.Prepare(query)

	if err != nil {
		return err
	}

	defer statement.Close()

	_, err = statement.Exec(
		event.Name, event.Description, event.Location, event.DateTime, event.ID,
	)

	return err
}

func PatchEvent(id int64, fields map[string]interface{}) (*Event, error) {
	query := "UPDATE events SET "
	params := []interface{}{}
	index := 1

	for key, value := range fields {
		if index > 1 {
			query += ", "
		}
		query += key + " = ?"
		params = append(params, value)
		index++
	}

	query += " WHERE id = ?"
	params = append(params, id)

	statement, err := db.DB.Prepare(query)

	if err != nil {
		return nil, err
	}

	defer statement.Close()

	_, err = statement.Exec(params...)

	if err != nil {
		return nil, err
	}

	updatedEvent, err := GetEventById(id)

	if err != nil {
		return nil, err
	}

	return updatedEvent, err
}

func (event *Event) Delete() error {
	query := "DELETE FROM events WHERE id = ?"
	statement, err := db.DB.Prepare(query)

	if err != nil {
		return err
	}

	defer statement.Close()

	_, err = statement.Exec(event.ID)

	return err
}
