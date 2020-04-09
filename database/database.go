package database

import (
	"errors"
	"fmt"
	"sync"

	"github.com/go-pg/pg/v9"
	"github.com/xoreo/yearbook/common"
	"github.com/xoreo/yearbook/models"
	"golang.org/x/crypto/ssh/terminal"
)

type connStatus int

const (
	CONNECTED    = iota
	DISCONNECTED = iota
	ERROR        = iota
)

var (
	errDisconnected = errors.New("not connected to database")
)

// Database represents a database.
type Database struct {
	DB     *pg.DB
	status connStatus
	mux    sync.Mutex
}

// Connect connects to the database.
func Connect() *Database {
	fmt.Printf("Password: ")
	pwd, _ := terminal.ReadPassword(0)
	fmt.Printf("\n")
	db := pg.Connect(&pg.Options{
		User:     "postgres",
		Password: string(pwd),
		Database: common.DatabaseName,
	})
	return &Database{db, CONNECTED, sync.Mutex{}}
}

// Disconnect disconnects from the database.
func (db *Database) Disconnect() error {
	err := db.DB.Close()
	if err != nil {
		db.status = ERROR
		return err
	}
	db.status = DISCONNECTED
	return nil
}

// AddPost adds a post to the database.
func (db *Database) AddPost(post *models.Post) error {
	if db.status != CONNECTED {
		return errDisconnected
	}
	db.mux.Lock()
	defer db.mux.Unlock()

	err := db.DB.Insert(post)
	if err != nil {
		return err
	}
	return nil
}

// AddUser adds a new user to the database.
func (db *Database) AddUser(email string) error { return nil }

// GetPost gets a post from the database.
func (db *Database) GetPost(postID string) error { return nil }

// GetUser gets a user from the database.
func (db *Database) GetUser(userID string) error { return nil }

// GetAllPosts gets all posts from the database.
func (db *Database) GetAllPosts() error { return nil }

// GetAllUsers gets all users from the database.
func (db *Database) GetAllUsers() error { return nil }

// GetNPosts gets n posts from the database.
func (db *Database) GetNPosts(n int) error { return nil }

// createSchema creates the database schema.
func (db *Database) createSchema() error {
	db.mux.Lock()
	defer db.mux.Unlock()
	for _, model := range []interface{}{
		(*models.User)(nil),
		(*models.Post)(nil)} {
		err := db.DB.CreateTable(model, nil)
		if err != nil {
			return err
		}
	}
	return nil
}
