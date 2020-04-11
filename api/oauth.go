package api

import (
	"github.com/xoreo/yearbook/common"
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

func (api *API) initializeOAuthRoutes() {

}

// getLoginURL gets the "Sign in with Google" URL.
func (api *API) getLoginURL(state string) string {
	return api.oauthConfig.AuthCodeURL(state)
}
