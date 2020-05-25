package common

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"unicode/utf8"

	"github.com/joho/godotenv"
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

	// LogsDir is the location where all log files are stored.
	LogsDir = "./data/logs"

	// APIPort represents the default api server port
	APIPort = 8081

	// envFile is the path to the file containing needed environment variables.
	envFile = "./.env"

	// Email notification info
	NotifEmail    = "mastersseniors2020.com@gmail.com"
	NotifProvider = "smtp.gmail.com"
	NotifPort     = 587
)

var (
	// DatabaseName is the name of the Postgres database.
	DatabaseName = GetEnv("DATABASE_NAME")

	// NotifPassword is the password for the gmail account that sends notifications.
	NotifPassword = GetEnv("NOTIF_PASSWORD")

	// NotifsEnabled turns email notifications on or off.
	NotifsEnabled = false
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

// GetEnv gets an environment variable
func GetEnv(key string) string {
	// Load .env file every time
	err := godotenv.Load(envFile)

	if err != nil {
		panic(errors.New("could not load .env file"))
	}

	return os.Getenv(key)
}

// StringToArray returns an array given a string.
func StringToArray(s string) []string {
	s = strings.TrimSuffix(s, "]")
	s = trimFirstRune(s)
	s = strings.ReplaceAll(s, "\"", "")
	return strings.Split(s, ", ")
}

func trimFirstRune(s string) string {
	_, i := utf8.DecodeRuneInString(s)
	return s[i:]
}
