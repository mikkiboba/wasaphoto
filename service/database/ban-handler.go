package database

func (db *appdbimpl) CheckBan(uid string, bid string) (bool, error) {
	var res int
	err := db.c.QueryRow(`SELECT COUNT(*) FROM bans WHERE userbanning = ? AND userbanned = ?`, uid, bid).Scan(&res)
	return res > 0, err
}

func (db *appdbimpl) BanUser(uid string, bid string) error {
	res, err := db.c.Exec(`INSERT INTO bans (userbanning, userbanned) VALUES (?, ?)`, uid, bid)
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

func (db *appdbimpl) UnbanUser(uid string, bid string) error {
	res, err := db.c.Exec(`DELETE FROM bans WHERE userbanning = ? AND userbanned = ?`, uid, bid)
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
