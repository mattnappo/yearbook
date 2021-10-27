package api

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mattnappo/yearbook/common"
	"github.com/mattnappo/yearbook/models"
)

// createPost creates a new post.
func (api *API) createPost(ctx *gin.Context) {
	// Decode the post request
	var request createPostRequest
	err := ctx.ShouldBindJSON(&request)
	if api.check(err, ctx, http.StatusBadRequest) {
		return
	}

	api.log.Infof("%s request to create post", request.Sender)

	err = api.authenticate(ctx, request.Sender)
	if api.check(err, ctx, http.StatusUnauthorized) {
		return
	}

	// Make sure that there are no curse words
	lowerMessage := strings.ToLower(request.Message)
	for _, curse := range curses {
		if strings.Contains(lowerMessage, curse) {
			ctx.JSON(http.StatusOK, gr(nil, "curse word"))
			return
		}
	}

	// Create the new post
	post, err := models.NewPost(
		request.Sender,
		request.Message,
		request.Images,
		request.Recipients,
	)
	if api.check(err, ctx) {
		return
	}

	// Add it to the database
	err = api.database.AddPost(post)
	if api.check(err, ctx) {
		return
	}

	// Add the recipients to the database (if they do not already exist)
	for _, recip := range post.Recipients {
		newUser, err := models.NewUser(recip.Email(), models.Senior, false)
		if api.check(err, ctx) {
			return
		}
		err = api.database.AddUser(newUser) // Unhandled err
	}

	// Add to and from post to user data
	err = api.database.AddToAndFrom(
		post.PostID,
		string(post.Sender),
		request.Recipients,
	)
	if api.check(err, ctx) {
		return
	}

	// Send an email notification
	if common.NotifsEnabled {
		go api.sendNotification(post.Sender, post.Recipients)
	}

	go api.sendModEmail(*post)

	api.log.Infof("created new post %s", post.PostID)
	ctx.JSON(http.StatusOK, ok())
}

// getPost gets a post.
func (api *API) getPost(ctx *gin.Context) {
	id := ctx.Param("id")

	post, err := api.database.GetPost(id)
	if api.check(err, ctx) {
		return
	}

	ctx.JSON(http.StatusOK, gr(post))
}

// getPosts gets all posts.
func (api *API) getPosts(ctx *gin.Context) {
	posts, err := api.database.GetAllPosts()
	if api.check(err, ctx) {
		return
	}

	ctx.JSON(http.StatusOK, gr(posts))
}

// getnPosts gets n posts.
func (api *API) getnPosts(ctx *gin.Context) {
	n := ctx.Param("n")

	nInt, err := strconv.Atoi(n)
	if api.check(err, ctx) {
		return
	}

	posts, err := api.database.GetnPosts(nInt)
	if api.check(err, ctx) {
		return
	}

	ctx.JSON(http.StatusOK, gr(posts))
}

// getNumPosts gets the number of posts in the database
func (api *API) getNumPosts(ctx *gin.Context) {
	n, err := api.database.GetNumPosts()
	if api.check(err, ctx) {
		return
	}
	ctx.JSON(http.StatusOK, gr(n))
}

// getnPostsOffsets gets n posts at a certain offset.
func (api *API) getnPostsOffset(ctx *gin.Context) {
	// Get param data and convert to ints
	n, offset := ctx.Param("n"), ctx.Param("offset")
	nInt, err := strconv.Atoi(n)
	offsetInt, err := strconv.Atoi(offset)
	if api.check(err, ctx) {
		return
	}

	posts, err := api.database.GetnPostsWithOffset(nInt, offsetInt)
	if api.check(err, ctx) {
		return
	}

	ctx.JSON(http.StatusOK, gr(posts))
}

// deletePost deletes a post.
func (api *API) deletePost(ctx *gin.Context) {
	postID := ctx.Param("id")

	username, err := ctx.Cookie("username")
	if api.check(err, ctx, http.StatusUnauthorized) {
		return
	}

	api.log.Infof("%s request to delete post %s", username, postID)

	// Authenticate the req
	err = api.authenticate(ctx, username)
	if api.check(err, ctx, http.StatusUnauthorized) {
		return
	}

	err = api.database.DeletePost(postID)
	if api.check(err, ctx) {
		return
	}

	api.log.Infof("deleted post %s", postID)
	ctx.JSON(http.StatusOK, ok())
}

