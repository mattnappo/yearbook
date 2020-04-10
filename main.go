package main

import (
	"flag"
	"fmt"

	"github.com/xoreo/yearbook/database"
	"github.com/xoreo/yearbook/models"
)

var (
	createSchemaFlag = flag.Bool("create-schema", false, "create the database schema")
)

func testNewPost() {
	sender, err := models.NewUser("sen.der@mastersny.org", models.Freshman)
	recip1, err := models.NewUser("recip.one@mastersny.org", models.Senior)
	recip2, err := models.NewUser("recip.two@mastersny.org", models.Senior)
	if err != nil {
		panic(err)
	}

	post, err := models.NewPost(
		sender,
		"Hi, this is a test message!",
		[]string{"./data/img1.jpg", "./data/img2.jpg"},
		recip1, recip2,
	)

	if err != nil {
		panic(err)
	}

	db := database.Connect(false)
	defer db.Disconnect()

	err = db.AddPost(post)
	if err != nil {
		panic(err)
	}

	fmt.Println(post)
}

func main() {
	flag.Parse()

	if *createSchemaFlag {
		db := database.Connect(true)
		defer db.Disconnect()
		db.CreateSchema()
		fmt.Println("created schema")
	}
	// for i := 0; i < 10; i++ {
	testNewPost()
	// time.Sleep(1 * time.Second)
	// }
}
