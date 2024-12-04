package models

import (
	"time"

	sq "github.com/Masterminds/squirrel"

	"example.com/events-api/db"
)

type Event struct {
	ID          int64
	Name        string    `binding:"required"`
	Description string    `binding:"required"`
	Location    string    `binding:"required"`
	DateTime    time.Time `binding:"required"`
	UserID      int64
}

type PatchEventRequest struct {
	Name        *string
	Description *string
	Location    *string
	DateTime    *time.Time
}

func GetAllEvents() ([]Event, error) {
	query, _, err := sq.Select("*").From("events").ToSql()

	if err != nil {
		return nil, err
	}

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
	query, params, err := sq.
		Select("*").
		From("events").
		Where(sq.Eq{"id": id}).
		ToSql()

	if err != nil {
		return nil, err
	}

	row := db.DB.QueryRow(query, params...)

	var event Event
	err = row.Scan(
		&event.ID, &event.Name, &event.Description, &event.Location, &event.DateTime, &event.UserID,
	)

	if err != nil {
		return nil, err
	}

	return &event, nil
}

func (event *Event) Save() error {
	query, params, err := sq.
		Insert("events").
		Columns("name", "description", "location", "dateTime", "userId").
		Values(event.Name, event.Description, event.Location, event.DateTime, event.UserID).
		ToSql()

	if err != nil {
		return err
	}

	statement, err := db.DB.Prepare(query)

	if err != nil {
		return err
	}

	defer statement.Close()

	result, err := statement.Exec(params...)

	if err != nil {
		return err
	}

	id, err := result.LastInsertId()

	if err != nil {
		return err
	}

	event.ID = id

	return nil
}

func (event *Event) Update() error {
	query, params, err := sq.
		Update("events").
		Set("name", event.Name).
		Set("description", event.Description).
		Set("location", event.Location).
		Set("dateTime", event.DateTime).
		Where(sq.Eq{"id": event.ID}).
		ToSql()

	if err != nil {
		return err
	}

	statement, err := db.DB.Prepare(query)

	if err != nil {
		return err
	}

	defer statement.Close()

	_, err = statement.Exec(params...)

	return err
}

func (event *Event) Patch(patchRequest PatchEventRequest) error {
	queryBuilder := sq.Update("events")
	noData :=
		patchRequest.Name == nil &&
			patchRequest.Description == nil &&
			patchRequest.Location == nil &&
			patchRequest.DateTime == nil

	if noData {
		return nil
	}

	if patchRequest.Name != nil {
		queryBuilder = queryBuilder.Set("name", *patchRequest.Name)
	}
	if patchRequest.Description != nil {
		queryBuilder = queryBuilder.Set("description", *patchRequest.Description)
	}
	if patchRequest.Location != nil {
		queryBuilder = queryBuilder.Set("location", *patchRequest.Location)
	}
	if patchRequest.DateTime != nil {
		queryBuilder = queryBuilder.Set("dateTime", *patchRequest.DateTime)
	}

	queryBuilder = queryBuilder.Where(sq.Eq{"id": event.ID})

	query, params, err := queryBuilder.ToSql()

	if err != nil {
		return err
	}

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
	query, params, err := sq.
		Delete("events").
		Where(
			sq.Eq{"id": event.ID},
		).
		ToSql()

	if err != nil {
		return err
	}

	statement, err := db.DB.Prepare(query)

	if err != nil {
		return err
	}

	defer statement.Close()

	_, err = statement.Exec(params...)

	return err
}
