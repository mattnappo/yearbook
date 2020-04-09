package main

import "flag"

import "github.com/xoreo/yearbook/database"

var (
	createSchemaFlag = flag.Bool("create-schema", false, "create the database schema")
)

func main() {
	flag.Parse()

	if *createSchemaFlag {
		db := database.Connect()
		defer db.Disconnect()
		db.CreateSchema()
	}
}
