// Package api implements an http api server for the entire backend (this entire repo).
package api

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
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
	callbackURL string
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
		// Secure: true,
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

		callbackURL: fmt.Sprintf(
			"http://localhost:%d/oauth", 3000,
		),
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

// check checks for an error. Returns true if the request shuold be
// terminated, false if it shold stay alive.
func (api *API) check(err error, ctx *gin.Context, status ...int) bool {
	if err != nil {
		api.log.Criticalf(err.Error()) // Log the error
		// Respond with correct status code and error
		var statusCode int
		switch len(status) {
		case 1: // If the user supplied a different status
			statusCode = status[0]
			break
		default:
			statusCode = http.StatusInternalServerError
			break
		}

		// if statusCode == http.StatusUnauthorized {
		// 	ctx.Redirect(http.StatusPermanentRedirect, "http://localhost:3000/")
		// 	return true
		// }

		ctx.AbortWithStatusJSON( // Respond with the error}
			statusCode, gr("", err.Error()),
		)
		return true
	}
	return false
}

// initializeRoutes initializes the necessary routes.
func (api *API) initializeRoutes() {
	// Create a group of protected routes (the main api routes)
	protectedRoutes := api.router.Group(api.root)

	// Require only authorized requests
	protectedRoutes.Use(api.authorizeRequest())
	{
		protectedRoutes.POST("createPost", api.createPost)
		protectedRoutes.GET("getPost/:id", api.getPost)
		protectedRoutes.GET("getPosts", api.getPosts)
		protectedRoutes.GET("getnPosts/:n", api.getnPosts)
		protectedRoutes.DELETE("deletePost/:id", api.deletePost)

		protectedRoutes.POST("createUser", api.createUser)
		protectedRoutes.PATCH("updateUser", api.updateUser)
		protectedRoutes.GET("getUser/:username", api.getUser)
		protectedRoutes.GET("getActivity/:username", api.getActivity)
		protectedRoutes.GET("getUserWithAuthentication/:username", api.getUserWithAuth)
		protectedRoutes.GET("getUsers", api.getUsers)
		protectedRoutes.GET("getSeniors", api.getSeniors)
		protectedRoutes.GET("getUsernames", api.getUsernames)
		protectedRoutes.DELETE("deleteUser/:username", api.deleteUser)
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
