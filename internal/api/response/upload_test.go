package response

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/systemli/ticker/internal/storage"
)

var u = storage.NewUpload("image/jpeg", 1)

type UploadResponseTestSuite struct {
	suite.Suite
}

func (s *UploadResponseTestSuite) TestUploadResponse() {
	response := UploadResponse(u)

	s.Equal("/api/media/"+u.FileName(), response.URL)
}

func (s *UploadResponseTestSuite) TestUploadsResponse() {
	response := UploadsResponse([]storage.Upload{u})

	s.Equal(1, len(response))
}

func TestUploadResponseTestSuite(t *testing.T) {
	suite.Run(t, new(UploadResponseTestSuite))
}