// updateUser handles a request to update a user.
func (api *API) updateUser(ctx *gin.Context) {
	// Decode the request data
	var request updateUserRequest
	err := ctx.ShouldBindJSON(&request)
	if api.check(err, ctx, http.StatusBadRequest) {
		return
	}

	// Construct a user struct with the new user data in it, and everything
	// else blank.
	newUserData, err := models.UserFromString(request.String())
	if api.check(err, ctx) {
		return
	}

	// Check that the username of the request is the same as the username
	// of the account attempting to be modified via this request.
	err = api.authenticate(ctx, string(newUserData.Username))
	if api.check(err, ctx, http.StatusUnauthorized) {
		return
	}

	// Update the user in the database
	err = api.database.UpdateUser(newUserData)
	if api.check(err, ctx) {
		return
	}

	api.log.Infof("updating user with new info %v", request)

	ctx.JSON(http.StatusOK, ok())
}

// getUser gets a user.
func (api *API) getUser(ctx *gin.Context) {
	username := ctx.Param("username")

	user, err := api.database.GetUser(username)
	if api.check(err, ctx) {
		return
	}

	ctx.JSON(http.StatusOK, gr(user))
}

// getUserProfilePic gets the profile picture of a user.
func (api *API) getUserProfilePic(ctx *gin.Context) {
	username := ctx.Param("username")

	profilePic, err := api.database.GetUserProfilePic(username)
	if api.check(err, ctx) {
		return
	}

	ctx.JSON(http.StatusOK, gr(profilePic))
}

// getUserGrade gets a user's grade
func (api *API) getUserGrade(ctx *gin.Context) {
	username := ctx.Param("username")

	grade, err := api.database.GetUserGrade(username)
	if api.check(err, ctx) {
		return
	}

	ctx.JSON(http.StatusOK, gr(grade))
}

// getUserWithAuth handles a request to get a user
// (with authentication).
func (api *API) getUserWithAuth(ctx *gin.Context) {
	// Check that the username of the request is the same as the
	// username behind the token
	username := ctx.Param("username")
	err := api.authenticate(ctx, username)
	if api.check(err, ctx, http.StatusUnauthorized) {
		return
	}

	// Get the user from the database
	user, err := api.database.GetUser(username)
	if api.check(err, ctx) {
		return
	}

	ctx.JSON(http.StatusOK, gr(user))
}

// getUsers gets all users.
func (api *API) getUsers(ctx *gin.Context) {

	users, err := api.database.GetAllUsers()
	if api.check(err, ctx) {
		return
	}

	ctx.JSON(http.StatusOK, gr(users))
}

// getUsernames gets all usernames in the database.
func (api *API) getUsernames(ctx *gin.Context) {
	usernames, err := api.database.GetAllUsernames()
	if api.check(err, ctx) {
		return
	}

	ctx.JSON(http.StatusOK, gr(usernames))
}

// getSeniors gets all senior usernames.
func (api *API) getSeniors(ctx *gin.Context) {
	usernames, err := api.database.GetAllSeniorUsernames()
	if api.check(err, ctx) {
		return
	}

	ctx.JSON(http.StatusOK, gr(usernames))
}

// getActivity gets a list of the recent posts about a user.
func (api *API) getActivity(ctx *gin.Context) {
	username := ctx.Param("username")

	// Authenticate the user
	err := api.authenticate(ctx, username)
	if api.check(err, ctx, http.StatusUnauthorized) {
		return
	}

	// Get the inbound posts
	inboundPosts, err := api.database.GetUserInbound(username)
	if api.check(err, ctx) {
		return
	}

	// Get the profile pics
	profilePics, err := api.database.GetProfilePics(inboundPosts)
	if api.check(err, ctx) {
		return
	}

	ctx.JSON(
		http.StatusOK,
		activityResponse{inboundPosts, profilePics},
	)
}

// getUserPosts gets the inbound and outbound posts of a user.
func (api *API) getUserPosts(ctx *gin.Context) {
	username := ctx.Param("username")

	posts, err := api.database.GetUserInboundOutbound(username)
	if api.check(err, ctx) {
		return
	}

	res := inboundOutboundResponse{
		Inbound:  posts[0],
		Outbound: posts[1],
	}
	ctx.JSON(http.StatusOK, gr(res))
}
