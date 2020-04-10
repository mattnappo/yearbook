package database

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/xoreo/yearbook/models"
)

func genRandUser() string {
	r := fmt.Sprintf("%f", rand.Float64())[2:]
	return fmt.Sprintf("first%s.last%s", r, r)
}

func genRandEmail() string {
	r := fmt.Sprintf("%f", rand.Float64())[2:]
	return fmt.Sprintf("first%s.last%s@mastersny.org", r, r)
}

func TestAddPost(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	db := Connect(false)
	defer db.Disconnect()

	post, err := models.NewPost(
		genRandUser(),
		"I am a message",
		[]string{"dW5pcXVlIGltYWdlIGRhdGEgb25l", "dW5pcXVlIGltYWdlIGRhdGEgdHdv"},
		[]string{genRandUser()},
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
	pid := "c596d281587fa801d1329558f39f9792a9a96ca9bf75f18296a8172b46393403"
	db := Connect(false)
	defer db.Disconnect()

	post, err := db.GetPost(pid)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(post)
}

func TestGetAllPosts(t *testing.T) {
	db := Connect(false)
	defer db.Disconnect()

	posts, err := db.GetAllPosts()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(posts)
}

func TestGetnPosts(t *testing.T) {
	db := Connect(false)
	defer db.Disconnect()

	posts, err := db.GetnPosts(5)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(posts)
}

func TestDeletePost(t *testing.T) {
	db := Connect(false)
	defer db.Disconnect()

	uid := "c596d281587fa801d1329558f39f9792a9a96ca9bf75f18296a8172b46393403"
	err := db.DeletePost(uid)
	if err != nil {
		t.Fatal(err)
	}
}

func TestAddUser(t *testing.T) {
	db := Connect(false)
	defer db.Disconnect()

	user, err := models.NewUser("cool.dude@mastersny.org", models.Sophomore)
	if err != nil {
		t.Fatal(err)
	}

	err = db.AddUser(user)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetUser(t *testing.T) {
	db := Connect(false)
	defer db.Disconnect()

	uid := "first252744.last252744"
	u, err := db.GetUser(uid)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(u)
}

func TestGetAllUsers(t *testing.T) {
	db := Connect(false)
	defer db.Disconnect()

	users, err := db.GetAllUsers()
	if err != nil {
		t.Fatal(err)
	}

	t.Log(users)
}

func TestDeleteUser(t *testing.T) {
	db := Connect(false)
	defer db.Disconnect()

	uid := "first252744.last252744"
	err := db.DeleteUser(uid)
	if err != nil {
		t.Fatal(err)
	}
}
