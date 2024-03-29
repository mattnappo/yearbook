package models

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/mattnappo/yearbook/common"
	"github.com/mattnappo/yearbook/crypto"
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
	ID       int32    `pg:",pk" json:"id"`
	Username Username `pg:",pk" json:"username"`

	Firstname    string    `pg:",notnull" json:"firstname"`
	Lastname     string    `pg:",notnull" json:"lastname"`
	Email        string    `pg:",notnull,unique" json:"email"`
	Grade        Grade     `pg:",use_zero" json:"grade"`
	RegisterDate time.Time `pg:",notnull" json:"register_date"`

	// Mutable fields
	Nickname   string `json:"nickname"`
	ProfilePic string `json:"profile_pic"` // A url
	Bio        string `json:"bio"`
	Will       string `json:"will"`
	Registered bool   `json:"registered"`

	OutboundPosts []string `json:"outbound_posts"` // postIDs from this user
	InboundPosts  []string `json:"inbound_posts"`  // postIDs to this user
}

// Post represents a post in the database.
type Post struct {
	ID     int32  `pg:",pk" json:"id"`
	PostID string `pg:",notnull,unique" json:"post_id"`

	Timestamp  time.Time  `pg:",notnull" json:"timestamp"`
	Sender     Username   `pg:",notnull" json:"sender"`
	Recipients []Username `pg:",notnull" json:"recipients"`

	Message string  `pg:",notnull" json:"message"`
	Images  []image `pg:",array" json:"images"`
}

// NewUser creates a *User given a valid email and grade.
func NewUser(email string, grade Grade, registered bool) (*User, error) {
	username, err := UsernameFromEmail(email)
	if err != nil {
		return nil, err
	}

	return &User{
		Username:     username,
		Firstname:    username.Firstname(),
		Lastname:     username.Lastname(),
		Email:        username.Email(),
		Grade:        grade,
		Registered:   registered,
		RegisterDate: time.Now(),
	}, nil
}

// NewPost creates a new post.
func NewPost(
	senderUsername string,
	message string,
	images []string,
	recipientsUsernames []string,
) (*Post, error) {
	// Check that all data for the post is valid
	if len(recipientsUsernames) > common.MaxRecipients ||
		len(recipientsUsernames) <= 0 || message == "" ||
		len(message) > common.MaxMessageLength ||
		len(images) > common.MaxImages {
		return nil, errors.New("too much or not enough data to construct post")
	}

	// Validate sender username
	sender, err := validateUsername(senderUsername)
	if err != nil {
		return nil, err
	}
	// Validate all recipient usernames
	var recipients []Username
	for _, recipient := range recipientsUsernames {
		validRecipient, err := validateUsername(recipient)
		if err != nil {
			return nil, err
		}
		// Add valid recipient to the slice of valid recipients
		recipients = append(recipients, validRecipient)
	}

	// Cast all base64 []string images to images
	var byteImages []image
	for _, base64Image := range images {
		byteImage, err := base64.StdEncoding.DecodeString(base64Image)
		if err != nil {
			return nil, err
		}
		byteImages = append(byteImages, image(byteImage))
	}

	post := &Post{
		Sender:     sender,
		Recipients: recipients,
		Message:    message,
		Images:     byteImages,
	}
	post.PostID = crypto.Sha3String(post.String())
	post.Timestamp = time.Now()
	return post, nil

}

// UserFromString returns a new User given a JSON/string representation
// of a user struct.
func UserFromString(data string) (*User, error) {
	user := &User{}
	err := json.Unmarshal([]byte(data), user)
	if err != nil {
		return nil, err
	}

	return user, err
}

// UsernameFromEmail constructs a username given an email.
func UsernameFromEmail(email string) (Username, error) {
	// Check that the email is not too long
	if len(email) > common.MaxEmailLength || len(email) <=
		len(common.EmailSuffix) {
		return Username(""), fmt.Errorf("email %s is too long or too short", email)
	}

	// Check that the email suffix is at the end and that there is only one first
	// name and one last name.
	if email[len(email)-len(common.EmailSuffix):] != common.EmailSuffix ||
		len(strings.Split(email, ".")) != 3 {
		return Username(""), errInvalidEmail
	}

	// Return firstname.lastname as a Username type
	return Username(email[0 : len(email)-len(common.EmailSuffix)]), nil
}

// Firstname returns the first name associated with the username.
func (u Username) Firstname() string {
	components := strings.Split(string(u), ".")

	return components[0]
}

// Lastname returns the last name associated with the username.
func (u Username) Lastname() string {
	components := strings.Split(string(u), ".")

	return strings.Split(components[1], "@")[0]
}

// Name returns the properly formatted name associated with the
// username.
func (u Username) Name() string {
	return strings.Title(
		fmt.Sprintf("%s %s", u.Firstname(), u.Lastname()),
	)
}

// Email returns the email associated with the username.
func (u Username) Email() string {
	return string(u) + common.EmailSuffix
}

// isValid checks if a given username is valid.
func (u Username) isValid() bool {
	return len(strings.Split(string(u), ".")) == 2
}

// validateUsername will attempt to validate a username.
func validateUsername(u string) (Username, error) {
	if Username(u).isValid() {
		return Username(u), nil
	}
	return Username(""), fmt.Errorf("invalid username '%s", u)
}

// String marshals a user to a string.
func (user *User) String() string {
	json, _ := json.MarshalIndent(*user, " ", "  ")
	return string(json)
}

// String marshals a post to a string.
func (post *Post) String() string {
	json, _ := json.MarshalIndent(*post, " ", "  ")
	return string(json)
}
