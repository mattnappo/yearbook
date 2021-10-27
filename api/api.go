// Package api implements an http api server for the entire backend (this entire repo).
package api

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/smtp"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/juju/loggo"
	"github.com/juju/loggo/loggocolor"
	"github.com/mattnappo/yearbook/common"
	"github.com/mattnappo/yearbook/database"
	"github.com/mattnappo/yearbook/models"
	"golang.org/x/oauth2"
)

const (
	// defaultAPIRoot is the default API root.
	defaultAPIRoot = "/api"

	// defaultOAuthRoot is the default API root for all OAuth2-related
	// requests.
	defaultOAuthRoot = "/api/oauth"

	// defaultSessionTimeout represenst the expiration time of a session cookie.
	defaultSessionTimeout = time.Minute * 30
)

var (
	protocol = common.GetEnv("PROTOCOL")
	// The frontend callback data
	callbackProvider = common.GetEnv("PROVIDER")
	callbackURL      = "%s://%s/oauth"
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
	r.Use(gin.Logger())

	// Construct the API
	api := &API{
		router:   r,
		database: nil,

		root:      defaultAPIRoot,
		oauthRoot: defaultOAuthRoot,

		port: port,

		callbackURL: fmt.Sprintf(
			callbackURL, protocol, callbackProvider,
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

		ctx.AbortWithStatusJSON( // Respond with the error
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
		protectedRoutes.GET("getNumPosts", api.getNumPosts)
		protectedRoutes.GET("getnPosts/:n", api.getnPosts)
		protectedRoutes.GET("getnPostsOffset/:n/:offset", api.getnPostsOffset)
		protectedRoutes.DELETE("deletePost/:id", api.deletePost)

		protectedRoutes.PATCH("updateUser", api.updateUser)
		protectedRoutes.GET("getUser/:username", api.getUser)
		protectedRoutes.GET("getUserProfilePic/:username", api.getUserProfilePic)
		protectedRoutes.GET("getUserGrade/:username", api.getUserGrade)
		protectedRoutes.GET("getActivity/:username", api.getActivity)
		protectedRoutes.GET("getUserPosts/:username", api.getUserPosts)
		protectedRoutes.GET("getUserWithAuthentication/:username", api.getUserWithAuth)
		protectedRoutes.GET("getUsers", api.getUsers)
		protectedRoutes.GET("getSeniors", api.getSeniors)
		protectedRoutes.GET("getUsernames", api.getUsernames)
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
	err := common.CreateDirIfDoesNotExist(
		filepath.FromSlash(common.LogsDir),
	)
	if err != nil {
		return err
	}

	logger.SetLogLevel(loggo.INFO)

	// Create the log file
	logFile, err := os.OpenFile(filepath.FromSlash(fmt.Sprintf(
		"%s/internal_log_%s.txt", common.LogsDir,
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
	err = loggo.RegisterWriter("logs", loggo.NewSimpleWriter(
		logFile, loggo.DefaultFormatter,
	))
	if err != nil {
		return err
	}

	api.log = &logger

	return nil
}

// genEmailBody generates the body of a notification email given
// a sender's username.
func genEmailBody(sender string) string {
	emailTemplate, _ := ioutil.ReadFile("email.txt")
	return strings.ReplaceAll(
		string(emailTemplate), "$$$SENDER$$$", sender,
	)
}

// sendNotification sends an email to a user that they have
// been congratulated.
func (api *API) sendNotification(
	sender models.Username,
	recipients []models.Username,
) error {
	// Setup the authentication
	auth := smtp.PlainAuth("",
		common.NotifEmail,
		common.NotifPassword,
		common.NotifProvider,
	)

	htmlBody := genEmailBody(sender.Name())

	// Setup the message
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"

	var to []string
	for _, recip := range recipients {
		to = append(to, recip.Email())
	}
	msg := fmt.Sprintf("To: %s\r\nSubject: %s Congratulated you!\r\n"+
		mime+"\r\n"+
		"%s\r\n",
		strings.Join(to, ","),
		sender.Name(),
		htmlBody,
	)

	// Actually send the email
	err := smtp.SendMail(
		fmt.Sprintf("%s:%d", common.NotifProvider, common.NotifPort),
		auth, common.NotifEmail, to, []byte(msg),
	)
	return err
}

func (api *API) sendModEmail(post models.Post) error {
	// Setup the authentication
	auth := smtp.PlainAuth("",
		common.NotifEmail,
		common.NotifPassword,
		common.NotifProvider,
	)

	body := fmt.Sprintf(
		"Sender: %s\nRecipients: %v\nMessage: %s\nID: %d\nPostID: %s\n",
		post.Sender, post.Recipients, post.Message, post.ID, post.PostID,
	)

	// Setup the message
	to := []string{"mattnappo@gmail.com"}
	msg := fmt.Sprintf("To: %s\r\nSubject: %s posted\r\n"+
		"\r\n"+
		"%s\r\n",
		to,
		post.Sender.Name(),
		body,
	)

	// Actually send the email
	err := smtp.SendMail(
		fmt.Sprintf("%s:%d", common.NotifProvider, common.NotifPort),
		auth, common.NotifEmail, to, []byte(msg),
	)
	return err

}
