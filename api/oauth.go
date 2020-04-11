package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"path"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/xoreo/yearbook/common"
	"github.com/xoreo/yearbook/crypto"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// errUnauthorized is thrown when a request could not be authorized.
var errUnauthorized = errors.New("failed to authorize request")

// user is a retrieved and authentiacted user.
type user struct {
	Sub           string `json:"sub"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Profile       string `json:"profile"`
	Picture       string `json:"picture"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Gender        string `json:"gender"`
}

// initializeOAuth configures the API's OAuth2 config.
func (api *API) initializeOAuth() {
	// Configure the OAuth2 client
	api.oauthConfig = &oauth2.Config{
		RedirectURL:  api.callbackURL,
		ClientID:     common.GetEnv("GOOGLE_CLIENT_ID"),
		ClientSecret: common.GetEnv("GOOGLE_CLIENT_SECRET"),
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
		},
		Endpoint: google.Endpoint,
	}

	api.initializeOAuthRoutes()
}

// initializeOAuthRoutes initializes the OAuth2-related API routes.
func (api *API) initializeOAuthRoutes() {
	api.router.GET("/home", api.home)

	api.router.GET(path.Join(api.oauthRoot, "login"), api.login)
	api.router.GET(path.Join(api.oauthRoot, "authorize"), api.authorize)

	api.log.Infof("initialized API server OAuth2 routes")

}

// getLoginURL gets the "Sign in with Google" URL.
func (api *API) getLoginURL(state string) string {
	return api.oauthConfig.AuthCodeURL(state)
}

// authorizeRequest is used to authorize a request for a certain
// endpoint group.
func (api *API) authorizeRequest() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Get the token from the header
		bearerToken := strings.Split(ctx.GetHeader("Authorization"), " ")[1]
		api.log.Infof("attempting to authorize request with token %s",
			bearerToken,
		)

		// This should make client.Get call on behalf of client with token,
		// get sub, then look up in a PG db with sub and check that the tokens
		// match.

		if bearerToken != "test" {
			api.log.Infof(errUnauthorized.Error())
			api.check(errUnauthorized, ctx)
			return
		}
		ctx.Next()
	}
}

// home handles requests to the home ("/home").
func (api *API) home(ctx *gin.Context) {
	ctx.Writer.Write([]byte("welcome home!"))
}

// login handles a request to login.
func (api *API) login(ctx *gin.Context) {
	api.log.Infof("request to login")

	// Generate a random token for the cookie session handler
	// Encrypt this cookie later
	state := crypto.GenRandomToken()
	session := sessions.Default(ctx)
	session.Set("state", state)
	err := session.Save()
	if api.check(err, ctx) {
		return
	}

	api.log.Debugf("generated random token %s", state)

	gotState := session.Get("state")
	api.log.Debugf("state: %v", gotState)

	// will returb
	ctx.Writer.Write([]byte("<html><title>Golang Google</title> <body> <a href='" + api.getLoginURL(state) + "'><button>Login with Google!</button> </a> </body></html>"))
	// Respond with the URL to "Sign in with Google"
	// ctx.JSON(http.StatusOK, api.getLoginURL(state))
}

// authorize is the Google authorization callback URL.
func (api *API) authorize(ctx *gin.Context) {
	api.log.Infof("request to authenticate")

	// Create session to check state validity
	session := sessions.Default(ctx)
	queryState := ctx.Query("state")
	sessionState := session.Get("state")

	api.log.Debugf("  query state: %s", queryState)
	api.log.Debugf("session state: %v", sessionState)
	if queryState != sessionState {
		api.check(fmt.Errorf("invalid session state: %v", sessionState), ctx)
		return
	}
	api.log.Infof("session is valid")

	// Handle the exchange code to initiate a transport
	token, err := api.oauthConfig.Exchange(
		oauth2.NoContext,
		ctx.Query("code"),
	)
	if api.check(err, ctx) {
		return
	}
	api.log.Debugf("transport initiated")

	// The following client logic will also happen in the postgres db lookups,
	// making calls on behalf of the user to obtain the "sub"
	// (primary key in the postgres database)

	// Construct the client
	client := api.oauthConfig.Client(oauth2.NoContext, token)

	// Query the Google API to get information about the user
	// Streamline this to use the api.oauthConfig.Scopes field
	userinfoReq, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if api.check(err, ctx) {
		return
	}
	defer userinfoReq.Body.Close()

	// Read that information
	userinfo, err := ioutil.ReadAll(userinfoReq.Body)
	if api.check(err, ctx) {
		return
	}
	api.log.Infof("got client data %s", string(userinfo))

	// Parse client data
	u := user{}
	err = json.Unmarshal(userinfo, &u)
	if api.check(err, ctx) {
		return
	}
	api.log.Debugf("parsed client %v", u)

	// next: clean up the client.(anything) code into a different func.
	// Add the postgress database tokenizing stuff

	// Replace the next few lines with adding uid, email, and token to a postgress database.
	// Put the exchange token in the cookie

	// session.Set("exchange_token", token)
	// err = session.Save()
	// if api.check(err, ctx) {
	// 	return
	// }

	// ctx.Status(http.StatusOK) // Do this
	ctx.JSON(http.StatusOK, u) // Don't do this
}
