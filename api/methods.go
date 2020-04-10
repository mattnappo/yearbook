package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/xoreo/yearbook/common"
	"github.com/xoreo/yearbook/models"
)

func createPost(w http.ResponseWriter, r *http.Request) {
	logger := common.NewLogger("api.createPost")
	w.Header().Set("Content-Type", "application/json")

	// Decode the post request
	var request createPostRequest
	json.NewDecoder(r.Body).Decode(&request)

	logger.Infof("request to create post")

	// Create the new post
	post, err := models.NewPost(
		request.Sender,
		request.Message,
		request.Images,
		request.Recipients,
	)
	if err != nil {
		logger.Criticalf(err.Error())
	}
	logger.Debugf("constructed new post %s", post.PostID)

	logger.Debugf("added post %s to the database", post.PostID)
	logger.Infof("created new post %s", post.PostID)

	json.NewEncoder(w).Encode(*post) // Write to the server
}

// getPost gets a post.
func getPost(w http.ResponseWriter, r *http.Request) {
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

func getPosts(w http.ResponseWriter, r *http.Request) {
	logger := common.NewLogger("api.getPosts")
	// Write the response to the server
	json.NewEncoder(w).Encode(res)
}

func getnPosts(w http.ResponseWriter, r *http.Request) {
	logger := common.NewLogger("api.getnPosts")
	// Write the response to the server
	json.NewEncoder(w).Encode(res)
}

func deletePost(w http.ResponseWriter, r *http.Request) {
	logger := common.NewLogger("api.deletePost")
	// Write the response to the server
	json.NewEncoder(w).Encode(res)
}

func createUser(w http.ResponseWriter, r *http.Request) {
	logger := common.NewLogger("api.createUser")
	// Write the response to the server
	json.NewEncoder(w).Encode(res)
}

func getUser(w http.ResponseWriter, r *http.Request) {
	logger := common.NewLogger("api.getUser")
	// Write the response to the server
	json.NewEncoder(w).Encode(res)
}

func getUsers(w http.ResponseWriter, r *http.Request) {
	logger := common.NewLogger("api.getUsers")
	// Write the response to the server
	json.NewEncoder(w).Encode(res)
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	logger := common.NewLogger("api.deleteUser")
	// Write the response to the server
	json.NewEncoder(w).Encode(res)
}
