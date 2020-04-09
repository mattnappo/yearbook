package main

import (
	"errors"
	"fmt"
	"github.com/go-pg/pg/v9"
	"github.com/go-pg/pg/v9/orm"
	"golang.org/x/crypto/ssh/terminal"
	"strings"
	"time"
)

const emailPrefix = "@mastersny.org"

type PostID string

type User struct {
	ID           string `pg:",notnull"`
	Firstname    string `pg:",notnull"`
	Lastname     string `pg:",notnull"`
	Username     string `pg:",notnull"`
	Email        string `pg:",notnull"`
	RegisterDate string `pg:",notnull"`
}

func NewUserFromEmail(email string) (*User, error) {
	if email[len(email)-len(emailPrefix):] == emailPrefix {
		userData := strings.Split(email, ".")
		return &User{
			ID:           "hash of email",
			Firstname:    userData[0],
			Lastname:     userData[1],
			Username:     email[0 : len(email)-len(emailPrefix)],
			Email:        email,
			RegisterDate: time.Now().String(),
		}, nil
	}
	return nil, errors.New("invalid email address")
}

type Post struct {
	PostID     PostID  `pg:",notnull"`
	Timestamp  string  `pg:",notnull"`
	Sender     *User   `pg:",notnull"`
	Recipients []*User `pg:",notnull"`

	Message string   `pg:",notnull"`
	Images  [][]byte `pg:",notnull"`
}

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

func createSchema(db *pg.DB) error {
	for _, model := range []interface{}{(*User)(nil), (*Post)(nil)} {
		err := db.CreateTable(model, &orm.CreateTableOptions{
			Temp: false,
		})
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
