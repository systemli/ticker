package response

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/systemli/ticker/internal/config"
	"github.com/systemli/ticker/internal/storage"
)

var (
	u = storage.NewUpload("image.jpg", "image/jpg", 1)
	c = config.Config{
		Upload: config.Upload{URL: "http://localhost:8080"},
	}
)

type UploadResponseTestSuite struct {
	suite.Suite
}

func (s *UploadResponseTestSuite) TestUploadResponse() {
	response := UploadResponse(u, c)

	s.Equal(fmt.Sprintf("%s/media/%s", c.Upload.URL, u.FileName()), response.URL)
}

func (s *UploadResponseTestSuite) TestUploadsResponse() {
	response := UploadsResponse([]storage.Upload{u}, c)

	s.Equal(1, len(response))
}

func TestUploadResponseTestSuite(t *testing.T) {
	suite.Run(t, new(UploadResponseTestSuite))
}
