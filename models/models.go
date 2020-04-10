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

var (
	errInvalidEmail = errors.New("invalid email address")
)

type image []byte

// Grade is a grade enum.
type Grade int

const (
	// Freshman represents a freshman.
	Freshman = iota
	// Sophomore represents a sophomore.
	Sophomore = iota
	// Junior represents a junior.
	Junior = iota
	// Senior represents a senior.
	Senior = iota
)

// User represents a user.
type User struct {
	ID           int32  `pg:",pk"`
	UserID       string `pg:",notnull,unique"`
	Firstname    string `pg:",notnull"`
	Lastname     string `pg:",notnull"`
	Username     string `pg:",notnull"`
	Email        string `pg:",notnull"`
	Grade        Grade
	RegisterDate time.Time `pg:",notnull"`
}

// Post represents a post in the database.
type Post struct {
	ID         int32     `pg:",pk"`
	PostID     string    `pg:",notnull,unique"`
	Timestamp  time.Time `pg:",notnull"`
	Sender     User      `pg:",notnull"`
	Recipients []*User   `pg:",notnull"`

	Message string  `pg:",notnull"`
	Images  []image `pg:",notnull,array"`
}

// NewUser creates a *User given a valid email and grade.
func NewUser(email string, grade Grade) (*User, error) {
	// Check that the last part of the email is the correct email suffix
	if email[len(email)-len(common.EmailSuffix):] == common.EmailSuffix {
		userData := strings.Split(email, ".")

		// Check that there is only one first name and one last name
		if len(userData) != 3 {
			return nil, errInvalidEmail
		}

		return &User{
			UserID:       calcUserID(email),
			Firstname:    userData[0],
			Lastname:     strings.Split(userData[1], "@")[0],
			Username:     email[0 : len(email)-len(common.EmailSuffix)],
			Email:        email,
			Grade:        grade,
			RegisterDate: time.Now(),
		}, nil
	}
	return nil, errInvalidEmail
}

// calcUserID calculates the user ID.
func calcUserID(email string) string {
	return crypto.Sha3String(email)
}

// isValidUser checks if a User struct is valid
func (user *User) isValid() bool {
	l := len(user.Email) - len(common.EmailSuffix)
	if user.Email[l:] == common.EmailSuffix &&
		len(strings.Split(user.Email, ".")) == 3 {
		return true
	}
	return false
}

// NewPost creates a new post.
func NewPost(
	sender *User,
	message string,
	imagePaths []string,
	recipients ...*User,
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
	var images []image
	for _, imagePath := range imagePaths {
		tImage, err := ioutil.ReadFile(imagePath)
		if err != nil {
			return nil, err
		}
		images = append(images, image(tImage))
	}

	timestamp := time.Now()
	return &Post{
		PostID: crypto.Sha3String(
			sender.Email + timestamp.String() + message,
		),
		Timestamp:  timestamp,
		Sender:     *sender,
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
