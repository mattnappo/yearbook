package main

import (
	"flag"
	"fmt"
	"math/rand"

	"github.com/xoreo/yearbook/database"
	"github.com/xoreo/yearbook/models"
)

var (
	createSchemaFlag = flag.Bool("create-schema", false, "create the database schema")
)

func genRandUser() *models.User {
	r := rand.Intn(999) + 1
	u, err := models.NewUserFromEmail(
		fmt.Sprintf("first%d.last%d@mastersny.org", r, r),
	)
	if err != nil {
		panic(err)
	}
	return u
}

func main() {
	flag.Parse()

	if *createSchemaFlag {
		db := database.Connect()
		defer db.Disconnect()
		db.CreateSchema()
	}

}
