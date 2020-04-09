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
	// CONNECTED means the connection to the database is established.
	CONNECTED = iota

	// DISCONNECTED means the connection to the database is not established.
	DISCONNECTED = iota

	// ERROR means the database is in some error state.
	ERROR = iota
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

// createSchema creates the database schema.
func (db *Database) createSchema() error {
	db.mux.Lock()
	defer db.mux.Unlock()
	for _, model := range []interface{}{
		(*models.User)(nil),   // Make the users table
		(*models.Post)(nil)} { // make the posts table
		err := db.DB.CreateTable(model, nil)
		if err != nil {
			return err
		}
	}
	return nil
}
