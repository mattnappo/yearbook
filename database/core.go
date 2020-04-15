package database

import (
	"github.com/xoreo/yearbook/models"
)

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
	post := &models.Post{}
	err := db.DB.Model(post).
		Where("post.post_id = ?", postID).
		Select()
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

// GetnPosts gets n posts from the database.
func (db *Database) GetnPosts(n int) ([]models.Post, error) {
	var posts []models.Post
	err := db.DB.Model(&posts).Order("id ASC").Limit(n).Select()
	if err != nil {
		return nil, err
	}
	return posts, nil
}

// DeletePost deletes a post from the database
func (db *Database) DeletePost(postID string) error {
	db.mux.Lock()
	defer db.mux.Unlock()

	_, err := db.DB.Model(&models.Post{}).
		Where("post.post_id = ?", postID).
		Delete()
	return err
}

// AddUser adds a new user to the database.
func (db *Database) AddUser(user *models.User) error {
	db.mux.Lock()
	defer db.mux.Unlock()
	err := db.DB.Insert(user)
	if err != nil {
		return err
	}
	return nil
}

// UpdateUser updates a user with the given new values in a user struct.
func (db *Database) UpdateUser(user *models.User) error {
	// HORRIBLE FUNC RIGHT HERE FIX THIS BOI
	// fmt.Printf("\n\n%v\n\n", user)
	db.mux.Lock()
	defer db.mux.Unlock()

	lookupUser, err := db.GetUser(string(user.Username))
	if err != nil {
		return err
	}

	if user.Bio != "" {
		lookupUser.Bio = user.Bio
	}
	if user.Will != "" {
		lookupUser.Will = user.Will
	}
	if user.ProfilePic != nil {
		lookupUser.ProfilePic = user.ProfilePic
	}
	if user.Nickname != "" {
		lookupUser.Nickname = user.Nickname
	}

	err = db.DB.Update(&lookupUser)
	if err != nil {
		return err
	}

	return nil
}

// GetUser gets a user from the database.
func (db *Database) GetUser(username string) (models.User, error) {
	user := &models.User{}
	err := db.DB.Model(user).
		Where("username = ?", username).
		Select()
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
func (db *Database) DeleteUser(username string) error {
	db.mux.Lock()
	defer db.mux.Unlock()
	_, err := db.DB.Model(&models.User{}).
		Where("username = ?", username).
		Delete()
	return err
}
