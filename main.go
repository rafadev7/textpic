package main

import (
	"log"
	"os"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)

const (
	tokenEnv = "TOKEN"
	ipEnv    = "OPENSHIFT_GO_IP"
	portEnv  = "OPENSHIFT_GO_PORT"
)

type State uint

const (
	InitState    = 0
	GalleryState = 1
)

func main() {

	token := os.Getenv(tokenEnv)
	if token == "" {
		log.Panicln("TOKEN ENV NOT FOUND!")
	}

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panicln(err)
	}

	// Start this webserver just to never puts this instance idle
	go StartWebServer()

	bot.Debug = true

	log.Printf("Authorized on account %s\n", bot.Self.UserName)

	news := tgbotapi.NewUpdate(0)
	news.Timeout = 60

	updates, err := bot.GetUpdatesChan(news)

	for update := range updates {

		//log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		u := GetUser(update)
		if u == nil {
			u = NewUser(bot, update)
			u.State = InitState
		}

		// Handle the actual command
		go Handle(u, update)

	}

}
