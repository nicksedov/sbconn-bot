package telegram

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/nicksedov/sbconn-bot/pkg/cli"
)

var bot *tgbotapi.BotAPI

func getBotAPI() (*tgbotapi.BotAPI, error) {
	if bot == nil {
		var err error
		bot, err = tgbotapi.NewBotAPI(*cli.FlagBotToken)
		if err != nil {
			return nil, fmt.Errorf("cannot create bot API: %w", err)
		}
	}
	return bot, nil
}

func SendMessageToChat(mc tgbotapi.Chattable) error {
	bot, err := getBotAPI()
	if err == nil {
		_, err = bot.Send(mc)
	}
	return err
}
