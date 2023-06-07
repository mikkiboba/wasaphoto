package database

/*
Comment a post
*/
func (db *appdbimpl) Comment(uid string, pid string, text string) error {
	res, err := db.c.Exec(`INSERT INTO comments (postid, userid, text) VALUES (?, ?, ?)`, pid, uid, text)
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
Delete a comment
*/
func (db *appdbimpl) Uncomment(uid string, pid string, cid string) error {
	res, err := db.c.Exec(`DELETE FROM comments WHERE id = ? AND userid = ? AND postid = ?`, cid, uid, pid)
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

/*
Delete all the comments of a user after a ban
*/
func (db *appdbimpl) DeleteCommentsBan(uid string, bid string) error {
	// Get the list of posts
	posts, err := db.GetStream([]string{uid})
	if err != nil {
		return err
	}

	for index := range posts {
		_, err := db.c.Exec(`DELETE FROM comments WHERE postid = ? AND userid = ?`, posts[index].PostID, bid)
		if err != nil {
			return err
		}
	}
	return nil
}

func (db *appdbimpl) GetComments(pid string) ([]Comment, error) {
	var commentList []Comment
	rows, err := db.c.Query(`SELECT id, text, userid FROM comments WHERE postid = ?`, pid)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var comment Comment
		err := rows.Scan(&comment.Id, &comment.Text, &comment.User)
		if err != nil {
			return nil, err
		}
		commentList = append(commentList, comment)
	}
	return commentList, nil
}