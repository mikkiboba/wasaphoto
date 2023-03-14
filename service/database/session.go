package database

/*
Login the user into the system (only if the username is in the database) returning the uuid (authorization token).
*/
func (db *appdbimpl) Login(username string) (string, error) {
	var uid string
	row := db.c.QueryRow(`SELECT id FROM users WHERE username = ?`, username)
	err := row.Scan(&uid)
	return uid, err
}

/*
Register the user into the database with the uuid and the username.
*/
func (db *appdbimpl) Registration(uid string, username string) error {
	_, err := db.c.Exec(`INSERT INTO users (id, username) VALUES (?,?)`, uid, username)
	return err
}
