package database

import (
	"database/sql"
	"errors"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/globaltime"
)

/*
Insert a photo into the database
*/
func (db *appdbimpl) UploadPhoto(user string, filename string) error {
	res, err := db.c.Exec(`INSERT INTO posts (user, filename, date, hour) VALUES (?, ?, ?, ?)`, user, filename, globaltime.Now().Format("2006-01-02"), globaltime.Now().Format("15:04"))
	if err != nil {
		return err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return err
	} else if affected == 0 {
		return ErrElementNotAdded
	}
	return err
}

func (db *appdbimpl) DeletePhoto(pid string) (string, error) {
	var filename string
	err := db.c.QueryRow(`SELECT filename FROM posts WHERE id = ?`, pid).Scan(&filename)
	if err != nil {
		return "", nil
	}

	res, err := db.c.Exec(`DELETE FROM posts WHERE id = ?`, pid)
	if err != nil {
		return "", nil
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return "", err
	} else if affected == 0 {
		return "", ErrElementNotDeleted
	}
	return filename, nil
}

/*
Get the user who made a post
*/
func (db *appdbimpl) GetPostOwner(pid string) (string, error) {
	var uid string
	err := db.c.QueryRow(`SELECT user FROM posts WHERE id = ?`, pid).Scan(&uid)
	return uid, err
}

/*
Get the number of likes for a post
*/
func (db *appdbimpl) GetLikesCounter(pid int) (int, error) {
	var nlikes int
	err := db.c.QueryRow(`SELECT COUNT(*) FROM likes WHERE postid = ?`, pid).Scan(&nlikes)
	return nlikes, err
}

/*
Get the number of comments for a post
*/
func (db *appdbimpl) GetCommentsCounter(pid int) (int, error) {
	var ncomments int
	err := db.c.QueryRow(`SELECT COUNT(*) FROM comments WHERE postid = ?`, pid).Scan(&ncomments)
	return ncomments, err
}

func (db *appdbimpl) GetPhoto(pid string) (string, error) {
	var filename string
	err := db.c.QueryRow("SELECT filename FROM posts WHERE id = ?", pid).Scan(&filename)
	if errors.Is(err, sql.ErrNoRows) {
		return "", ErrPhotoNotFound
	} else if err != nil {
		return "", err
	}
	return filename, err
}
