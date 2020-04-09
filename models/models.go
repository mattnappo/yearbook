package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"github.com/xoreo/yearbook/common"
	"github.com/xoreo/yearbook/crypto"
)

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
	PostID     string  `pg:",notnull"`
	Timestamp  string  `pg:",notnull"`
	Sender     *User   `pg:",notnull"`
	Recipients []*User `pg:",notnull"`

	Message string   `pg:",notnull"`
	Images  [][]byte `pg:",notnull"`
}

// NewUserFromEmail creates a *User given a valid email.
func NewUserFromEmail(email string) (*User, error) {
	if email[len(email)-len(common.EmailSuffix):] == common.EmailSuffix {
		userData := strings.Split(email, ".")
		usernameLen := len(email) - len(common.EmailSuffix)
		return &User{
			ID:           crypto.Sha3String(email),
			Firstname:    userData[0],
			Lastname:     strings.Split(userData[1], "@")[0],
			Username:     email[0:usernameLen],
			Email:        email,
			RegisterDate: time.Now().String(),
		}, nil
	}
	return nil, errors.New("invalid email address")
}

// isValidUser checks if a User struct is valid
func (user *User) isValid() bool {
	l := len(user.Email) - len(common.EmailSuffix)
	if user.Email[l:] == common.EmailSuffix {

		return true
	}
	return false
}

// NewPost creates a new post.
func NewPost(
	sender *User,
	recipients []*User,
	message string,
	imagePaths []string,
) (*Post, error) {
	// Check that all data for the post is valid
	if sender == nil || len(recipients) > common.MaxRecipients ||
		len(recipients) <= 0 || message == "" ||
		len(message) > common.MaxMessageLength ||
		len(imagePaths) > common.MaxImages {
		return nil, errors.New("invalid data to construct post")
	}

	// Check that the users are valid
	for _, user := range recipients {
		if !user.isValid() {
			return nil, fmt.Errorf("invalid user '%s'", user.Email)
		}
	}
	if !sender.isValid() {
		return nil, fmt.Errorf("invalid user '%s'", sender.Email)
	}

	// Load the images into a [][]byte
	var images [][]byte
	for _, imagePath := range imagePaths {
		image, err := ioutil.ReadFile(imagePath)
		if err != nil {
			return nil, err
		}
		images = append(images, image)
	}

	timestamp := time.Now().String()
	return &Post{
		PostID:     crypto.Sha3String(sender.Email + timestamp),
		Timestamp:  timestamp,
		Sender:     sender,
		Recipients: recipients,
		Message:    message,
		Images:     images,
	}, nil
}

// String marshals auser to a string.
func (user *User) String() string {
	json, _ := json.MarshalIndent(*user, " ", "  ")
	return string(json)
}

// String marshals a post to a string.
func (post *Post) String() string {
	json, _ := json.MarshalIndent(*post, " ", "  ")
	return string(json)
}
