package models

import (
	sq "github.com/Masterminds/squirrel"

	"example.com/events-api/db"
	"example.com/events-api/utils"
)

type User struct {
	ID       int64
	Email    string `binding:"required"`
	Password string `binding:"required"`
}

func (user *User) Save() error {
	hashedPassword, err := utils.HashPassword(user.Password)

	if err != nil {
		return err
	}

	query, args, err := sq.
		Insert("users").
		Columns("email", "password").
		Values(user.Email, hashedPassword).
		ToSql()

	if err != nil {
		return err
	}

	statement, err := db.DB.Prepare(query)

	if err != nil {
		return err
	}

	defer statement.Close()

	result, err := statement.Exec(args...)

	if err != nil {
		return err
	}

	userId, err := result.LastInsertId()

	if err != nil {
		return err
	}

	user.ID = userId

	return nil
}
