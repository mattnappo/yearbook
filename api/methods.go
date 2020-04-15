package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/xoreo/yearbook/models"
	"golang.org/x/oauth2"
)

// createPost creates a new post.
func (api *API) createPost(ctx *gin.Context) {
	api.log.Infof("request to create post")
	// Decode the post request
	var request createPostRequest
	err := ctx.ShouldBindJSON(&request)
	if api.check(err, ctx, http.StatusBadRequest) {
		return
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

	api.log.Debugf("constructed new post %s", post.String())

	// Add it to the database
	err = api.database.AddPost(post)
	if api.check(err, ctx) {
		return
	}

	api.log.Infof("created new post %s", post.String())
	ctx.JSON(http.StatusOK, ok())
}

// getPost gets a post.
func (api *API) getPost(ctx *gin.Context) {
	id := ctx.Param("id")
	api.log.Infof("request to get post %s", id)

	post, err := api.database.GetPost(id)
	if api.check(err, ctx) {
		return
	}

	api.log.Infof("got post %s from database", post.String())
	ctx.JSON(http.StatusOK, gr(post))
}

// getPosts gets all posts.
func (api *API) getPosts(ctx *gin.Context) {
	api.log.Infof("request to get all posts")

	posts, err := api.database.GetAllPosts()
	if api.check(err, ctx) {
		return
	}

	api.log.Infof("got all posts")
	ctx.JSON(http.StatusOK, gr(posts))
}

// getnPosts gets n posts.
func (api *API) getnPosts(ctx *gin.Context) {
	n := ctx.Param("n")
	api.log.Infof("request to get %s posts", n)

	nInt, err := strconv.Atoi(n)
	if api.check(err, ctx) {
		return
	}

	posts, err := api.database.GetnPosts(nInt)
	if api.check(err, ctx) {
		return
	}

	api.log.Infof("got %d posts", nInt)
	ctx.JSON(http.StatusOK, gr(posts))
}

// deletePost deletes a post.
func (api *API) deletePost(ctx *gin.Context) {
	postID := ctx.Param("id")
	api.log.Infof("request to delete post %s", postID)

	err := api.database.DeletePost(postID)
	if api.check(err, ctx) {
		return
	}

	api.log.Infof("deleted post %s", postID)
	ctx.JSON(http.StatusOK, ok())
}

// createUser creates a user.
func (api *API) createUser(ctx *gin.Context) {
	api.log.Infof("request to create user")

	// Decode the new user request
	var request createUserRequest
	err := ctx.ShouldBindJSON(&request)
	if api.check(err, ctx, http.StatusBadRequest) {
		return
	}

	// Determine the grade from the request
	var grade models.Grade
	switch request.Grade {
	case "freshman":
		grade = models.Freshman
		break
	case "sophomore":
		grade = models.Sophomore
		break
	case "junior":
		grade = models.Junior
		break
	case "senior":
		grade = models.Senior
		break
	default:
		api.check(fmt.Errorf("invalid grade %v", request.Grade), ctx)
		return
	}

	// Create the new user
	user, err := models.NewUser(request.Email, grade)
	if api.check(err, ctx) {
		return
	}
	api.log.Debugf("constructed new user %s", user.String())

	// Add it to the database
	err = api.database.AddUser(user)
	if api.check(err, ctx) {
		return
	}

	api.log.Infof("created new user %s", user.String())
	ctx.JSON(http.StatusOK, ok())
}

// updateUser handles a request to update a user.
func (api *API) updateUser(ctx *gin.Context) {
	api.log.Infof("request to update user")

	// Decode the request data
	var request updateUserRequest
	err := ctx.ShouldBindJSON(&request)
	if api.check(err, ctx, http.StatusBadRequest) {
		return
	}
	api.log.Debugf("updating user with new info %v", request)

	// Construct a user struct with the new user data in it, and everything
	// else blank.
	newUserData, err := models.UserFromString(request.String())
	if api.check(err, ctx) {
		return
	}

	// Check that the username of the request is the same as the username
	// of the account being modified.
	bearerToken, err := extractBearerToken(ctx)
	if api.check(err, ctx, http.StatusUnauthorized) {
		return
	}

	// Query the Google API to get the email associated with the token
	// in the request header.
	headerToken := &oauth2.Token{AccessToken: bearerToken}
	googleUser, err := api.getUserInfo(headerToken)
	if api.check(err, ctx, http.StatusUnauthorized) {
		return
	}

	// Marshal the request data into a new models.user
}

// getUser gets a user.
func (api *API) getUser(ctx *gin.Context) {
	username := ctx.Param("username")
	api.log.Infof("request to get user %s", username)

	user, err := api.database.GetUser(username)
	if api.check(err, ctx) {
		return
	}

	api.log.Infof("got user %s", user.String())
	ctx.JSON(http.StatusOK, gr(user))
}

// getUsers gets all users.
func (api *API) getUsers(ctx *gin.Context) {
	api.log.Infof("request to get all users")

	users, err := api.database.GetAllUsers()
	if api.check(err, ctx) {
		return
	}

	api.log.Infof("got all users")
	ctx.JSON(http.StatusOK, gr(users))

}

// deleteUser deletes a user.
func (api *API) deleteUser(ctx *gin.Context) {
	username := ctx.Param("username")
	api.log.Infof("request to delete user %s", username)

	err := api.database.DeleteUser(username)
	if api.check(err, ctx) {
		return
	}

	api.log.Infof("deleted user %s", username)
	ctx.JSON(http.StatusOK, ok())
}
