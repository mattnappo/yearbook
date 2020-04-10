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
	} else {
		g = models.Freshman
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
	if err != nil {
		t.Fatal(err)
	}
	err = db.AddPost(post)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetPost(t *testing.T) {
	pid := "813ccc6451537315143363e69c9cf936445c2cb80a96fa0a06e0649e57877272"
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

func TestGetNPosts(t *testing.T) {
	db := Connect(false)
	defer db.Disconnect()

	posts, err := db.GetNPosts(5)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(posts)
}

func TestDeletePost(t *testing.T) {
	db := Connect(false)
	defer db.Disconnect()

	uid := "813ccc6451537315143363e69c9cf936445c2cb80a96fa0a06e0649e57877272"
	err := db.DeletePost(uid)
	if err != nil {
		t.Fatal(err)
	}
}

func TestAddUser(t *testing.T) {
	db := Connect(false)
	defer db.Disconnect()

	u := genRandUser()
	t.Log(u)

	err := db.AddUser(u.Email, models.Freshman)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetUser(t *testing.T) {
	db := Connect(false)
	defer db.Disconnect()

	uid := "5999d445bb37c4282cd7c75fce5d954fe83d060d6a4c9d70b51cb00f89e3edb6"
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

	uid := "5999d445bb37c4282cd7c75fce5d954fe83d060d6a4c9d70b51cb00f89e3edb6"
	err := db.DeleteUser(uid)
	if err != nil {
		t.Fatal(err)
	}
}
