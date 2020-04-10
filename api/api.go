// Package api implements an http api server for the entire backend (this entire repo).
package api

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/juju/loggo"
	"github.com/juju/loggo/loggocolor"
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
func newAPI(port int64) (*API, error) {
	api := &API{
		router:   gin.New(),
		database: nil,

		root: common.DefaultAPIRoot,
		port: port,
	}

	err := api.initLogger()
	if err != nil {
		return nil, err
	}
	api.setupRoutes()

	api.log.Infof("API server initialization complete")
	return api, nil
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
	api, err := newAPI(port)
	if err != nil {
		return err
	}

	api.database = database.Connect(false)
	defer api.database.Disconnect()

	api.log.Infof("API server to listen on port %d", port)

	return api.router.Run("0.0.0.0:" + strconv.FormatInt(port, 10))
}

// initLogger initializes the api's logger.
func (api *API) initLogger() error {
	logger := loggo.GetLogger("api")
	err := common.CreateDirIfDoesNotExist(filepath.FromSlash(common.LogsDir))
	if err != nil {
		return err
	}

	// Create the log file
	logFile, err := os.OpenFile(filepath.FromSlash(fmt.Sprintf(
		"%s/logs_%s.txt", common.LogsDir,
		time.Now().Format("2006-01-02_15-04-05"))),
		os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return err
	}

	// Enable colors
	_, err = loggo.ReplaceDefaultWriter(loggocolor.NewWriter(os.Stderr))
	if err != nil {
		return err
	}

	// Register file writer
	err = loggo.RegisterWriter("logs", loggo.NewSimpleWriter(logFile, loggo.DefaultFormatter))
	if err != nil {
		return err
	}

	api.log = &logger // Get a pointer to the logger

	return nil
}
