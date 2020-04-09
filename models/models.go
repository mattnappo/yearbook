package models

import (
	"errors"
	"io/ioutil"
	"strings"
	"time"
)

// PostID is the main form of identification for a post.
type PostID string

// User represents a user.
type User struct {
	ID           string `pg:",notnull"`
	Firstname    string `pg:",notnull"`
	Lastname     string `pg:",notnull"`
	Username     string `pg:",notnull"`
	Email        string `pg:",notnull"`
	RegisterDate string `pg:",notnull"`
}

// Post represents a post in the database.
type Post struct {
	PostID     PostID  `pg:",notnull"`
	Timestamp  string  `pg:",notnull"`
	Sender     *User   `pg:",notnull"`
	Recipients []*User `pg:",notnull"`

	Message string   `pg:",notnull"`
	Images  [][]byte `pg:",notnull"`
}

// NewUserFromEmail creates a *User given a valid email.
func NewUserFromEmail(email string) (*User, error) {
	if email[len(email)-len(emailPrefix):] == emailPrefix {
		userData := strings.Split(email, ".")
		return &User{
			ID:           "hash of email",
			Firstname:    userData[0],
			Lastname:     userData[1],
			Username:     email[0 : len(email)-len(emailPrefix)],
			Email:        email,
			RegisterDate: time.Now().String(),
		}, nil
	}
	return nil, errors.New("invalid email address")
}

// NewPost creates a new post.
func NewPost(
	sender *User,
	recipients []*User,
	message string,
	imagePaths []string,
) (*Post, error) {
	// Check that no data is nil
	if sender == nil || len(recipients) > common.MaxRecipients ||
		len(recipients) <= 0 || message == "" || len(images) > common.MaxImages {
		return nil, errors.New("cannot supply nil data to construct post")
	}

	// Load the images into a [][]byte
	var images [][]byte
	for _, imagePath := range images {
		image, err := ioutil.ReadFile(imagePath)
		if err != nil {
			return nil, err
		}
		images = append(images, image)
	}
}
