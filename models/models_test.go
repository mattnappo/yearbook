package models

import "testing"

func TestNewUserFromEmail(t *testing.T) {
	user, err := NewUserFromEmail("first.last@mastersny.org")
	if err != nil {
		t.Fatal(err)
	}

	t.Log(user.String())
}

func TestIsValidUser(t *testing.T) {
	vUser, err := NewUserFromEmail("first.last@mastersny.org")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(vUser.isValid())

	// ivUser, _ := NewUserFromEmail("first.last@gmail.com")
	// t.Log(ivUser.isValid())
}

func TestNewPost(t *testing.T) {
	sender, err := NewUserFromEmail("sen.der@mastersny.org")
	recip1, err := NewUserFromEmail("recip.one@mastersny.org")
	recip2, err := NewUserFromEmail("recip.two@mastersny.org")
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
