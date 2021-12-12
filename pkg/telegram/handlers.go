package telegram

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func (b *MyBot) handleMessage(msg tgbotapi.Message) {

	// if the message was sent by admin
	if msg.Chat.ID == conf.AdminChatID {
		b.sendToSubs(msg)
	}
}

func (b *MyBot) sendToSubs(msg tgbotapi.Message) {

	// open file
	file, err := os.Open(conf.StorageFileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {

		pair := strings.Split(scanner.Text(), "=")

		if pair[1] == "1" {
			chatID, err := strconv.Atoi(pair[0])
			if err != nil {
				log.Fatal(err)
			}
			fwd := tgbotapi.NewForward(int64(chatID), conf.AdminChatID, msg.MessageID)
			b.bot.Send(fwd)
		}
	}
}

func (b *MyBot) handleCommand(msg tgbotapi.Message) {

	switch msg.Command() {

	case "start":
		b.answer(b.subsStorage.Start, msg)

	case "sub":
		b.answer(b.subsStorage.Subscribe, msg)

	case "unsub":
		b.answer(b.subsStorage.Unsubscribe, msg)

	case "stop":
		os.Exit(0)
	}
}

func (b *MyBot) answer(callback func(int64) string, msg tgbotapi.Message) {
	text := callback(msg.Chat.ID)
	answer := tgbotapi.NewMessage(msg.Chat.ID, text)

	b.bot.Send(answer)
}
