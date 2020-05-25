package database

import (
	"errors"
	"sync"

	"github.com/go-pg/pg/v9"
	"github.com/xoreo/yearbook/common"
	"github.com/xoreo/yearbook/models"
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
func Connect(fromStdin bool) *Database {
	db := pg.Connect(&pg.Options{
		User:     "postgres",
		Password: common.GetEnv("DB_PASSWORD"),
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

// CreateSchema creates the database schema.
func (db *Database) CreateSchema() error {
	db.mux.Lock()
	defer db.mux.Unlock()
	for _, model := range []interface{}{
		(*models.User)(nil), // Make the users table
		(*models.Post)(nil), // make the posts table
		(*token)(nil)} {     // make the tokens table
		err := db.DB.CreateTable(model, nil)
		if err != nil {
			return err
		}
	}
	return nil
}
