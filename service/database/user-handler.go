package database

/*
Check if the username is already in the database
*/
func (db *appdbimpl) CheckUsername(username string) (bool, error) {
	var isThere int
	row := db.c.QueryRow(`SELECT COUNT(*) FROM users WHERE username = ?`, username)
	err := row.Scan(&isThere)
	return isThere > 0, err
}

/*
Change the username of the specified user trough the uid
*/
func (db *appdbimpl) SetUsername(uid string, username string) error {
	_, err := db.c.Exec(`UPDATE users SET username = ? WHERE id = ?`, username, uid)
	return err
}

/*
Get the number of followers
*/
func (db *appdbimpl) GetFollowNumber(uid string) (int, error) {
	var followNumber int
	row := db.c.QueryRow(`SELECT COUNT(*) FROM follows WHERE userfollowed = ?`, uid)
	err := row.Scan(&followNumber)
	return followNumber, err
}

/*
Get the number of following
*/
func (db *appdbimpl) GetFollowingNumber(uid string) (int, error) {
	var followingNumber int
	row := db.c.QueryRow(`SELECT COUNT(*) FROM follows WHERE userfollowing = ?`, uid)
	err := row.Scan(&followingNumber)
	return followingNumber, err
}

/*
Get the number of posts
*/
func (db *appdbimpl) GetPostNumber(uid string) (int, error) {
	var postNumber int
	row := db.c.QueryRow(`SELECT COUNT(*) FROM posts WHERE user = ?`, uid)
	err := row.Scan(&postNumber)
	return postNumber, err
}

/*
Get the username from the uid
*/
func (db *appdbimpl) GetUsername(uid string) (string, error) {
	var username string
	err := db.c.QueryRow(`SELECT username FROM users WHERE id = ?`, uid).Scan(&username)
	return username, err
}
