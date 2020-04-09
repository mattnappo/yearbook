package database

import "github.com/xoreo/yearbook/models"

// AddUser adds a new user to the database.
func (db *Database) AddUser(email string) error {
	db.mux.Lock()
	defer db.mux.Unlock()

	user, err := models.NewUserFromEmail(email)
	if err != nil {
		return err
	}

	err = db.DB.Insert(user)
	if err != nil {
		return err
	}
	return nil
}

// GetUser gets a user from the database.
func (db *Database) GetUser(userID string) (models.User, error) {
	user := &models.User{UserID: userID}
	err := db.DB.Select(user)
	if err != nil {
		return models.User{}, err
	}

	return *user, nil
}

// GetAllUsers gets all users from the database.
func (db *Database) GetAllUsers() error { return nil }

// DeleteUser deletes a user from the database
func (db *Database) DeleteUser(userID string) error { return nil }
