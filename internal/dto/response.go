package dto

type APIResponse struct {
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

type APIError struct {
	Message string `json:"message"`
	Error   string `json:"error"`
}
