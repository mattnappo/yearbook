package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/xoreo/yearbook/models"
)

// createPost creates a new post.
func (api *API) createPost(ctx *gin.Context) {
	api.log.Infof("request to create post")
	// Decode the post request
	var request createPostRequest
	err := ctx.ShouldBindJSON(&request)
	if api.check(err, ctx, http.StatusBadRequest) {
		return
	}

	api.log.Infof("\n\nrequest to create post: %v", request)

	// Check that the sender in the request is the same as the username
	// associated with that username's token
	err = api.authenticate(ctx, request.Sender)
	if api.check(err, ctx, http.StatusUnauthorized) {
		return
	}

	// Create the new post
	post, err := models.NewPost(
		request.Sender,
		request.Message,
		request.Images,
		request.Recipients,
	)
	if api.check(err, ctx) {
		return
	}

	api.log.Debugf("constructed new post %s", post.String())

	// Add it to the database
	err = api.database.AddPost(post)
	if api.check(err, ctx) {
		return
	}

	// Add the recipients to the database (if they do not already exist)
	for _, recip := range post.Recipients {
		newUser, err := models.NewUser(recip.Email(), models.Senior, false)
		if api.check(err, ctx) {
			return
		}
		err = api.database.AddUser(newUser) // Unhandled err
	}

	// Add to and from post to user data
	err = api.database.AddToAndFrom(
		post.PostID,
		string(post.Sender),
		request.Recipients,
	)
	if api.check(err, ctx) {
		return
	}

	api.log.Infof("created new post %s", post.String())
	ctx.JSON(http.StatusOK, ok())
}

// getPost gets a post.
func (api *API) getPost(ctx *gin.Context) {
	id := ctx.Param("id")
	api.log.Infof("request to get post %s", id)

	post, err := api.database.GetPost(id)
	if api.check(err, ctx) {
		return
	}

	api.log.Infof("got post %s from database", post.String())
	ctx.JSON(http.StatusOK, gr(post))
}

// getPosts gets all posts.
func (api *API) getPosts(ctx *gin.Context) {
	api.log.Infof("request to get all posts")

	posts, err := api.database.GetAllPosts()
	if api.check(err, ctx) {
		return
	}

	api.log.Infof("got all posts")
	ctx.JSON(http.StatusOK, gr(posts))
}

// getnPosts gets n posts.
func (api *API) getnPosts(ctx *gin.Context) {
	n := ctx.Param("n")
	api.log.Infof("request to get %s posts", n)

	nInt, err := strconv.Atoi(n)
	if api.check(err, ctx) {
		return
	}

	posts, err := api.database.GetnPosts(nInt)
	if api.check(err, ctx) {
		return
	}

	api.log.Infof("got %d posts", nInt)
	ctx.JSON(http.StatusOK, gr(posts))
}

// deletePost deletes a post.
func (api *API) deletePost(ctx *gin.Context) {
	postID := ctx.Param("id")
	api.log.Infof("request to delete post %s", postID)

	err := api.database.DeletePost(postID)
	if api.check(err, ctx) {
		return
	}

	api.log.Infof("deleted post %s", postID)
	ctx.JSON(http.StatusOK, ok())
}

// createUser creates a user.
func (api *API) createUser(ctx *gin.Context) {
	api.log.Infof("request to create user")

	// Decode the new user request
	var request createUserRequest
	err := ctx.ShouldBindJSON(&request)
	if api.check(err, ctx, http.StatusBadRequest) {
		return
	}

	// Determine the grade from the request
	var grade models.Grade
	switch request.Grade {
	case "freshman":
		grade = models.Freshman
		break
	case "sophomore":
		grade = models.Sophomore
		break
	case "junior":
		grade = models.Junior
		break
	case "senior":
		grade = models.Senior
		break
	default:
		api.check(fmt.Errorf("invalid grade %v", request.Grade), ctx)
		return
	}

	// Create the new user
	user, err := models.NewUser(request.Email, grade, true)
	if api.check(err, ctx) {
		return
	}
	api.log.Debugf("constructed new user %s", user.String())

	// Add it to the database
	err = api.database.AddUser(user)
	if api.check(err, ctx) {
		return
	}

	api.log.Infof("created new user %s", user.String())
	ctx.JSON(http.StatusOK, ok())
}

