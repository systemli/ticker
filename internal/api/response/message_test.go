package response

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/systemli/ticker/internal/config"
	"github.com/systemli/ticker/internal/storage"
)

func TestMessagesResponse(t *testing.T) {
	config := config.Config{UploadURL: "https://upload.example.com"}
	message := storage.NewMessage()
	message.Attachments = []storage.Attachment{{UUID: "uuid", Extension: "jpg"}}

	response := MessagesResponse([]storage.Message{message}, config)

	assert.Equal(t, 1, len(response))
	assert.Empty(t, response[0].TwitterURL)
	assert.Empty(t, response[0].TelegramURL)
	assert.Empty(t, response[0].MastodonURL)
	assert.Equal(t, `{"type":"FeatureCollection","features":[]}`, response[0].GeoInformation)
	assert.Equal(t, 1, len(response[0].Attachments))

	attachments := response[0].Attachments

	assert.Equal(t, "https://upload.example.com/media/uuid.jpg", attachments[0].URL)
}
