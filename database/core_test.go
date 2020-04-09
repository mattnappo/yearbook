package database

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/xoreo/yearbook/models"
)

func genRandUser(grade ...models.Grade) *models.User {
	var g models.Grade
	if len(grade) == 1 {
		g = grade[0]
	}

	r := fmt.Sprintf("%f", rand.Float64())[2:]
	u, err := models.NewUser(
		fmt.Sprintf("first%s.last%s@mastersny.org", r, r),
		g,
	)
	if err != nil {
		panic(err) // Panic because panic is better here
	}
	return u
}

func TestAddPost(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	db := Connect(false)
	defer db.Disconnect()
	post, err := models.NewPost(
		genRandUser(),
		"I am a messagessadsd",
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
	pid := "f4f470a193e2b94b2a22b4484693692dc8c4a9cbc3a7c06a1aa866264b71362e"
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
