package models

import "testing"

func TestNewUser(t *testing.T) {
	user, err := NewUser("first.last@mastersny.org", Freshman)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(user.String())
}

func TestIsValidUser(t *testing.T) {
	vUser, err := NewUser("first.last@mastersny.org", Freshman)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(vUser.isValid())
}

func TestNewPost(t *testing.T) {
	sender, err := NewUser("sen.der@mastersny.org", Freshman)
	recip1, err := NewUser("recip.one@mastersny.org", Senior)
	recip2, err := NewUser("recip.two@mastersny.org", Senior)
	if err != nil {
		t.Fatal(err)
	}

	post, err := NewPost(
		sender,
		"Hi, this is a test message!",
		[]string{"../data/img1.jpg", "../data/img2.jpg"},
		recip1, recip2,
	)

	if err != nil {
		t.Fatal(err)
	}

	t.Log(post)
}
