package database

import (
	"golang.org/x/oauth2"
)

// token describes the schema for the token table.
type token struct {
	Sub   string `pg:",pk"`
	Token string `pg:",notnull,unique"`
	Email string `pg:",notnull,unique"`
}

// InsertToken inserts a token into the database.
func (db *Database) InsertToken(
	sub string, oauthToken *oauth2.Token, email ...string,
) error {
	db.mux.Lock()
	defer db.mux.Unlock()
	// Construct a new *token
	token := &token{
		Sub:   sub,
		Token: oauthToken.AccessToken,
	}
	if len(email) > 0 {
		token.Email = email[0]
	}

	// Insert if it is not there, update if it is.
	_, err := db.DB.Model(token).
		OnConflict("(sub) DO UPDATE").
		Set("token = EXCLUDED.token").
		Insert()
	if err != nil {
		return err
	}
	return nil
}

// GetToken gets a token in the token table.
func (db *Database) GetToken(sub string) (string, error) {
	token := &token{Sub: sub}  // Init the buffer
	err := db.DB.Select(token) // Select token from database
	if err != nil {
		return "", err
	}

	return token.Token, nil
}
