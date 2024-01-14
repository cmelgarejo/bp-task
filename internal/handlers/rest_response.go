package handlers

// GenericResponse is a generic response for the API (currently one use, handleIPFSPost)
type GenericResponse struct {
	Message string `json:"message"`
}

// NewResponse creates a new GenericResponse
func NewResponse(message string) *GenericResponse {
	return &GenericResponse{
		Message: message,
	}
}
