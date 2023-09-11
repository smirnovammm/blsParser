package src

import (
	"fmt"
	tgbotapi "github.com/Syfaro/telegram-bot-api"
)

type TgClient struct {
	bot *tgbotapi.BotAPI
}

func NewTgClient(token string) *TgClient {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		panic(err)
	}
	return &TgClient{bot: bot}
}

func (s TgClient) SendToTg(message string) error {
	msg := tgbotapi.NewMessage(246073487, message)
	if _, err := s.bot.Send(msg); err != nil {
		return fmt.Errorf("Unable to send message to tg: ", err)
	}
	return nil
}
