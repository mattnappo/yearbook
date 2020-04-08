package main

import (
	"fmt"

	"github.com/go-pg/pg/v9"
	"github.com/go-pg/pg/v9/orm"
	"golang.org/x/crypto/ssh/terminal"
)

type User struct {
	ID           int64  `sql:"notnull"`
	FirstName    string `sql:"notnull"`
	LastName     string `sql:"notnull"`
	Username     string `sql:"notnull"`
	Email        string `sql:"notnull"`
	RegisterDate string `sql:"notnull"`
}

type Post struct {
	PostID     PostID `sql:"notnull"`
	Timestamp  string `sql:"notnull"`
	Sender     User   `sql:"notnull"`
	Recipients []User `sql:"notnull"`

	Message string `sql:"notnull"`
	Images  [][]byte
}

func Connect() *pg.DB {
	fmt.Printf("Password: ")
	pwd, _ := terminal.ReadPassword(0)
	fmt.Printf("\n")
	db := pg.Connect(&pg.Options{
		User:     "postgres",
		Password: string(pwd),
	})
	return db
}

func Model() error {
	db := Connect()
	defer db.Close()

	err := createSchema(db)
	if err != nil {
		return err
	}

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
	Model()
}
