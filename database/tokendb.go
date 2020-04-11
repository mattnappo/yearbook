package database

import "golang.org/x/oauth2"

// token describes the schema for the token table.
type token struct {
	sub   float64 `pg:",pk"`
	token string  `pg:",unique"`
	email string  `pg:",unique"`
}

// InsertToken inserts a token into the database.
func (db *Database) InsertToken(
	sub float64, oauthToken *oauth2.Token, email ...string,
) error {
	db.mux.Lock()
	defer db.mux.Unlock()
	// Construct a new *token
	token := &token{
		sub:   sub,
		token: oauthToken.AccessToken,
	}
	if len(email) > 0 {
		token.email = email[0]
	}

	// Put it in the database
	err := db.DB.Insert(token)
	if err != nil {
		return err
	}
	return nil
}

// GetToken gets a token in the token table.
func (db *Database) GetToken(sub int64) (string, error) {
	return "", nil
}
