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
	"github.com/mattnappo/yearbook/common"
	"github.com/mattnappo/yearbook/crypto"
	"github.com/mattnappo/yearbook/models"
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
	api.router.GET(path.Join(api.oauthRoot, "login"), api.login)
	api.router.POST(path.Join(api.oauthRoot, "authorize"), api.authorize)

	api.log.Infof("initialized API server OAuth2 routes")

}

// getLoginURL gets the "Sign in with Google" URL.
func (api *API) getLoginURL(state string) string {
	return api.oauthConfig.AuthCodeURL(state)
}

// getUserInfo querys the Google API to get user info given a token.
func (api *API) getUserInfo(token *oauth2.Token) (user, error) {
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

	// Parse client data
	u := user{}
	err = json.Unmarshal(userinfo, &u)
	if err != nil {
		return user{}, err
	}

	// Check that there were no errors
	if u.Sub == "" {
		return user{}, errors.New("invalid credentials to query Google API")
	}

	return u, nil
}

// extractBearerToken extracts the bearer token from a *gin.Context.
func extractBearerToken(ctx *gin.Context) (string, error) {
	authHeader := strings.Split(ctx.GetHeader("Authorization"), "bearer")
	// Check that the authorization header exists
	if len(authHeader) <= 1 {
		return "", errors.New("no authorization header")
	}
	// Remove all spaces from the parsed auth header
	return strings.ReplaceAll(authHeader[1], " ", ""), nil
}

// authorizeRequest is the middleware used to authorize a request for a
// certain endpoint group.
func (api *API) authorizeRequest() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Extract the bearer token
		headerTokenString, err := extractBearerToken(ctx)
		if api.check(err, ctx, http.StatusUnauthorized) {
			return
		}

		// Construct the actual token which is needed to get the correct
		// token from the database because with this "headerToken" can a
		// connection with the Google API be established. This connection
		// is how we obtain the sub, which is then used to lookup the
		// correct token in the Postgres database.
		headerToken := oauth2.Token{AccessToken: headerTokenString}

		// Get user info to obtain the sub
		u, err := api.getUserInfo(&headerToken)
		if api.check(err, ctx) {
			return
		}
		// Query the database to get the token, given the sub
		correctToken, err := api.database.GetToken(u.Sub)
		if api.check(err, ctx) {
			return
		}

		// Check if the token provided in the authorization header equals the
		// token from the database. If an only if this is true will the user
		// gain authorization to the protected resources.
		if headerTokenString != correctToken {
			api.log.Infof(errUnauthorized.Error())
			api.check(errUnauthorized, ctx, http.StatusUnauthorized)
			return
		}
		ctx.Next()
	}
}

// authenticate authenticates a username and a given token.
func (api *API) authenticate(ctx *gin.Context, username string) error {
	errNoAuthentication := fmt.Errorf("authentication for %s failed", username)

	// Get the token from the Authorization bearer header
	bearerToken, err := extractBearerToken(ctx)
	if err != nil {
		return errNoAuthentication
	}

	// Query the Google API to get the email associated with the token
	// in the request header.
	headerToken := &oauth2.Token{AccessToken: bearerToken}
	googleUser, err := api.getUserInfo(headerToken)
	if err != nil {
		return errNoAuthentication
	}

	// Get the username of the email responded by the Google API
	googleUsername, err := models.UsernameFromEmail(googleUser.Email)
	if string(googleUsername) != username {
		return errNoAuthentication
	}

	return nil
}

// login handles a request to login.
func (api *API) login(ctx *gin.Context) {
	api.log.Infof("request to login")

	// Generate a random token for the cookie session handler
	state := crypto.GenRandomToken()
	http.SetCookie(ctx.Writer, &http.Cookie{
		Name:   "state",
		Value:  state,
		Path:   "/",
		MaxAge: 30 * 60,
		Secure: false,
	})

	// Store the state in a session (still dont know where this is)
	session := sessions.Default(ctx)
	session.Set("state", state)
	err := session.Save()
	if api.check(err, ctx, http.StatusUnauthorized) {
		return
	}

	// Respond with the URL to "Sign in with Google"
	ctx.JSON(http.StatusOK, api.getLoginURL(state))
}

// authorize is the Google authorization callback URL.
func (api *API) authorize(ctx *gin.Context) {
	// Decode the request
	var request authorizeRequest
	err := ctx.ShouldBindJSON(&request)
	if api.check(err, ctx, http.StatusBadRequest) {
		return
	}

	// Get the state from the client-side cookie
	queryState, err := ctx.Request.Cookie("state")
	if api.check(err, ctx, http.StatusUnauthorized) {
		return
	}

	// Create session to check back-end state
	session := sessions.Default(ctx)
	sessionState := session.Get("state")

	// Compare the two states
	if queryState.Value != sessionState { // Compare session and query states
		api.check(fmt.Errorf("invalid query state: %v", queryState.Value),
			ctx, http.StatusUnauthorized)
		return
	}

	// If they have a token already, don't issue a new one
	_, err = ctx.Request.Cookie("token")
	if err == nil {
		// Maybe search PG database to see if its valid
		ctx.JSON(http.StatusOK, ok())
		return
	}

	// Handle the exchange code to initiate a transport and get a token
	token, err := api.oauthConfig.Exchange(
		oauth2.NoContext,
		request.Code,
	)
	if api.check(err, ctx, http.StatusUnauthorized) {
		return
	}

	// Construct the client to get the sub.
	u, err := api.getUserInfo(token)
	if api.check(err, ctx) {
		return
	}

	// Insert the token into the database
	err = api.database.InsertToken(u.Sub, token, u.Email)
	if api.check(err, ctx) {
		return
	}

	// Try and get user from database to determine whether to update,
	// add, or do nothing
	stringUsername, err := models.UsernameFromEmail(u.Email)
	if api.check(err, ctx) {
		return
	}
	dbUser, err := api.database.GetUser(string(stringUsername))
	cookieUsername := dbUser.Username
	if string(dbUser.Username) == "" { // If the user does not exist
		newUser, err := models.NewUser(u.Email, models.Freshman, true)
		if api.check(err, ctx) {
			return
		}
		newUser.ProfilePic = u.Picture // Get and set the profile picture
		api.database.AddUser(newUser)
		cookieUsername = newUser.Username
		api.log.Infof("added user %s to database (in authorization)", u.Email)
	} else { // If the user EXISTS
		if dbUser.Registered == false { // If this is their first login
			// Update the account's profile picture and registration
			// status
			err = api.database.InitAccount(string(stringUsername), u.Picture)
			if api.check(err, ctx) {
				return
			}
		}
	}
	// Set the token in a cookie
	http.SetCookie(ctx.Writer, &http.Cookie{
		Name:   "token",
		Value:  token.AccessToken,
		Path:   "/",
		MaxAge: 30 * 60,
		Secure: false,
	})
	// Set the username in a cookie
	http.SetCookie(ctx.Writer, &http.Cookie{
		Name:   "username",
		Value:  string(cookieUsername),
		Path:   "/",
		MaxAge: 30 * 60,
		Secure: false,
	})

	api.log.Infof("authorized %s", u.Email)

	ctx.JSON(http.StatusOK, ok())
}
