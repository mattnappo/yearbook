package common

import (
	"os"

	"github.com/juju/loggo"
	"github.com/juju/loggo/loggocolor"
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

	// DatabaseName is the name of the Postgres database.
	DatabaseName = "new_tests"

	// PasswordFile is the location of the file containing the Postgres
	// password.
	PasswordFile = "../password.pwd"

	// DefaultAPIRoot is the default API root.
	DefaultAPIRoot = "/api"
)

// NewLogger will create a new default loggo.Logger.
func NewLogger(context string) *loggo.Logger {
	// Create and setup a new logger
	logger := loggo.GetLogger(context)
	logger.SetLogLevel(loggo.DEBUG)
	loggo.ReplaceDefaultWriter(loggocolor.NewWriter(os.Stderr)) // Add colors

	return &logger
}
