package api

// createUserRequest is the structure of a request to create a new user.
type createUserRequest struct {
	Email string `json:"email"`
}

type createPostRequest struct {
	Sender     string   `json:"sender"`
	Recipients []string `json:"recipients"`
	Message    string   `json:"message"`
	Images     [][]byte `json:"images"`
}

// genericResponse represents a generic response.
type genericResponse struct {
	Data   string   `json:"data"`
	Errors []string `json:"errors"`
}
