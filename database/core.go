package database

import "github.com/xoreo/yearbook/models"

// AddPost adds a post to the database.
func (db *Database) AddPost(post *models.Post) error {
	db.mux.Lock()
	defer db.mux.Unlock()
	err := db.DB.Insert(post)
	if err != nil {
		return err
	}
	return nil
}

// GetPost gets a post from the database.
func (db *Database) GetPost(postID string) (models.Post, error) {
	post := &models.Post{PostID: postID}
	err := db.DB.Select(post)
	if err != nil {
		return models.Post{}, err
	}

	return *post, nil
}

// GetAllPosts gets all posts from the database.
func (db *Database) GetAllPosts() ([]models.Post, error) {
	var posts []models.Post
	err := db.DB.Model(&posts).Select()
	if err != nil {
		return nil, err
	}
	return posts, nil
}

// GetNPosts gets n posts from the database.
func (db *Database) GetNPosts(n int) error { return nil }

// DeletePost deletes a post from the database
func (db *Database) DeletePost(postID string) error {
	db.mux.Lock()
	defer db.mux.Unlock()
	return db.DB.Delete(&models.Post{PostID: postID})
}

// AddUser adds a new user to the database.
func (db *Database) AddUser(email string, grade models.Grade) error {
	db.mux.Lock()
	defer db.mux.Unlock()

	user, err := models.NewUser(email, grade)
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
func (db *Database) GetAllUsers() ([]models.User, error) {
	var users []models.User
	err := db.DB.Model(&users).Select()
	if err != nil {
		return nil, err
	}
	return users, nil
}

// DeleteUser deletes a user from the database
func (db *Database) DeleteUser(userID string) error {
	db.mux.Lock()
	defer db.mux.Unlock()
	return db.DB.Delete(&models.User{UserID: userID})
}
