package database

import "github.com/xoreo/yearbook/models"

// AddPost adds a post to the database.
func (db *Database) AddPost(post *models.Post) error {
	if db.status != CONNECTED {
		return errDisconnected
	}
	db.mux.Lock()
	defer db.mux.Unlock()

	err := db.DB.Insert(post)
	if err != nil {
		return err
	}
	return nil
}

// GetPost gets a post from the database.
func (db *Database) GetPost(postID string) error { return nil }

// GetAllPosts gets all posts from the database.
func (db *Database) GetAllPosts() error { return nil }

// GetNPosts gets n posts from the database.
func (db *Database) GetNPosts(n int) error { return nil }
