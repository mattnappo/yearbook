package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/xoreo/yearbook/common"
)

func createPost(w http.ResponseWriter, r *http.Request) {
	logger := common.NewLogger("api.CreateServer")
	w.Header().Set("Content-Type", "application/json") // Set the proper header

	// Decode the post request
	var requestData CreateServerRequest
	json.NewDecoder(r.Body).Decode(&requestData)

	// Extract the data from the request
	port, err := strconv.Atoi(requestData.Port)
	if err != nil {
		logger.Criticalf(err.Error())
	}

	ram, err := strconv.Atoi(requestData.RAM)
	if err != nil {
		logger.Criticalf(err.Error())
	}

	logger.Infof("request to create server with specs:\n%s", requestData.String())

	// Create the new server
	server, err := types.NewServer(requestData.Version, requestData.Name, port, ram)
	if err != nil {
		logger.Criticalf(err.Error())
	}

	logger.Debugf("created new server entry %s", server.Hash.String())

	// Initialize the server
	err = commands.InitializeServer(server)
	if err != nil {
		logger.Criticalf(err.Error())
	}

	logger.Debugf("server initialization complete")

	// Add the newly-created server to the database
	serverDB, err := types.LoadDB()
	if err != nil {
		logger.Criticalf(err.Error())
	}

	err = serverDB.AddServer(server)
	if err != nil {
		logger.Criticalf(err.Error())
	}
	serverDB.Close()

	logger.Debugf("added new server to the database")
	logger.Infof("created new server %s", server.Hash.String())

	json.NewEncoder(w).Encode(*server) // Write to the server
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
