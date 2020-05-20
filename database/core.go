package database

import (
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
	err := db.DB.Model(&posts).Order("id DESC").Limit(n).Select()
	if err != nil {
		return nil, err
	}
	return posts, nil
}

// GetnPostsWithOffset gets n posts at a ceratin offset.
func (db *Database) GetnPostsWithOffset(
	n, offset int,
) ([]models.Post, error) {
	var posts []models.Post
	err := db.DB.Model(&posts).
		Limit(n).
		Offset(offset).
		Order("id DESC").
		Select()

	return posts, err
}

// GetNumPosts returns the number of posts in the database.
func (db *Database) GetNumPosts() (int, error) {
	return db.DB.Model((*models.Post)(nil)).Count()
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
	lookupUser.Grade = user.Grade
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

// GetUserInbound returns the inbound posts of a user.
func (db *Database) GetUserInbound(username string) ([]models.Post, error) {
	var inboundPostIDs models.User

	// Get the list of inbound postIDs
	err := db.DB.Model((*models.User)(nil)).
		Column("inbound_posts").
		Where("username = ?", username).
		Select(&inboundPostIDs)
	if err != nil {
		return nil, err
	}

	// Get all of the posts given the postIDs
	var posts []models.Post
	for _, inboundPostID := range inboundPostIDs.InboundPosts {
		post, _ := db.GetPost(inboundPostID) // CHECK THIS ERROR SOMEHOW
		posts = append(posts, post)
	}
	return posts, nil
}

// GetUserInboundOutbound returns partial information about the inbound
// and outbound posts of a user.
func (db *Database) GetUserInboundOutbound(username string) ([][]models.Post, error) {
	var postIDs models.User

	// Get the list of inbound and outbound postIDs
	err := db.DB.Model((*models.User)(nil)).
		Column("inbound_posts", "outbound_posts").
		Where("username = ?", username).
		Select(&postIDs)
	if err != nil {
		return nil, err
	}

	// Get all of the necessary post data given the post IDs
	// Really these should throw errors
	inboundPosts := db.traversePosts(postIDs.InboundPosts)
	outboundPosts := db.traversePosts(postIDs.OutboundPosts)
	return [][]models.Post{inboundPosts, outboundPosts}, nil
}

// GetUserProfilePic gets a user's profile pic given username.
func (db *Database) GetUserProfilePic(username string) (string, error) {
	var profilePic string
	err := db.DB.Model((*models.User)(nil)).
		Column("profile_pic").
		Where("username = ?", username).
		Select(&profilePic)

	return profilePic, err
}

// GetUserGrade gets a user's grade given a username.
func (db *Database) GetUserGrade(username string) (models.Grade, error) {
	var grade models.Grade
	err := db.DB.Model((*models.User)(nil)).
		Column("grade").
		Where("username = ?", username).
		Select(&grade)

	return grade, err
}

// GetAllUsers gets all users from the database.
func (db *Database) GetAllUsers() ([]models.User, error) {
	var users []models.User
	err := db.DB.Model(&users).Select()
	return users, err
}

// GetAllUsernames gets all usernames in the database.
func (db *Database) GetAllUsernames() ([]string, error) {
	var usernames []string
	err := db.DB.Model((*models.User)(nil)).
		Column("username").
		Select(&usernames)
	return usernames, err
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

// traversePosts returns some data about posts given a []string of postIDs.
func (db *Database) traversePosts(postIDs []string) []models.Post {
	var posts []models.Post
	for _, postID := range postIDs {
		var post models.Post
		db.DB.Model((*models.Post)(nil)).
			Where("post_id = ?", postID).
			Select(&post)
		posts = append(posts, post)
	}
	return posts
}
