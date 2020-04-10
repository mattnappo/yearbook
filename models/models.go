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
	errInvalidEmail = errors.New("malformed email address")
)

// Username represents a username.
type Username string

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
	ID       int32    `pg:",pk"`
	Username Username `pg:",notnull,unique"`

	Firstname    string `pg:",notnull"`
	Lastname     string `pg:",notnull"`
	Email        string `pg:",notnull,unique"`
	Grade        Grade
	RegisterDate time.Time `pg:",notnull"`
}

// Post represents a post in the database.
type Post struct {
	ID     int32  `pg:",pk"`
	PostID string `pg:",notnull,unique"`

	Timestamp  time.Time `pg:",notnull"`
	Sender     User      `pg:",notnull"`
	Recipients []*User   `pg:",notnull"`

	Message string  `pg:",notnull"`
	Images  []image `pg:",notnull,array"`
}

// NewUser creates a *User given a valid email and grade.
func NewUser(email string, grade Grade) (*User, error) {
	userData := strings.Split(email, ".")

	// Check that there is only one first name and one last name
	if len(userData) != 3 {
		return nil, errInvalidEmail
	}

	username, err := UsernameFromEmail(email)
	if err != nil {
		return nil, err
	}

	return &User{
		Username:     username,
		Firstname:    username.firstname(),
		Lastname:     username.lastname(),
		Email:        username.email(),
		Grade:        grade,
		RegisterDate: time.Now(),
	}, nil
}

// NewPost creates a new post.
func NewPost(
	sender UserID,
	message string,
	imagePaths [][]byte,
	recipients []string,
) (*Post, error) {
	// Check that all data for the post is valid
	if sender == "" || len(recipients) > common.MaxRecipients ||
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

// UsernameFromEmail constructs a username given an email.
func UsernameFromEmail(email string) (Username, error) {
	if len(email) > common.MaxEmailLength {
		return errors.New("email is too long")
	}

	// Check that the email suffix is at the end
	if email[len(email)-len(common.EmailSuffix):] != common.EmailSuffix {
		return errInvalidEmail
	}
	return email[0 : len(email)-len(common.EmailSuffix)]
}

func (u Username) firstname() string {
	components := strings.Split(string(u), ".")
	return components[0]
}

func (u Username) lastname() string {
	components := strings.Split(string(u), ".")
	return components[0]
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
