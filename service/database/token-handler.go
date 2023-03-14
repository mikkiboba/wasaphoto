package database

import (
	"database/sql"
	"errors"
)

/*
Check if the token is inside the database
*/
func (db *appdbimpl) CheckToken(token string) (bool, error) {
	var id string
	err := db.c.QueryRow(`SELECT id FROM users WHERE id = ?`, token).Scan(&id)
	if errors.Is(err, sql.ErrNoRows) {
		return false, err
	}
	return true, err
}

/*
Get the token related to the input username
*/
func (db *appdbimpl) GetToken(username string) (string, error) {
	var token string
	err := db.c.QueryRow(`SELECT id FROM users WHERE username = ?`, username).Scan(&token)
	return token, err
}
