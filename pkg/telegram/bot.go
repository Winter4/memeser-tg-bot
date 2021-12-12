package telegram

import (
	"encoding/json"
	"log"
	"memeserBot/pkg/storage"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type MyBot struct {
	bot         tgbotapi.BotAPI
	subsStorage *storage.Storage
}

type config struct {
	AdminChatID     int64  `json:"adminChatID"`
	StorageFileName string `json:"storageFileName"`
	BotApiToken     string `json:"botApiToken"`
}

var conf config

// ________________________________________________________________

func getConfig() {

	file, err := os.Open("conf.json")
	if err != nil {
		log.Println("Error on openinng json conf")
		log.Fatal(err)
	}
	defer file.Close()

	jsonDecoder := json.NewDecoder(file)
	jsonDecoder.Decode(&conf)
}

func NewBot() *MyBot {

	getConfig()

	bot, err := tgbotapi.NewBotAPI(conf.BotApiToken)
	if err != nil {
		log.Println("Error on authing bot")
		log.Fatal(err)
	}

	return &MyBot{
		bot:         *bot,
		subsStorage: storage.NewStorage(conf.StorageFileName),
	}
}

func (b *MyBot) Start() {

	// auth TG server
	log.Printf("Authrized on account %s", b.bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := b.bot.GetUpdatesChan(u)
	if err != nil {
		log.Fatal(err)
	}

	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		if update.Message.IsCommand() {
			b.handleCommand(*update.Message)
		} else {
			b.handleMessage(*update.Message)
		}
	}
}