// updateUser handles a request to update a user.
func (api *API) updateUser(ctx *gin.Context) {
	api.log.Infof("request to update user")

	api.log.Infof("\n\n%v\n\n", ctx.Request.Body)

	// Decode the request data
	var request updateUserRequest
	err := ctx.ShouldBindJSON(&request)
	if api.check(err, ctx, http.StatusBadRequest) {
		return
	}
	api.log.Debugf("updating user with new info %v", request)

	// Construct a user struct with the new user data in it, and everything
	// else blank.
	newUserData, err := models.UserFromString(request.String())
	if api.check(err, ctx) {
		return
	}

	// Check that the username of the request is the same as the username
	// of the account attempting to be modified via this request.
	err = api.authenticate(ctx, string(newUserData.Username))
	if api.check(err, ctx, http.StatusUnauthorized) {
		return
	}

	api.log.Infof("authenticated %s", string(newUserData.Username))

	// Update the user in the database
	err = api.database.UpdateUser(newUserData)
	if api.check(err, ctx) {
		return
	}

	api.log.Infof("updated user %s", string(newUserData.Username))
	ctx.JSON(http.StatusOK, ok())
}

// getUser gets a user.
func (api *API) getUser(ctx *gin.Context) {
	username := ctx.Param("username")
	api.log.Infof("request to get user %s", username)

	user, err := api.database.GetUser(username)
	if api.check(err, ctx) {
		return
	}

	api.log.Infof("got user %s", user.String())
	ctx.JSON(http.StatusOK, gr(user))
}

// getUserProfilePic gets the profile picture of a user.
func (api *API) getUserProfilePic(ctx *gin.Context) {
	username := ctx.Param("username")
	api.log.Infof("request to get user profile pic %s", username)

	profilePic, err := api.database.GetUserProfilePic(username)
	if api.check(err, ctx) {
		return
	}

	api.log.Infof("got user profile pic for %s", username)
	ctx.JSON(http.StatusOK, gr(profilePic))
}

// getUserWithAuth handles a request to get a user
// (with authentication).
func (api *API) getUserWithAuth(ctx *gin.Context) {
	api.log.Infof("request to get user with authentication")

	// Check that the username of the request is the same as the
	// username
	// behind the token
	username := ctx.Param("username")
	err := api.authenticate(ctx, username)
	if api.check(err, ctx, http.StatusUnauthorized) {
		return
	}

	api.log.Infof("authenticated %s", username)

	// Get the user from the database
	user, err := api.database.GetUser(username)
	if api.check(err, ctx) {
		return
	}

	api.log.Infof("got user %s", username)
	ctx.JSON(http.StatusOK, gr(user))
}

// getUsers gets all users.
func (api *API) getUsers(ctx *gin.Context) {
	api.log.Infof("request to get all users")

	users, err := api.database.GetAllUsers()
	if api.check(err, ctx) {
		return
	}

	api.log.Infof("got all users")
	ctx.JSON(http.StatusOK, gr(users))
}

// getUsernames gets all usernames in the database.
func (api *API) getUsernames(ctx *gin.Context) {
	api.log.Infof("request to get all usernames")

	usernames, err := api.database.GetAllUsernames()
	if api.check(err, ctx) {
		return
	}

	api.log.Infof("got all usernames")
	ctx.JSON(http.StatusOK, gr(usernames))
}

// getSeniors gets all senior usernames.
func (api *API) getSeniors(ctx *gin.Context) {
	api.log.Infof("request to get all senior usernames")

	usernames, err := api.database.GetAllSeniorUsernames()
	if api.check(err, ctx) {
		return
	}

	api.log.Infof("got all senior usernames")
	ctx.JSON(http.StatusOK, gr(usernames))
}

// deleteUser deletes a user.
func (api *API) deleteUser(ctx *gin.Context) {
	username := ctx.Param("username")
	api.log.Infof("request to delete user %s", username)

	err := api.database.DeleteUser(username)
	if api.check(err, ctx) {
		return
	}

	api.log.Infof("deleted user %s", username)
	ctx.JSON(http.StatusOK, ok())
}

// getActivity gets a list of the recent posts about a user.
func (api *API) getActivity(ctx *gin.Context) {
	username := ctx.Param("username")
	api.log.Infof("request to get activity for %s", username)

	// Authenticate the user
	err := api.authenticate(ctx, username)
	if api.check(err, ctx, http.StatusUnauthorized) {
		return
	}

	api.log.Infof("authenticated %s", username)

	inboundPosts, err := api.database.GetUserInbound(username)
	if api.check(err, ctx) {
		return
	}

	api.log.Infof("got %s's inbound posts", username)
	ctx.JSON(http.StatusOK, gr(inboundPosts))
}

// getUserPosts gets the inbound and outbound posts of a user.
func (api *API) getUserPosts(ctx *gin.Context) {
	username := ctx.Param("username")
	api.log.Infof("request to get user posts for user %s", username)

	posts, err := api.database.GetUserInboundOutbound(username)
	if api.check(err, ctx) {
		return
	}

	api.log.Infof("got %s's inbound and outbound posts", username)
	res := inboundOutboundResponse{
		Inbound:  posts[0],
		Outbound: posts[1],
	}
	ctx.JSON(http.StatusOK, gr(res))
}
