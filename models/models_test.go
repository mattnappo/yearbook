package models

import "testing"

func TestNewUserFromEmail(t *testing.T) {
	user, err := NewUserFromEmail("first.last@mastersny.org")
	if err != nil {
		t.Fatal(err)
	}

	t.Log(user.String())
}
