package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xoreo/yearbook/models"
)

func (api *API) createPost(ctx *gin.Context) {
	// Decode the post request
	var request createPostRequest
	err := ctx.ShouldBindJSON(&request)
	if err != nil {
		api.log.Criticalf(err.Error())
		ctx.JSON(
			http.StatusInternalServerError,
			newGenericResponse("", err.Error()),
		)
	}

	api.log.Infof("request to create post")

	// Create the new post
	post, err := models.NewPost(
		request.Sender,
		request.Message,
		request.Images,
		request.Recipients,
	)
	if err != nil {
		api.log.Criticalf(err.Error())
		ctx.JSON(
			http.StatusInternalServerError,
			newGenericResponse("", err.Error()),
		)
	}
	api.log.Debugf("constructed new post %s", post.PostID)

	// Add it to the database
	err = api.database.AddPost(post)
	if err != nil {
		api.log.Criticalf(err.Error())
		ctx.JSON(
			http.StatusInternalServerError,
			newGenericResponse("", err.Error()),
		)
	}
	api.log.Debugf("added post %s to the database", post.PostID)
	api.log.Infof("created new post %s", post.PostID)

	ctx.JSON(http.StatusOK, ok())
}

// getPost gets a post.
func (api *API) getPost(ctx *gin.Context) {

}

// getPosts gets all posts.
func (api *API) getPosts(ctx *gin.Context) {

}

// getnPosts gets n posts.
func (api *API) getnPosts(ctx *gin.Context) {

}

// deletePost deletes a post.
func (api *API) deletePost(ctx *gin.Context) {

}

// createUser creates a user.
func (api *API) createUser(ctx *gin.Context) {

}

// getUser gets a user.
func (api *API) getUser(ctx *gin.Context) {

}

// getUsers gets all users.
func (api *API) getUsers(ctx *gin.Context) {

}

// deleteUser deletes a user.
func (api *API) deleteUser(ctx *gin.Context) {

}
