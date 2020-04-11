// Package api implements an http api server for the entire backend (this entire repo).
package api

import (
	"fmt"
	"os"
	"os/signal"
	"path"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/juju/loggo"
	"github.com/juju/loggo/loggocolor"
	"github.com/xoreo/yearbook/common"
	"github.com/xoreo/yearbook/database"
	"golang.org/x/oauth2"
)

const (
	// defaultAPIRoot is the default API root.
	defaultAPIRoot = "/api"

	// defaultOAuthRoot is the default API root for all OAuth2-related
	// requests.
	defaultOAuthRoot = "/oauth"

	// defaultSessionTimeout represenst the expiration time of a session cookie.
	defaultSessionTimeout = time.Minute * 30
)

// API contains the API layer.
type API struct {
	router   *gin.Engine
	database *database.Database
	log      *loggo.Logger

	root      string
	oauthRoot string

	port int64

	oauthConfig *oauth2.Config
	cookieStore cookie.Store
}

// newAPI constructs a new API struct.
func newAPI(port int64) (*API, error) {
	// Generate the store
	// Should use common.GenRandomToken
	cookieStore := cookie.NewStore(
		[]byte(common.GetEnv("COOKIE_SECRET")),
	)
	cookieStore.Options(sessions.Options{
		Path:   "/",
		MaxAge: int(defaultSessionTimeout.Seconds()),
	})

	// Initialize the router
	r := gin.New()
	r.Use(sessions.Sessions("go_session", cookieStore))
	r.Use(gin.Recovery())

	// Construct the API
	api := &API{
		router:   r,
		database: nil,

		root:      defaultAPIRoot,
		oauthRoot: defaultOAuthRoot,

		port: port,

		cookieStore: cookieStore,
	}

	// Setup the logger
	err := api.initLogger()
	if err != nil {
		return nil, err
	}
	api.initializeRoutes()
	api.initializeOAuth()

	api.log.Infof("API server initialization complete")
	return api, nil
}

// initializeRoutes initializes the necessary routes.
func (api *API) initializeRoutes() {
	// Create a group of protected routes (the main api routes)
	protectedRoutes := api.router.Group(api.root)

	// Require only authorized requests
	protectedRoutes.Use(api.authorizeRequest())
	{
		protectedRoutes.GET(path.Join(api.root, "createPost"), api.createPost)
		protectedRoutes.POST(path.Join(api.root, "createPost"), api.createPost)
		protectedRoutes.GET(path.Join(api.root, "getPost/:id"), api.getPost)
		protectedRoutes.GET(path.Join(api.root, "getPosts"), api.getPosts)
		protectedRoutes.GET(path.Join(api.root, "getnPosts/:n"), api.getnPosts)
		protectedRoutes.DELETE(path.Join(api.root, "deletePost/:id"), api.deletePost)

		api.router.POST(path.Join(api.root, "createUser"), api.createUser)
		api.router.GET(path.Join(api.root, "getUser/:username"), api.getUser)
		api.router.GET(path.Join(api.root, "getUsers"), api.getUsers)
		api.router.DELETE(path.Join(api.root, "deleteUser/:username"), api.deleteUser)
	}

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

	// Catch intrerupt
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func(api *API) {
		for sig := range c {
			// Shutdown the API when ctrl-C is pressed
			api.shutdown(sig)
			os.Exit(0)
		}
	}(api)

	return api.router.Run(":" + strconv.FormatInt(port, 10))
}

// shutdown shuts down the API.
func (api *API) shutdown(sig os.Signal) {
	api.log.Debugf("caught %v", sig)
	api.log.Infof("shutting down API server")

	api.database.Disconnect()
	api.log.Debugf("disconnected from database")

	api.router = nil
	api.log.Debugf("destroyed router")

	api.log.Infof("API server shut down")
}

// initLogger initializes the api's logger.
func (api *API) initLogger() error {
	logger := loggo.GetLogger("api")
	err := common.CreateDirIfDoesNotExist(filepath.FromSlash(common.LogsDir))
	if err != nil {
		return err
	}

	logger.SetLogLevel(loggo.DEBUG)

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
