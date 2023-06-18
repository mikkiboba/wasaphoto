package database

/*
Check if a user (uid) already follows another (fid)
*/
func (db *appdbimpl) CheckFollow(uid string, fid string) (bool, error) {
	var res int
	err := db.c.QueryRow(`SELECT COUNT(*) FROM follows WHERE userfollowing = ? AND userfollowed = ?`, uid, fid).Scan(&res)
	return res > 0, err
}

/*
Follow the user
*/
func (db *appdbimpl) FollowUser(uid string, fid string) error {
	res, err := db.c.Exec(`INSERT INTO follows (userfollowing, userfollowed) VALUES (?,?)`, uid, fid)
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

/*
Unfollow the user
*/
func (db *appdbimpl) UnfollowUser(uid string, fid string) error {
	res, err := db.c.Exec(`DELETE FROM follows WHERE userfollowing = ? AND userfollowed = ?`, uid, fid)
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

func (db *appdbimpl) GetFollowingList(uid string) ([]string, error) {
	var followingList []string
	rows, err := db.c.Query(`SELECT userfollowed FROM follows WHERE userfollowing = ?`, uid)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var followed string
		err := rows.Scan(&followed)
		if err != nil {
			return nil, err
		}
		followingList = append(followingList, followed)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return followingList, nil
}
