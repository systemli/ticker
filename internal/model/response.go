package model

const (
	ErrorUnspecified             = 1000
	ErrorNotFound                = 1001
	ErrorCredentials             = 1002
	ErrorInsufficientPermissions = 1003

	ResponseSuccess = `success`
	ResponseError   = `error`
)

//JSONResponse represents response structure
type JSONResponse struct {
	Data   map[string]interface{} `json:"data"`
	Status string                 `json:"status"`
	Error  interface{}            `json:"error"`
}

//NewJSONSuccessResponse returns a successfully response
func NewJSONSuccessResponse(name string, data interface{}) JSONResponse {
	return JSONResponse{
		Data:   map[string]interface{}{name: data},
		Status: ResponseSuccess,
	}
}

//NewJSONErrorResponse returns a erroneous response
func NewJSONErrorResponse(code int, message string) JSONResponse {
	return JSONResponse{
		Data:   map[string]interface{}{},
		Status: ResponseError,
		Error: map[string]interface{}{
			"code":    code,
			"message": message,
		},
	}
}
