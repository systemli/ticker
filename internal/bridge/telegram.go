package bridge

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/sirupsen/logrus"
	"github.com/systemli/ticker/internal/model"
	"github.com/systemli/ticker/internal/storage"
)

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

func SendTelegramMessage(ticker *model.Ticker, message *model.Message) error {
	if ticker.Telegram.ChannelName == "" || !model.Config.TelegramBotEnabled() {
		return nil
	}

	bot, err := tgbotapi.NewBotAPI(model.Config.TelegramBotToken)
	if err != nil {
		return err
	}

	if len(message.Attachments) == 0 {
		msgConfig := tgbotapi.NewMessageToChannel(ticker.Telegram.ChannelName, message.Text)
		msg, err := bot.Send(msgConfig)
		if err != nil {
			return err
		}
		message.Telegram = model.TelegramMeta{Messages: []tgbotapi.Message{msg}}
	} else {
		var photos []interface{}
		for _, attachment := range message.Attachments {
			upload := &model.Upload{}
			err := storage.DB.One("UUID", attachment.UUID, upload)
			if err != nil {
				log.WithError(err).Error("failed to find upload")
				continue
			}

			media := tgbotapi.FilePath(upload.FullPath())
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
		message.Telegram = model.TelegramMeta{Messages: msgs}
	}

	return nil
}

func DeleteTelegramMessage(ticker *model.Ticker, message *model.Message) error {
	if ticker.Telegram.ChannelName == "" || !model.Config.TelegramBotEnabled() {
		return nil
	}

	if len(message.Telegram.Messages) == 0 {
		return nil
	}

	bot, err := tgbotapi.NewBotAPI(model.Config.TelegramBotToken)
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
