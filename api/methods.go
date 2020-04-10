package api

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/mux"
	"github.com/xoreo/yearbook/common"
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
	logger := common.NewLogger("api.GetServer")
	w.Header().Set("Content-Type", "application/json")

	// Extract the server hash from the request
	hashString := mux.Vars(r)["post_id"]

	logger.Infof("request to get %s", hashString)

	logger.Debugf("got post %s from database", server.Hash.String())

	// Prepare the response
	res := GETServerResponse{
		server.TimeCreated,
		server.ID,
		*server.Properties,
		server.GetCoreProperties(),
	}

	// Write the response to the server
	json.NewEncoder(w).Encode(res)
}

func (api *API) getPosts(ctx *gin.Context) {
	logger := common.NewLogger("api.getPosts")
	// Write the response to the server
	json.NewEncoder(w).Encode(res)
}

func (api *API) getnPosts(ctx *gin.Context) {
	logger := common.NewLogger("api.getnPosts")
	// Write the response to the server
	json.NewEncoder(w).Encode(res)
}

func (api *API) deletePost(ctx *gin.Context) {
	logger := common.NewLogger("api.deletePost")
	// Write the response to the server
	json.NewEncoder(w).Encode(res)
}

func (api *API) createUser(ctx *gin.Context) {
	logger := common.NewLogger("api.createUser")
	// Write the response to the server
	json.NewEncoder(w).Encode(res)
}

func (api *API) getUser(ctx *gin.Context) {
	logger := common.NewLogger("api.getUser")
	// Write the response to the server
	json.NewEncoder(w).Encode(res)
}

func (api *API) getUsers(ctx *gin.Context) {
	logger := common.NewLogger("api.getUsers")
	// Write the response to the server
	json.NewEncoder(w).Encode(res)
}

func (api *API) deleteUser(ctx *gin.Context) {
	logger := common.NewLogger("api.deleteUser")
	// Write the response to the server
	json.NewEncoder(w).Encode(res)
}
