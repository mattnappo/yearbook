package api

import "encoding/json"

// createUserRequest is the structure of a request to create a new user.
type createUserRequest struct {
	Email string `json:"email"`
}

type createPostRequest struct {
	Sender     string   `json:"sender"`     // In the form first.last
	Recipients []string `json:"recipients"` // In the form first.last
	Message    string   `json:"message"`
	Images     [][]byte `json:"images"`
}

// genericResponse represents a generic response.
type genericResponse struct {
	Data   string   `json:"data"`
	Errors []string `json:"errors"`
}

// newGenericResponse constructs a new genericResponse as a string.
func newGenericResponse(data string, errors ...string) string {
	gr := genericResponse{data, errors}
	json, _ := json.MarshalIndent(gr, " ", "  ")
	return string(json)
}

// ok returns an OK generic response.
func ok() string {
	return "success"
}
