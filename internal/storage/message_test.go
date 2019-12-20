package storage_test

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/systemli/ticker/internal/model"
	"github.com/systemli/ticker/internal/storage"
	"github.com/systemli/ticker/internal/util"
)

func TestFindByTicker(t *testing.T) {
	setup()

	ticker := &model.Ticker{
		ID:          1,
		Active:      true,
		Title:       "Demoticker",
		Description: "Description",
		Domain:      "demoticker.org",
	}

	_ = storage.DB.Save(&ticker)

	c := createContext("")
	pagination := util.NewPagination(&c)
	messages, err := storage.FindByTicker(ticker, pagination)
	if err != nil {
		t.Fail()
	}

	assert.Equal(t, len(messages), 0)

	m1 := model.NewMessage()
	m1.Ticker = ticker.ID
	m1.Text = "First Message"

	err = storage.DB.Save(m1)

	messages, err = storage.FindByTicker(ticker, pagination)
	if err != nil {
		t.Fail()
	}

	assert.Equal(t, len(messages), 1)

	after := m1.ID
	c = createContext(fmt.Sprintf(`after=%d`, after))
	pagination = util.NewPagination(&c)

	messages, err = storage.FindByTicker(ticker, pagination)
	if err != nil {
		t.Fail()
	}

	assert.Equal(t, len(messages), 0)

	before := m1.ID
	c = createContext(fmt.Sprintf(`before=%d`, before))
	pagination = util.NewPagination(&c)

	messages, err = storage.FindByTicker(ticker, pagination)
	if err != nil {
		t.Fail()
	}

	assert.Equal(t, len(messages), 0)

	m2 := model.NewMessage()
	m2.Ticker = ticker.ID
	m2.Text = "Second Message"

	err = storage.DB.Save(m2)

	c = createContext("")
	pagination = util.NewPagination(&c)

	messages, err = storage.FindByTicker(ticker, pagination)
	if err != nil {
		t.Fail()
	}

	assert.Equal(t, len(messages), 2)

	c = createContext(fmt.Sprintf(`before=%d`, m2.ID))
	pagination = util.NewPagination(&c)

	messages, err = storage.FindByTicker(ticker, pagination)
	if err != nil {
		t.Fail()
	}

	assert.Equal(t, len(messages), 1)
	assert.Equal(t, messages[0].ID, 1)
	assert.Equal(t, messages[0].Text, "First Message")

	c = createContext(fmt.Sprintf(`after=%d`, m1.ID))
	pagination = util.NewPagination(&c)

	messages, err = storage.FindByTicker(ticker, pagination)
	if err != nil {
		t.Fail()
	}

	assert.Equal(t, len(messages), 1)
	assert.Equal(t, messages[0].ID, 2)
	assert.Equal(t, messages[0].Text, "Second Message")
}

func TestFindByTickerInactive(t *testing.T) {
	setup()

	ticker := &model.Ticker{
		ID:     1,
		Active: false,
	}

	_ = storage.DB.Save(&ticker)

	c := createContext("")
	pagination := util.NewPagination(&c)
	messages, err := storage.FindByTicker(ticker, pagination)
	if err != nil {
		t.Fail()
	}

	assert.Equal(t, len(messages), 0)
}

func TestDeleteMessage(t *testing.T) {
	setup()

	ticker, message, _ := initialMessageTestData(t)

	err := storage.DeleteMessage(ticker, message)
	if err != nil {
		t.Fail()
	}

	var m *model.Message
	err = storage.DB.Find("ID", message.ID, &m)
	if err == nil {
		t.Fail()
	}
}

func TestDeleteMessageNonExisting(t *testing.T) {
	setup()

	ticker, message, upload := initialMessageTestData(t)

	err := storage.DeleteMessage(ticker, model.NewMessage())
	if err == nil {
		t.Fail()
	}

	_ = storage.DB.DeleteStruct(upload)

	err = storage.DeleteMessage(ticker, message)
	if err == nil {
		t.Fail()
	}
}

func TestDeleteMessages(t *testing.T) {
	setup()

	ticker, message, _ := initialMessageTestData(t)

	err := storage.DeleteMessages(ticker)
	if err != nil {
		t.Fail()
	}

	var m *model.Message
	err = storage.DB.Find("ID", message.ID, &m)
	if err == nil {
		t.Fail()
	}
}

func TestDeleteMessagesNonExisting(t *testing.T) {
	setup()

	ticker, message, _ := initialMessageTestData(t)

	_ = storage.DB.DeleteStruct(message)

	err := storage.DeleteMessages(ticker)
	if err == nil {
		t.Fail()
	}
}

func initialMessageTestData(t *testing.T) (*model.Ticker, *model.Message, *model.Upload) {
	ticker := model.NewTicker()
	err := storage.DB.Save(ticker)
	if err != nil {
		t.Fail()
	}

	upload := model.NewUpload("name.jpg", "image/jpeg", ticker.ID)
	err = storage.DB.Save(upload)
	if err != nil {
		t.Fail()
	}

	message := model.NewMessage()
	attachment := model.Attachment{UUID: upload.UUID, Extension: upload.Extension, ContentType: upload.Extension}
	message.Attachments = []model.Attachment{attachment}
	err = storage.DB.Save(message)
	if err != nil {
		t.Fail()
	}

	return ticker, message, upload
}

func createContext(query string) gin.Context {
	req := http.Request{
		URL: &url.URL{
			RawQuery: query,
		},
	}

	return gin.Context{Request: &req}
}

func setup() {
	gin.SetMode(gin.TestMode)

	model.Config = model.NewConfig()
	model.Config.UploadPath = os.TempDir()

	if storage.DB == nil {
		storage.DB = storage.OpenDB(fmt.Sprintf("%s/ticker_%d.db", os.TempDir(), time.Now().Nanosecond()))
	}
	_ = storage.DB.Drop("Ticker")
	_ = storage.DB.Drop("Message")
	_ = storage.DB.Drop("Upload")
	_ = storage.DB.Drop("User")
	_ = storage.DB.Drop("Setting")
}
