package api

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"path"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/xoreo/yearbook/common"
	"github.com/xoreo/yearbook/crypto"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

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
		RedirectURL:  "http://localhost:8080/oauth/callback",
		ClientID:     common.GetEnv("GOOGLE_CLIENT_ID"),
		ClientSecret: common.GetEnv("GOOGLE_CLIENT_SECRET"),
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}

	api.initializeOAuthRoutes()
}

// getLoginURL gets the "Sign in with Google" URL.
func (api *API) getLoginURL(state string) string {
	return api.oauthConfig.AuthCodeURL(state)
}

// initializeOAuthRoutes initializes the OAuth2-related API routes.
func (api *API) initializeOAuthRoutes() {
	api.router.GET(path.Join(api.oauthRoot, "login"), api.login)
	api.router.GET(path.Join(api.oauthRoot, "auth"), api.auth)

	api.log.Infof("initialized API server OAuth2 routes")

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

	// Respond with the URL to "Sign in with Google"
	ctx.JSON(http.StatusOK, api.getLoginURL(state))
}

// auth handles a request to authenticate.
func (api *API) auth(ctx *gin.Context) {
	api.log.Infof("request to authenticate")

	// Check state validity
	session := sessions.Default(ctx)
	retrievedState := session.Get("state")
	if retrievedState != ctx.Query("state") {
		api.check(fmt.Errorf("invalid session state: %s", retrievedState), ctx)
		return
	}
	api.log.Debugf("session is valid")

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
	client := api.oauthConfig.Client(oauth2.NoContext, token)

	// Query the Google API to get information about the user
	userinfo, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if api.check(err, ctx) {
		return
	}
	defer userinfo.Body.Close()

	// Read that information
	data, err := ioutil.ReadAll(userinfo.Body)
	if api.check(err, ctx) {
		return
	}
	api.log.Infof("got client data %s", string(data))

	session.Set("user-id", data)

	ctx.Status(http.StatusOK)
}