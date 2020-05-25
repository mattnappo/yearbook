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
	apiPort          = flag.Int64("start-api", common.APIPort, "start the API server on a given port")
)

func main() {
	flag.Parse()

	if *createSchemaFlag {
		db := database.Connect(false)
		defer db.Disconnect()
		db.CreateSchema()
		fmt.Println("created schema")
	}

	if *apiPort > 0 {
		err := api.StartAPIServer(*apiPort)
		if err != nil {
			panic(err)
		}
	}
}
