package database

import (
	"fmt"

	"github.com/go-pg/pg/v9"
	"github.com/xoreo/yearbook/models"
	"golang.org/x/crypto/ssh/terminal"
)

// Connect connects to the database.
func Connect() *pg.DB {
	fmt.Printf("Password: ")
	pwd, _ := terminal.ReadPassword(0)
	fmt.Printf("\n")
	db := pg.Connect(&pg.Options{
		User:     "postgres",
		Password: string(pwd),
		Database: "gotests",
	})
	return db
}

// createSchema creates the database schema.
func createSchema(db *pg.DB) error {
	for _, model := range []interface{}{(*models.User)(nil), (*models.Post)(nil)} {
		err := db.CreateTable(model, nil)
		if err != nil {
			return err
		}
	}
	return nil
}
