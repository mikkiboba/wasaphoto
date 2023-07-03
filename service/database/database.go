/*
Package database is the middleware between the app database and the code. All data (de)serialization (save/load) from a
persistent database are handled here. Database specific logic should never escape this package.

To use this package you need to apply migrations to the database if needed/wanted, connect to it (using the database
data source name from config), and then initialize an instance of AppDatabase from the DB connection.

For example, this code adds a parameter in `webapi` executable for the database data source name (add it to the
main.WebAPIConfiguration structure):

	DB struct {
		Filename string `conf:""`
	}

This is an example on how to migrate the DB and connect to it:

	// Start Database
	logger.Println("initializing database support")
	db, err := sql.Open("sqlite3", "./foo.db")
	if err != nil {
		logger.WithError(err).Error("error opening SQLite DB")
		return fmt.Errorf("opening SQLite: %w", err)
	}
	defer func() {
		logger.Debug("database stopping")
		_ = db.Close()
	}()

Then you can initialize the AppDatabase and pass it to the api package.
*/
package database

import (
	"database/sql"
	"errors"
	"fmt"
)

// AppDatabase is the high level interface for the DB
type AppDatabase interface {
	Login(string) (string, error)      // Login the user into the system (only if the username is in the database) returning the uuid (authorization token). Login(username string)
	Registration(string, string) error // Register the user into the database with the uuid and the username. Registration(uid string, username string)

	CheckUsername(string) (bool, error)       // Check if the username is already in the database. CheckUsername(username string)
	CheckToken(string) (bool, error)          // Check if the token is inside the database. CheckToken(token string)
	CheckFollow(string, string) (bool, error) // Check if a user already follows another.
	CheckBan(string, string) (bool, error)    // Check if a user is already banned.

	GetToken(string) (string, error)           // Get the token related to the input username. GetToken(username string)
	GetFollowNumber(string) (int, error)       // Get the number of follows
	GetFollowingNumber(string) (int, error)    // Get the number of following
	GetPostNumber(string) (int, error)         // Get the number of posts
	GetPostOwner(string) (string, error)       // Get the id of the owner of the post
	GetFollowingList(string) ([]string, error) // Get the list of following of a user
	GetLikesCounter(int) (int, error)          // Get the number of likes for a post
	GetCommentsCounter(int) (int, error)       // Get the number of comments
	GetStream([]string) ([]Post, error)        // Get the stream of a user
	GetUsername(string) (string, error)        // Get the username from the uid

	SetUsername(string, string) error // Change the username of the specified user trough the uid.

	UploadPhoto(string, string) error   // Insert a photo into the database
	DeletePhoto(string) (string, error) // Delete a photo from the database

	FollowUser(string, string) error   // Follow a user
	UnfollowUser(string, string) error // Unfollow a user

	BanUser(string, string) error   // Ban a user
	UnbanUser(string, string) error // Unban a user

	LikePhoto(string, string) error   // Like a photo
	UnlikePhoto(string, string) error // Unlike a photo
	DeleteLikesBan(string, string) error // Delete all the likes after a ban

	Comment(string, string, string) error   // Comment a post
	Uncomment(string, string, string) error // Uncomment a comment
	DeleteCommentsBan(string, string) error // Delete all the comments after a ban
	CheckLike(string, string) (bool, error)

	// EXTRA

	GetPhoto(string) (string, error)           // Get the photo from the database
	GetUsers(string, string) ([]string, error) // Get the users from the database with the specified prefix

	// Transaction handling

	StartTransaction() error               // Start the transaction
	Commit() error                         // Commit the transaction
	Rollback() error                       // Rollback the transaction
	GetPosts(string) ([]Post, error)       // Get the list of posts of a user
	GetComments(string) ([]Comment, error) // Get the list of comments of a post

	Ping() error
}

// STRUCTS
type Post struct {
	PostID         int
	User           string
	Date           string
	Hour           string
	LikesNumber    int
	CommentsNumber int
}

type Comment struct {
	Id   string
	User string
	Text string
}

// ERRORS

// Error that indicates that an element was not added to the database
var ErrElementNotAdded = errors.New("the element wasn't added to the database")

// Error that indicates that an element was not deleted from the database
var ErrElementNotDeleted = errors.New("the elements wasn't deleted from the database")

// Error that indicates that a photo has not been found
var ErrPhotoNotFound = errors.New("the photo was not found")

type appdbimpl struct {
	c *sql.DB
}

// New returns a new instance of AppDatabase based on the SQLite connection `db`.
// `db` is required - an error will be returned if `db` is `nil`.
func New(db *sql.DB) (AppDatabase, error) {
	if db == nil {
		return nil, errors.New("database is required when building a AppDatabase")
	}

	// Check if table exists. If not, the database is empty, and we need to create the structure
	var tableName string
	err := db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='users';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {
		sqlStmt := ` PRAGMA foreign_keys = ON;
		BEGIN TRANSACTION;

		CREATE TABLE users (
			id TEXT NOT NULL PRIMARY KEY,
			username TEXT
		);

		CREATE TABLE posts (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user TEXT NOT NULL REFERENCES users (id) ON UPDATE CASCADE ON DELETE CASCADE,
			filename TEXT NOT NULL,
			date TEXT NOT NULL,
			hour TEXT NOT NULL
		);

		CREATE TABLE follows (
			userfollowing TEXT,
			userfollowed TEXT,
			PRIMARY KEY (userfollowing, userfollowed),
			FOREIGN KEY (userfollowing) REFERENCES users (id) ON UPDATE CASCADE ON DELETE CASCADE,
			FOREIGN KEY (userfollowed) REFERENCES users (id) ON UPDATE CASCADE ON DELETE CASCADE
		);

		CREATE TABLE bans (
			userbanning TEXT,
			userbanned TEXT,
			PRIMARY KEY (userbanning, userbanned),
			FOREIGN KEY (userbanning) REFERENCES users (id) ON UPDATE CASCADE ON DELETE CASCADE,
			FOREIGN KEY (userbanned) REFERENCES users (id) ON UPDATE CASCADE ON DELETE CASCADE
		);

		CREATE TABLE likes (
			postid INT,
			userliking TEXT,
			PRIMARY KEY (postid, userliking),
			FOREIGN KEY (postid) REFERENCES posts (id) ON UPDATE CASCADE ON DELETE CASCADE,
			FOREIGN KEY (userliking) REFERENCES users (id) ON UPDATE CASCADE ON DELETE CASCADE
		);

		CREATE TABLE comments(
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			postid INTEGER REFERENCES posts (id) ON UPDATE CASCADE ON DELETE CASCADE,
			userid TEXT REFERENCES users (id) ON UPDATE CASCADE ON DELETE CASCADE,
			text TEXT NOT NULL
		);

		COMMIT;
		`
		_, err = db.Exec(sqlStmt)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure: %w", err)
		}
	}

	return &appdbimpl{
		c: db,
	}, nil
}

func (db *appdbimpl) Ping() error {
	return db.c.Ping()
}
