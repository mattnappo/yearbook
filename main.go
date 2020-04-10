package main

import (
	"flag"
	"fmt"

	"github.com/xoreo/yearbook/api"
	"github.com/xoreo/yearbook/database"
)

var (
	createSchemaFlag = flag.Bool("create-schema", false, "create the database schema")
	startAPIFlag     = flag.Bool("start-api", false, "start the API server")
)

func main() {
	flag.Parse()

	if *createSchemaFlag {
		db := database.Connect(false)
		defer db.Disconnect()
		db.CreateSchema()
		fmt.Println("created schema")
	}

	if *startAPIFlag {
		err := api.StartAPIServer(8081)
		if err != nil {
			panic(err)
		}
	}
}
