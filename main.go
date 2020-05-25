package main

import (
	"flag"
	"fmt"

	"github.com/xoreo/yearbook/api"
	"github.com/xoreo/yearbook/common"
	"github.com/xoreo/yearbook/database"
)

var (
	createSchemaFlag = flag.Bool("create-schema", false, "create the database schema")
	addSeniorsFlag   = flag.Bool("add-seniors", false, "add the seniors to the database")
	apiPort          = flag.Int64("start-api", common.APIPort, "start the API server on a given port")
)

func main() {
	flag.Parse()

	if *createSchemaFlag {
		db := database.Connect(false)
		defer db.Disconnect()
		err := db.CreateSchema()
		if err != nil {
			panic(err)
		}
		fmt.Println("created schema")
	}

	if *addSeniorsFlag {
		db := database.Connect(false)
		defer db.Disconnect()
		err := db.AddSeniors()
		if err != nil {
			panic(err)
		}
		fmt.Println("added the seniors to the database")
	}

	if *apiPort > 0 {
		err := api.StartAPIServer(*apiPort)
		if err != nil {
			panic(err)
		}
	}
}
