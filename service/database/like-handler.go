package database

func (db *appdbimpl) LikePhoto(uid string, pid string) error {
	res, err := db.c.Exec(`INSERT INTO likes (postid, userliking) VALUES (?, ?)`, pid, uid)
	if err != nil {
		return err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return err
	} else if affected == 0 {
		return ErrElementNotAdded
	}
	return nil
}

func (db *appdbimpl) UnlikePhoto(uid string, pid string) error {
	res, err := db.c.Exec(`DELETE FROM likes WHERE postid = ? AND userliking = ?`, pid, uid)
	if err != nil {
		return err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return err
	} else if affected == 0 {
		return ErrElementNotDeleted
	}
	return nil
}

func (db *appdbimpl) CheckLike(uid string, pid string) (bool, error) {
	var res int
	err := db.c.QueryRow(`SELECT COUNT(*) FROM likes WHERE postid = ? AND userliking = ?`, pid, uid).Scan(&res)
	return res>0, err
}