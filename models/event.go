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

type PatchEventRequest struct {
	Name        *string
	Description *string
	Location    *string
	DateTime    *time.Time
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

func (event *Event) Patch(patchRequest PatchEventRequest) error {
	params := []interface{}{}
	query := "UPDATE events SET "

	if patchRequest.Name != nil {
		query += "name = ?"
		params = append(params, *patchRequest.Name)
	}
	if patchRequest.Description != nil {
		if len(params) > 0 {
			query += ", "
		}
		query += "description = ?"
		params = append(params, *patchRequest.Description)
	}
	if patchRequest.Location != nil {
		if len(params) > 0 {
			query += ", "
		}
		query += "location = ?"
		params = append(params, *patchRequest.Location)
	}
	if patchRequest.DateTime != nil {
		if len(params) > 0 {
			query += ", "
		}
		query += "dateTime = ?"
		params = append(params, *patchRequest.DateTime)
	}

	if len(params) == 0 {
		return nil
	}

	query += " WHERE id = ?"
	params = append(params, event.ID)

	statement, err := db.DB.Prepare(query)

	if err != nil {
		return err
	}

	defer statement.Close()

	_, err = statement.Exec(params...)

	if err != nil {
		return err
	}

	if patchRequest.Name != nil {
		event.Name = *patchRequest.Name
	}
	if patchRequest.Description != nil {
		event.Description = *patchRequest.Description
	}
	if patchRequest.Location != nil {
		event.Location = *patchRequest.Location
	}
	if patchRequest.DateTime != nil {
		event.DateTime = *patchRequest.DateTime
	}

	return nil
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
