package database

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/xoreo/yearbook/models"
)

func genRandUser(grade ...models.Grade) *models.User {
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
		"I am a message",
		[]string{"../data/img1.jpg", "../data/img2.jpg"},
		genRandUser(models.Senior),
	)
	if err != nil {
		t.Fatal(err)
	}
	err = db.AddPost(post)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetPost(t *testing.T) {
	pid := "522f46a76a8ed12824c10559d01a436994a3f5147d99a73e96f88245952ecc93"
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
