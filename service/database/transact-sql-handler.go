package database

/*
Begin a transaction.
Remember to close it with Commit() or Rollback()
*/
func (db *appdbimpl) StartTransaction() error {
	_, err := db.c.Exec(`BEGIN TRANSACTION;`)
	return err
}

/*
End the transaction applying the changes
*/
func (db *appdbimpl) Commit() error {
	_, err := db.c.Exec(`COMMIT;`)
	return err
}

/*
End the transaction cancelling the changes
*/
func (db *appdbimpl) Rollback() error {
	_, err := db.c.Exec(`ROLLBACK;`)
	return err
}
