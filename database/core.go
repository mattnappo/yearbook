package database

import (
	"fmt"

	"github.com/go-pg/pg"
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
	if user.Grade != 0 {
		lookupUser.Grade = user.Grade
	}
	// if user.ProfilePic != nil {
	// 	lookupUser.ProfilePic = user.ProfilePic
	// }
	if user.Nickname != "" {
		lookupUser.Nickname = user.Nickname
	}

	err = db.DB.Update(&lookupUser)
	if err != nil {
		return err
	}
	return nil
}

// AddToAndFrom populates the InboundPosts and OutboundPosts data
// within a user
func (db *Database) AddToAndFrom(
	postID, senderUsername string,
	recipientUsernames []string,
) error {
	// Get the list of outbound posts from the user and append the new
	// outbound postID to the array in the sender db entry
	sender, err := db.GetUser(senderUsername)
	if err != nil {
		return err
	}
	outbound := sender.OutboundPosts
	outbound = append(outbound, postID)
	// Update the array in the database
	_, err = db.DB.Model(&sender).
		Set("outbound_posts = ?", outbound).
		Where("id = ?", sender.ID).
		Update()
	if checkIntegrity(err) != nil {
		return err
	}

	// Do the same thing as above, but for each recipient
	for _, recipientUsername := range recipientUsernames {
		recipient, err := db.GetUser(recipientUsername)
		if err != nil {
			return err
		}
		inbound := sender.InboundPosts
		inbound = append(inbound, postID)
		// Update the array in the database
		_, err = db.DB.Model(&recipient).
			Set("inbound_posts = ?", inbound).
			Where("id = ?", recipient.ID).
			Update()
		if checkIntegrity(err) != nil {
			return err
		}

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

// GetUserInbound returns the inbound posts to a user.
func (db *Database) GetUserInbound(username string) ([]models.Post, error) {
	var inboundPostIDs []string

	// Get the list of inbound postIDs
	/*
		err := db.DB.Model((*models.User)(nil)).
			Column("inbound_posts").
			Where("username = ?", username).
			Select(&inboundPostIDs)
	*/
	var inboundPost *interface{}
	err := db.DB.Model((*models.User)(nil)).
		Column("inbound_posts").
		Where("username = ?", username).
		Select(&inboundPostIDs)
	fmt.Printf("THING: \n%v\n", inboundPost)
	if err != nil {
		return nil, err
	}
	fmt.Printf("\n\n %s\n\n", inboundPostIDs[0])
	var posts []models.Post
	post, err := db.GetPost("a6b9d19f01c0205d5da39b734902273384a5b493d422b6240f953ba521438c85")
	if err != nil {
		return nil, err
	}
	posts = append(posts, post)
	/*
		// Lookup each post in that database
		var posts []models.Post
		for _, inboundPostID := range inboundPostIDs {
			fmt.Printf("\n\n%s\n\n", inboundPostID)
			post, err := db.GetPost(inboundPostID)
			if err != nil {
				return nil, err
			}
			posts = append(posts, post)
		}
	*/
	return posts, nil
}

// GetAllUsers gets all users from the database.
func (db *Database) GetAllUsers() ([]models.User, error) {
	var users []models.User
	err := db.DB.Model(&users).Select()
	return users, err
}

// GetAllSeniorUsernames gets all of the usernames of all of the seniors.
func (db *Database) GetAllSeniorUsernames() ([]string, error) {
	var usernames []string
	err := db.DB.Model((*models.User)(nil)).
		Column("username").
		Where("grade = 3").
		Select(&usernames)

	return usernames, err
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

// checkIntegrity checks the integrity of a postgres model function return.
func checkIntegrity(err error) error {
	// Return the error as long as it is not a duplicate key violation.
	if err != nil {
		pgErr, ok := err.(pg.Error)
		if ok && pgErr.IntegrityViolation() {
			return nil
		}
		return err
	}
	return nil
}

/*
// checkNoResults checks if the error is a postgres null result error.
func checkNoResults(err error) error {
	// Return the error as long as it is not a no results in set error.
	if err != nil {
		pgErr, ok := err.(pg.Error)
		if ok && pgErr.Field() == pg.ErrNoRows.Error() {
			return nil
		}
		return err
	}
	return nil
}
*/
