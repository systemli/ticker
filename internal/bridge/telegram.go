package bridge

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/systemli/ticker/internal/config"
	"github.com/systemli/ticker/internal/storage"
)

type TelegramBridge struct {
	config  config.Config
	storage storage.TickerStorage
}

func (tb *TelegramBridge) Send(ticker storage.Ticker, message storage.Message) error {
	if ticker.Telegram.ChannelName == "" || !tb.config.TelegramEnabled() {
		return nil
	}

	bot, err := tgbotapi.NewBotAPI(tb.config.TelegramBotToken)
	if err != nil {
		return err
	}

	if len(message.Attachments) == 0 {
		msgConfig := tgbotapi.NewMessageToChannel(ticker.Telegram.ChannelName, message.Text)
		msg, err := bot.Send(msgConfig)
		if err != nil {
			return err
		}
		message.Telegram = storage.TelegramMeta{Messages: []tgbotapi.Message{msg}}
	} else {
		var photos []interface{}
		for _, attachment := range message.Attachments {
			upload, err := tb.storage.FindUploadByUUID(attachment.UUID)
			if err != nil {
				log.WithError(err).Error("failed to find upload")
				continue
			}

			media := tgbotapi.FilePath(upload.FullPath(tb.config.UploadPath))
			if upload.ContentType == "image/gif" {
				photo := tgbotapi.NewInputMediaDocument(media)
				photo.Caption = message.Text
				photos = append(photos, photo)
			} else {
				photo := tgbotapi.NewInputMediaPhoto(media)
				photo.Caption = message.Text
				photos = append(photos, photo)
			}
		}

		mediaGroup := tgbotapi.MediaGroupConfig{
			ChannelUsername: ticker.Telegram.ChannelName,
			Media:           photos,
		}

		msgs, err := bot.SendMediaGroup(mediaGroup)
		if err != nil {
			return err
		}
		message.Telegram = storage.TelegramMeta{Messages: msgs}
	}

	return nil
}

func (tb *TelegramBridge) Delete(ticker storage.Ticker, message storage.Message) error {
	if ticker.Telegram.ChannelName == "" || !tb.config.TelegramEnabled() {
		return nil
	}

	if len(message.Telegram.Messages) == 0 {
		return nil
	}

	bot, err := tgbotapi.NewBotAPI(tb.config.TelegramBotToken)
	if err != nil {
		return err
	}

	for _, message := range message.Telegram.Messages {
		deleteMessageConfig := tgbotapi.DeleteMessageConfig{MessageID: message.MessageID, ChatID: message.Chat.ID}
		_, err = bot.Request(deleteMessageConfig)
		if err != nil {
			log.WithError(err).Error("failed to delete telegram message")
			continue
		}
	}

	return nil
}

func BotUser(token string) (tgbotapi.User, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return tgbotapi.User{}, err
	}

	user, err := bot.GetMe()
	if err != nil {
		return tgbotapi.User{}, err
	}

	return user, nil
}
