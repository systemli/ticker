package api

type errorResponse struct {
	Data   interface{} `json:"data"`
	Status string      `json:"status"`
	Error  errorData   `json:"error"`
}

type errorData struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}
