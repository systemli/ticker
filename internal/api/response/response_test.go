package response

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type ResponseTestSuite struct {
	suite.Suite
}

func (s *ResponseTestSuite) TestResponse() {
	s.Run("when status is success", func() {
		d := []string{"value1", "value2"}
		r := SuccessResponse(map[string]interface{}{"user": d})

		s.Equal(StatusSuccess, r.Status)
		s.Equal(Data(map[string]interface{}{"user": d}), r.Data)
		s.Equal(Error{}, r.Error)
	})

	s.Run("when status is error", func() {
		r := ErrorResponse(CodeDefault, InsufficientPermissions)

		s.Equal(StatusError, r.Status)
		s.Equal(Data(nil), r.Data)
		s.Equal(Error{Code: CodeDefault, Message: InsufficientPermissions}, r.Error)
	})
}

func TestResponseTestSuite(t *testing.T) {
	suite.Run(t, new(ResponseTestSuite))
}
