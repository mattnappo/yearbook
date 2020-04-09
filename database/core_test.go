package database

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/xoreo/yearbook/models"
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

func TestAddPost(t *testing.T) {
	db := Connect(false)
	defer db.Disconnect()
	post, err := models.NewPost(
		genRandUser(),
		"I am a message",
		[]string{"../data/img1.jpg", "../data/img2.jpg"},
		genRandUser(),
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
