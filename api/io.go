package api

import (
	"encoding/json"
	"github.com/mattnappo/yearbook/models"
)

// createUserRequest is the structure of a request to create a new user.
type createUserRequest struct {
	Email string `json:"email"`
	Grade string `json:"grade"`
}

type createPostRequest struct {
	Sender     string   `json:"sender"`     // In the form first.last
	Recipients []string `json:"recipients"` // In the form first.last
	Message    string   `json:"message"`    // Just a regular string
	Images     []string `json:"images"`     // Slice of images in base64
}

type authorizeRequest struct {
	Code  string `json:"code"`
	State string `json:"state"`
}

type updateUserRequest struct {
	Username string `json:"username"` // For authentication purposes only
	Grade    int    `json:"grade"`
	Bio      string `json:"bio"`
	Will     string `json:"will"`
	// ProfilePic string `json:"profile_pic"` // Base64 image string
	Nickname string `json:"nickname"`
}

// genericResponse represents a generic response.
type genericResponse struct {
	Data   interface{} `json:"data"`
	Errors []string    `json:"errors"`
}

// loginResponse is the response of a request to /login.
type loginResponse struct {
	Link  string `json:"link"`
	State string `json:"state"`
}

// activityResponse is the response of a request to get a user's activity.
type activityResponse struct {
	Activity    []models.Post `json:"activity"`     // The activity posts
	ProfilePics []string      `json:"profile_pics"` // The profile pics of the senders
}

// inboundOutboundResponse is the response of a getUserPosts/:username request
type inboundOutboundResponse struct {
	Inbound  []models.Post `json:"inbound"`
	Outbound []models.Post `json:"outbound"`
}

// gr constructs a new genericResponse.
func gr(data interface{}, errors ...string) genericResponse {
	return genericResponse{data, errors}
}

// ok returns an OK generic response.
func ok() genericResponse {
	return genericResponse{"", nil}
}

// String marshals a updateUserRequest as a string.
func (uur updateUserRequest) String() string {
	json, _ := json.MarshalIndent(&uur, " ", "  ")
	return string(json)
}
