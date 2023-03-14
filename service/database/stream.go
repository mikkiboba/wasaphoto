package database

import "database/sql"

func (db *appdbimpl) GetStream(followingList []string) ([]Post, error) {

	// Tansform the followingList array to a string (to use in a query)
	var followingS string
	for index := range followingList {
		if index == 0 {
			followingS = "'" + followingList[index] + "'"
		} else {
			followingS += "," + "'" + followingList[index] + "'"
		}
	}

	var posts []Post
	rows, err := db.c.Query(`SELECT id, user, date, hour FROM posts WHERE user IN (` + followingS + `) ORDER BY date DESC`)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var post Post
		err = rows.Scan(&post.PostID, &post.User, &post.Date, &post.Hour)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	if len(posts) == 0 {
		return nil, sql.ErrNoRows
	}

	for index := range posts {
		// Get the number of likes
		likes, err := db.GetLikesCounter(posts[index].PostID)
		if err != nil {
			return nil, err
		}

		// Get the number of comments
		comments, err := db.GetCommentsCounter(posts[index].PostID)
		if err != nil {
			return nil, err
		}

		// Get the username of the post's owner
		username, err := db.GetUsername(posts[index].User)
		if err != nil {
			return nil, err
		}

		posts[index].LikesNumber = likes
		posts[index].CommentsNumber = comments
		posts[index].User = username
	}
	return posts, nil

}
