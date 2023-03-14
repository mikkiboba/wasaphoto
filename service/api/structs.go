package api

import "git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/database"

type Username struct {
	Username string `json:"username"`
}

/*
Checing if the username is valid. A valid username has a length in the range [3,10].
*/
func (u *Username) isValid() bool {
	return len(u.Username) >= 3 && len(u.Username) <= 10
}

type Token struct {
	ID string `json:"token"`
}

type User struct {
	Username         string `json:"username"`
	FollowCounter    int    `json:"follow-counter"`
	FollowingCounter int    `json:"following-counter"`
	PostCounter      int    `json:"post-counter"`
}

type Post struct {
	PostID         int    `json:"id"`
	User           string `json:"user"`
	Date           string `json:"date"`
	Hour           string `json:"hour"`
	LikesNumber    int    `json:"likes"`
	CommentsNumber int    `json:"comments"`
}

func FromDatabase(post database.Post) Post {
	var p Post
	p.PostID = post.PostID
	p.User = post.User
	p.Date = post.Date
	p.Hour = post.Hour
	p.LikesNumber = post.LikesNumber
	p.CommentsNumber = post.CommentsNumber
	return p
}

type Comment struct {
	Text string `json:"comment"`
}
