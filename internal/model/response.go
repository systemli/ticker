package model

const (
	ErrorCodeDefault                 = 1000
	ErrorCodeNotFound                = 1001
	ErrorCodeCredentials             = 1002
	ErrorCodeInsufficientPermissions = 1003

	ErrorInsufficientPermissions = "insufficient permissions"
	ErrorUserIdentifierMissing   = "user identifier not found"
	ErrorUserNotFound            = "user not found"
	ErrorTickerNotFound          = "ticker not found"
	ErrorSettingNotFound         = "setting not found"

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
