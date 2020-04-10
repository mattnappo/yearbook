// Package api implements an http api server for the entire backend (this entire repo).
package api

import (
	"net/http"
	"path"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/juju/loggo"
	"github.com/xoreo/yearbook/common"
)

// API contains the API layer.
type API struct {
	router *mux.Router
	log    *loggo.Logger

	root string
	port int
}

// NewAPI constructs a new API struct.
func NewAPI(port int) *API {
	api := &API{
		router: mux.NewRouter(),
		log:    common.NewLogger("api"),

		root: common.DefaultAPIRoot,
		port: port,
	}

	api.setupRoutes()

	api.log.Infof("API server initialization complete")
	return api
}

// SetupRoutes initializes the necessary routes.
func (api *API) setupRoutes() {
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
	api := NewAPI(port)

	api.log.Infof("API server to listen on port %d", port)
	return http.ListenAndServe(":"+strconv.Itoa(api.port), api.router)
}
