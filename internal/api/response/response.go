package response

const (
	CodeDefault                 ErrorCode = 1000
	CodeNotFound                ErrorCode = 1001
	CodeBadCredentials          ErrorCode = 1002
	CodeInsufficientPermissions ErrorCode = 1003

	InsufficientPermissions ErrorMessage = "insufficient permissions"
	Unauthorized            ErrorMessage = "unauthorized"
	UserIdentifierMissing   ErrorMessage = "user identifier not found"
	TickerIdentifierMissing ErrorMessage = "ticker identifier not found"
	MessageNotFound         ErrorMessage = "message not found"
	FilesIdentifierMissing  ErrorMessage = "files identifier not found"
	TooMuchFiles            ErrorMessage = "upload limit exceeded"
	UserNotFound            ErrorMessage = "user not found"
	TickerNotFound          ErrorMessage = "ticker not found"
	SettingNotFound         ErrorMessage = "setting not found"
	MessageFetchError       ErrorMessage = "messages couldn't fetched"
	FormError               ErrorMessage = "invalid form values"
	StorageError            ErrorMessage = "failed to save"
	UploadsNotFound         ErrorMessage = "uploads not found"
	MastodonError           ErrorMessage = "unable to connect to mastodon"
	BlueskyError            ErrorMessage = "unable to connect to bluesky"
	SignalGroupError        ErrorMessage = "unable to connect to signal"
	PasswordError           ErrorMessage = "could not authenticate password"

	StatusSuccess Status = `success`
	StatusError   Status = `error`
)

type ErrorCode int
type ErrorMessage string
type Data map[string]interface{}
type Status string

type Response struct {
	Data   Data   `json:"data" swaggertype:"object,string"`
	Status Status `json:"status"`
	Error  Error  `json:"error,omitempty"`
}

type Error struct {
	Code    ErrorCode    `json:"code,omitempty"`
	Message ErrorMessage `json:"message,omitempty"`
}

func SuccessResponse(data map[string]interface{}) Response {
	return Response{
		Data:   data,
		Status: StatusSuccess,
	}
}

func ErrorResponse(code ErrorCode, message ErrorMessage) Response {
	return Response{
		Error: Error{
			Code:    code,
			Message: message,
		},
		Status: StatusError,
	}
}
