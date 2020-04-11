package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"path"
	"strconv"
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

// getUserInfo querys the Google API to get user info given a token.
func (api *API) getUserInfo(token *oauth2.Token) (user, error) {
	api.log.Infof("querying Google API to get userinfo")
	client := api.oauthConfig.Client(oauth2.NoContext, token)

	// Query the Google API to get information about the user
	userinfoReq, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		return user{}, err
	}
	defer userinfoReq.Body.Close()

	// Read that information
	userinfo, err := ioutil.ReadAll(userinfoReq.Body)
	if err != nil {
		return user{}, err
	}
	api.log.Infof("got client data %s", string(userinfo))

	// Parse client data
	u := user{}
	err = json.Unmarshal(userinfo, &u)
	if err != nil {
		return user{}, err
	}
	api.log.Debugf("parsed client info %v", u)
	return u, nil
}

// authorizeRequest is used to authorize a request for a certain
// endpoint group.
func (api *API) authorizeRequest() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := strings.Split(ctx.GetHeader("Authorization"), " ")
		// Check that the authorization header exists
		if len(authHeader) <= 0 {
			api.check(errors.New("no authorization header"), ctx)
			return
		}
		// Get the token from the header and put it in an oauth2.Token
		tokenString := authHeader[1]
		token := oauth2.Token{AccessToken: tokenString}
		api.log.Infof("attempting to authorize request with token %s",
			tokenString,
		)

		// Get user info to obtain the sub
		u, err := api.getUserInfo(&token)
		if api.check(err, ctx) {
			return
		}
		sub, err := strconv.ParseInt(u.Sub, 2, 64)
		if api.check(err, ctx) {
			return
		}

		// Query the database to get the token given the sub
		token, err := api.database.GetToken(sub)

		// Construct the client to lookup the sub associated with the token
		// given in the authentication bearer header

		// Check if the token provided in the authorization header equals the
		// token from the database.
		if tokenString != token {
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

	// will return just url: react app will query for url and render button
	// on react side (client side)
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

	// Construct the client
	u, err := api.getUserInfo(token)
	if api.check(err, ctx) {
		return
	}

	// Insert the token into the database
	sub, err := strconv.ParseInt(u.Sub, 2, 64)
	err = api.database.InsertToken(sub, token, u.Email)
	if api.check(err, ctx) {
		return
	}

	// Store the OAuth2 exchaneg token in a cookie
	session.Set("google_oauth2_token", token)
	err = session.Save()
	if api.check(err, ctx) {
		return
	}

	// ctx.Status(http.StatusOK) // Do this
	ctx.JSON(http.StatusOK, u) // Don't do this
}
