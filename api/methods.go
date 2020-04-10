package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xoreo/yearbook/models"
)

// check checks for an error.
func (api *API) check(err error, ctx *gin.Context) {
	if err != nil {
		api.log.Criticalf(err.Error())
		ctx.JSON(
			http.StatusInternalServerError,
			newGenericResponse("", err.Error()),
		)
	}
}

// createPost creates a new post.
func (api *API) createPost(ctx *gin.Context) {
	// Decode the post request
	var request createPostRequest
	err := ctx.ShouldBindJSON(&request)
	api.check(err, ctx)

	api.log.Infof("request to create post")

	// Create the new post
	post, err := models.NewPost(
		request.Sender,
		request.Message,
		request.Images,
		request.Recipients,
	)
	api.check(err, ctx)

	api.log.Debugf("constructed new post %s", post.String())

	// Add it to the database
	err = api.database.AddPost(post)
	api.check(err, ctx)

	api.log.Debugf("added post %s to the database", post.String())
	api.log.Infof("created new post %s", post.PostID)

	ctx.JSON(http.StatusOK, ok())
}

// getPost gets a post.
func (api *API) getPost(ctx *gin.Context) {
	id := ctx.Param("id")
	api.log.Infof("request to get post %s", id)

	post, err := api.database.GetPost(id)
	api.check(err, ctx)

	api.log.Debugf("got post %s from database", post.String())
	api.log.Infof("got post %s", post.ID)

	ctx.JSON(http.StatusOK, newGenericResponse(post.String()))
}

// getPosts gets all posts.
func (api *API) getPosts(ctx *gin.Context) {}

// getnPosts gets n posts.
func (api *API) getnPosts(ctx *gin.Context) {}

// deletePost deletes a post.
func (api *API) deletePost(ctx *gin.Context) {}

// createUser creates a user.
func (api *API) createUser(ctx *gin.Context) {}

// getUser gets a user.
func (api *API) getUser(ctx *gin.Context) {}

// getUsers gets all users.
func (api *API) getUsers(ctx *gin.Context) {}

// deleteUser deletes a user.
func (api *API) deleteUser(ctx *gin.Context) {}
