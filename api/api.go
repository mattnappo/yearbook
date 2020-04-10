// Package api implements an http api server for the entire backend (this entire repo).
package api

import (
	"net/http"
	"path"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/juju/loggo"
	"github.com/xoreo/yearbook/common"
	"github.com/xoreo/yearbook/database"
)

// API contains the API layer.
type API struct {
	router   *mux.Router
	database *database.Database
	log      *loggo.Logger

	root string
	port int
}

// newAPI constructs a new API struct.
func newAPI(port int) *API {
	api := &API{
		router:   mux.NewRouter(),
		log:      common.NewLogger("api"),
		database: nil,

		root: common.DefaultAPIRoot,
		port: port,
	}

	api.setupRoutes()

	api.log.Infof("API server initialization complete")
	return api
}

// setupRoutes initializes the necessary routes.
func (api *API) setupRoutes() {
	// Oh boi do I need to clean up this bad boi
	api.router.HandleFunc(path.Join(api.root, "createPost"), createPost).
		Methods("POST")
	api.router.HandleFunc(path.Join(api.root, "getPost/{id}"), getPost).
		Methods("GET")
	api.router.HandleFunc(path.Join(api.root, "getPosts"), getPosts).
		Methods("GET")
	api.router.HandleFunc(path.Join(api.root, "getnPosts/{n}"), getnPosts).
		Methods("GET")
	api.router.HandleFunc(path.Join(api.root, "deletePost/{id}"), deletePost).
		Methods("DELETE")

	api.router.HandleFunc(path.Join(api.root, "createUser"), createUser).
		Methods("POST")
	api.router.HandleFunc(path.Join(api.root, "getUser/{id}"), getUser).
		Methods("GET")
	api.router.HandleFunc(path.Join(api.root, "getUsers"), getUsers).
		Methods("GET")
	api.router.HandleFunc(path.Join(api.root, "deleteUser/{id}"), deleteUser).
		Methods("DELETE")

	api.log.Infof("initialized API server routes")
}

// StartAPIServer starts the API server.
func StartAPIServer(port int) error {
	api := newAPI(port)

	api.database = database.Connect(false)
	defer api.database.Disconnect()

	api.log.Infof("API server to listen on port %d", port)
	return http.ListenAndServe(":"+strconv.Itoa(api.port), api.router)
}
