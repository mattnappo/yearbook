package api

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/xoreo/yearbook/models"
)

// check checks for an error. Returns true if the request shuold be
// terminated, false if it shold stay alive.
func (api *API) check(err error, ctx *gin.Context) bool {
	if err != nil {
		api.log.Criticalf(err.Error())
		ctx.JSON( // Respond with the error
			http.StatusInternalServerError, gr("", err.Error()),
		)
		return true
	}
	return false
}

// createPost creates a new post.
func (api *API) createPost(ctx *gin.Context) {
	api.log.Infof("request to create post")

	// Decode the post request
	var request createPostRequest
	err := ctx.ShouldBindJSON(&request)
	if api.check(err, ctx) {
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

	api.log.Infof("got %d posts", n)
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
	if api.check(err, ctx) {
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
		if api.check(errors.New("invalid grade"), ctx) {
			return
		}
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
