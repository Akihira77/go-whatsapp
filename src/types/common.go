package types

type ErrorResponse struct {
	Field string `json:"field"`
	Error string `json:"error"`
}
