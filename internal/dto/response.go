package dto

type SuccessResponse struct {
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func OK(message string, data any) SuccessResponse {
	return SuccessResponse{Message: message, Data: data}
}

func Err(message string) ErrorResponse {
	return ErrorResponse{Error: message}
}
