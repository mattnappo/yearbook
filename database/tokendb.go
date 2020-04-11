package database

import "golang.org/x/oauth2"

// token describes the schema for the token table.
type token struct {
	sub   int64  `pg:",pk"`
	token string `pg:",unique"`
	email string `pg:",unique"`
}

// InsertToken inserts a token into the database.
func (db *Database) InsertToken(
	sub int64, token *oauth2.Token, email ...string,
) error {
	return nil
}

// GetToken gets a token in the token table.
func (db *Database) GetToken(sub int64) {

}
