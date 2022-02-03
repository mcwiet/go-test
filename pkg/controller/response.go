package controller

// Standard response format
type Response struct {
	Data  interface{} `json:"data"`
	Error error       `json:"error"`
}
