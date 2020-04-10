package main

import (
	"flag"
	"fmt"

	"github.com/xoreo/yearbook/database"
)

var (
	createSchemaFlag = flag.Bool("create-schema", false, "create the database schema")
)

func main() {
	flag.Parse()

	if *createSchemaFlag {
		db := database.Connect(false)
		defer db.Disconnect()
		db.CreateSchema()
		fmt.Println("created schema")
	}
}
