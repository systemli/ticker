package response

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSuccessResponse(t *testing.T) {
	d := []string{"value1", "value2"}
	r := SuccessResponse(map[string]interface{}{"user": d})

	assert.Equal(t, StatusSuccess, r.Status)
	assert.Equal(t, Data{"user": d}, r.Data)
	assert.Equal(t, Error{}, r.Error)
}

func TestErrorResponse(t *testing.T) {
	r := ErrorResponse(CodeDefault, InsufficientPermissions)

	assert.Equal(t, StatusError, r.Status)
	assert.Equal(t, Error{Code: CodeDefault, Message: InsufficientPermissions}, r.Error)
}
