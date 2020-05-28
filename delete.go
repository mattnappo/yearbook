package main

import (
	"os"

	"github.com/xoreo/yearbook/database"
)

func main() {
	db := database.Connect(false)
	defer db.Disconnect()

	if len(os.Args) != 2 {
		panic("no postid to delete")
	}

	uid := os.Args[1]
	err := db.DeletePost(uid)
	if err != nil {
		panic(err)
	}

}
