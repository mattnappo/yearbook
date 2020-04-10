// Package api implements an http api server for the entire backend (this entire repo).
package api

import (
	"path"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/juju/loggo"
	"github.com/xoreo/yearbook/common"
	"github.com/xoreo/yearbook/database"
)

// API contains the API layer.
type API struct {
	router   *gin.Engine
	database *database.Database
	log      *loggo.Logger

	root string
	port int64
}

// newAPI constructs a new API struct.
func newAPI(port int64) *API {
	api := &API{
		router:   gin.New(),
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
	api.router.POST(path.Join(api.root, "createPost"), api.createPost)
	api.router.GET(path.Join(api.root, "getPost/:id"), api.getPost)
	api.router.GET(path.Join(api.root, "getPosts"), api.getPosts)
	api.router.GET(path.Join(api.root, "getnPosts/:n"), api.getnPosts)
	api.router.DELETE(path.Join(api.root, "deletePost/:id"), api.deletePost)

	api.router.POST(path.Join(api.root, "createUser"), api.createUser)
	api.router.GET(path.Join(api.root, "getUser/:username"), api.getUser)
	api.router.GET(path.Join(api.root, "getUsers"), api.getUsers)
	api.router.DELETE(path.Join(api.root, "deleteUser/:username"), api.deleteUser)

	api.log.Infof("initialized API server routes")
}

// StartAPIServer starts the API server.
func StartAPIServer(port int64) error {
	api := newAPI(port)

	api.database = database.Connect(false)
	defer api.database.Disconnect()

	api.log.Infof("API server to listen on port %d", port)

	return api.router.Run("0.0.0.0:" + strconv.FormatInt(port, 10))
}
