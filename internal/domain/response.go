package domain

type Response struct {
	Message string `json:"message,omitempty"`
	Data    any    `json:"data",omitempty`
}

type ErrorResponse struct {
	Message string `json:"message,omitempty"`
	Errors  any    `json:"errors",omitempty`
}
