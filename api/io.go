package api

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
