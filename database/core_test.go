package database

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/xoreo/yearbook/models"
)

func genRandUser(grade ...models.Grade) *models.User {
	rand.Seed(time.Now().UTC().UnixNano())
	var g models.Grade
	if len(grade) == 1 {
		g = grade[0]
	}

	r := rand.Intn(999) + 1
	u, err := models.NewUser(
		fmt.Sprintf("first%d.last%d@mastersny.org", r, r),
		g,
	)
	if err != nil {
		panic(err) // Panic because panic is better here
	}
	return u
}

func TestAddPost(t *testing.T) {
	db := Connect(false)
	defer db.Disconnect()
	post, err := models.NewPost(
		genRandUser(),
		"I am a messagess",
		[]string{"../data/img1.jpg", "../data/img2.jpg"},
		genRandUser(models.Senior),
	)
	fmt.Println(post)
	if err != nil {
		t.Fatal(err)
	}
	err = db.AddPost(post)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetPost(t *testing.T) {
	pid := "f4a1414e255e65d52d124a685f441566900a96e220d7d691b8654669d9008107"
	db := Connect(false)
	defer db.Disconnect()

	post, err := db.GetPost(pid)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(post)
}

func TestGetAllPosts(t *testing.T) {

}

func TestGetNPosts(t *testing.T) {

}

func TestDeletePost(t *testing.T) {

}

func TestAddUser(t *testing.T) {

}

func TestGetUser(t *testing.T) {

}

func TestGetAllUsers(t *testing.T) {

}

func TestDeleteUser(t *testing.T) {

}
