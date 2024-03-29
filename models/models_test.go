package models

import "testing"

func TestNewUser(t *testing.T) {
	user, err := NewUser("first.last@mastersny.org", Freshman, false)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(user.String())
}

func TestNewPost(t *testing.T) {
	post, err := NewPost(
		"sen.der",
		"Hi, this is a test message!",
		[]string{"dW5pcXVlIGltYWdlIGRhdGEgb25l", "dW5pcXVlIGltYWdlIGRhdGEgdHdv"},
		[]string{"recip.one", "recip.two"},
	)

	if err != nil {
		t.Fatal(err)
	}

	t.Log(post)
}

func TestUserFromString(t *testing.T) {
	user, err := NewUser("first.last@mastersny.org", Freshman, false)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(user.String())

	stringedUser, err := UserFromString(user.String())
	if err != nil {
		t.Fatal(err)
	}
	t.Log(stringedUser.String())
}
