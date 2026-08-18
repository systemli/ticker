package response

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/systemli/ticker/internal/storage"
)

type MessagesResponseTestSuite struct {
	suite.Suite
}

func (s *MessagesResponseTestSuite) TestMessagesResponse() {
	message := storage.NewMessage()
	message.Attachments = []storage.Attachment{{UUID: "uuid", Extension: "jpg"}}

	response := MessagesResponse([]storage.Message{message})

	s.Equal(1, len(response))
	s.Empty(response[0].TelegramURL)
	s.Empty(response[0].MastodonURL)
	s.Empty(response[0].BlueskyURL)
	s.Equal(1, len(response[0].Attachments))

	attachments := response[0].Attachments

	s.Equal("/api/media/uuid.jpg", attachments[0].URL)
}

func TestMessagesResponseTestSuite(t *testing.T) {
	suite.Run(t, new(MessagesResponseTestSuite))
}
