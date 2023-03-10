package telegram

import (
	"fmt"
	"log"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/nicksedov/sbconn-bot/pkg/cli"
	"github.com/nicksedov/sbconn-bot/pkg/openai"
)

var bot *tgbotapi.BotAPI

func InitBot() error {
	var err error
	bot, err = tgbotapi.NewBotAPI(*cli.FlagBotToken)
	if err != nil {
		return fmt.Errorf("cannot create bot API: %w", err)
	}
	upd := tgbotapi.NewUpdate(0)
	upd.Timeout = 60
	go updatesListener(bot.GetUpdatesChan(upd))
	return nil
}

func SendMessageToChat(mc tgbotapi.Chattable) error {
	bot, err := getOrInitBot()
	if err == nil {
		_, err = bot.Send(mc)
	}
	return err
}

func getOrInitBot() (*tgbotapi.BotAPI, error) {
	if bot == nil {
		err := InitBot()
		if err != nil {
			return nil, fmt.Errorf("cannot create bot API: %w", err)
		}
	}
	return bot, nil
}

func updatesListener(updates tgbotapi.UpdatesChannel) {
	// `for {` means the loop is infinite until we manually stop it
	for {
		// execution thread locks until event received
		update := <-updates
		if update.Message != nil {
			handleMessage(update.Message)
		}
	}
}

func handleMessage(message *tgbotapi.Message) {
	user := message.From
	text := message.Text
	chatId := message.Chat.ID

	if user == nil {
		return
	}

	// Print to console
	log.Printf("%s wrote %s", user.FirstName, text)

	var err error
	if strings.HasPrefix(text, "/") {
		err = handleCommand(chatId, text)
	} else {
		resp := openai.SendRequest(user.ID, text)
		if len(resp.Choices) > 0 {
			msg := tgbotapi.NewMessage(chatId, resp.Choices[0].Message.Content)
			bot.Send(msg)
		}
	}
	if err != nil {
		log.Printf("An error occured: %s", err.Error())
	}
}

// When we get a command, we react accordingly
func handleCommand(chatId int64, command string) error {
	var err error

	switch command {
	case "/help":
		msg := tgbotapi.NewMessage(chatId, "GPT-3 chatbot")
		_, err = bot.Send(msg)
	}
	return err
}