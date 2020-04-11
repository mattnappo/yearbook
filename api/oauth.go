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

// initializeOAuth configures the API's OAuth2 config.
func (api *API) initializeOAuth() {
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

func (api *API) initializeOAuthRoutes() {
	api.router.GET(path.Join(api.oauthRoot, "login"), api.login)
	api.router.GET(path.Join(api.oauthRoot, "auth"), api.auth)

	api.log.Infof("initialized API server OAuth2 routes")

}

func (api *API) login(ctx *gin.Context) {
	state := crypto.GenRandomToken()
	session := sessions.Default(ctx)
	session.Set("state", state)
	session.Save()

	ctx.JSON(http.StatusOK, getLoginURL(state))
}

func (api *API) auth(ctx *gin.Context) {
	// Check state validity
	session := sessions.Default(ctx)
	retrievedState := session.Get("state")
	if retrievedState != ctx.Query("state") {
		api.check(fmt.Errorf("invalid session state: %s", retrievedState), ctx)
		return
	}

	// Handle the exchange code to initiate a transport
	token, err := api.oauthConfig.Exchange(
		oauth2.NoContext,
		ctx.Query("code"),
	)
	if api.check(err, ctx) {
		return
	}

	// Construct the client
	client := api.oauthConfig.Client(oauth2.NoContext, token)

	// Get information about the user
	res, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if api.check(err, ctx) {
		return
	}

	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	if api.check(err, ctx) {
		return
	}
}
