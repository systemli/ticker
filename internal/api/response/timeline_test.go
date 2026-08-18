package response

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/systemli/ticker/internal/storage"
)

type TimelineTestSuite struct {
	suite.Suite
}

func (s *TimelineTestSuite) TestTimelineResponse() {
	message := storage.NewMessage()
	message.Attachments = []storage.Attachment{{UUID: "uuid", Extension: "jpg"}}

	response := TimelineResponse([]storage.Message{message})

	s.Equal(1, len(response))
	s.Equal(1, len(response[0].Attachments))

	attachments := response[0].Attachments

	s.Equal("/api/media/uuid.jpg", attachments[0].URL)
}

func TestTimelineTestSuite(t *testing.T) {
	suite.Run(t, new(TimelineTestSuite))
}
