package main

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/go-pg/pg/v9"
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

func Test(shouldCreateSchema bool) error {
	db := Connect()
	defer db.Close()

	if shouldCreateSchema {
		err := createSchema(db)
		if err != nil {
			return err
		}
	}

	user1, err := NewUserFromEmail("first1.last1@mastersny.org")
	if err != nil {
		return err
	}
	user2, err := NewUserFromEmail("first2.last2@mastersny.org")
	user3, err := NewUserFromEmail("first3.last3@mastersny.org")
	testPost := &Post{
		PostID:     PostID("some_post_id"),
		Timestamp:  "I'm the time!",
		Sender:     user1,
		Recipients: []*User{user2, user3},
		Message:    "message message message text",
		Images:     [][]byte{[]byte("img1"), []byte("img2")},
	}
	err = db.Insert(testPost)
	if err != nil {
		return err
	}
	return nil
}

// createSchema creates the database schema.
func createSchema(db *pg.DB) error {
	for _, model := range []interface{}{(*User)(nil), (*Post)(nil)} {
		err := db.CreateTable(model, nil)
		if err != nil {
			return err
		}
	}
	return nil
}

func main() {
	err := Test(true)
	if err != nil {
		panic(err)
	}
}
