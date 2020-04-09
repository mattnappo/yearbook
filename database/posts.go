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
	return db.DB.Delete(&models.Post{PostID: postID})
}
