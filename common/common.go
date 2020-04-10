package common

import (
	"os"
	"path/filepath"
)

const (
	// EmailSuffix is the accepted email suffix.
	EmailSuffix = "@mastersny.org"

	// MaxRecipients is the maximum amount of recipients on one post.
	MaxRecipients = 10

	// MaxImages is the maximum amount of images on one post.
	MaxImages = 5

	// MaxMessageLength is the maximum amount of characters in a post message.
	MaxMessageLength = 2000

	// MaxEmailLength is the maximum amount of characters in an email.
	MaxEmailLength = 255

	// DatabaseName is the name of the Postgres database.
	DatabaseName = "new_tests"

	// PasswordFile is the location of the file containing the Postgres
	// password.
	PasswordFile = "./password.pwd"

	// DefaultAPIRoot is the default API root.
	DefaultAPIRoot = "/api"

	// LogsDir is the location where all log files are stored.
	LogsDir = "./data/logs"
)

// CreateDirIfDoesNotExist creates a directory if it does not already exist.
func CreateDirIfDoesNotExist(dir string) error {
	dir = filepath.FromSlash(dir)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			return err
		}
	}
	return nil
}
