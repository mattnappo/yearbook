package api

// createUserRequest is the structure of a request to create a new user.
type createUserRequest struct {
	Email string `json:"email"`
}

type createPostRequest struct {
	Sender     string   `json:"sender"`     // In the form first.last
	Recipients []string `json:"recipients"` // In the form first.last
	Message    string   `json:"message"`    // Just a regular string
	Images     []string `json:"images"`     // Slice of images in base64
}

// genericResponse represents a generic response.
type genericResponse struct {
	Data   interface{} `json:"data"`
	Errors []string    `json:"errors"`
}

// gr constructs a new genericResponse.
func gr(data interface{}, errors ...string) genericResponse {
	return genericResponse{data, errors}
}

// ok returns an OK generic response.
func ok() genericResponse {
	return genericResponse{"", nil}
}
