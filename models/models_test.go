package models

import "testing"

func TestNewUser(t *testing.T) {
	user, err := NewUser("first.last@mastersny.org", Freshman)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(user.String())
}

func TestNewPost(t *testing.T) {
	post, err := NewPost(
		"sen.der",
		"Hi, this is a test message!",
		[][]byte{[]byte("image data one"), []byte("image data two")},
		[]string{"recip.one", "recip.two"},
	)

	if err != nil {
		t.Fatal(err)
	}

	t.Log(post)
}
