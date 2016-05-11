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
		log.Panic("TOKEN ENV NOT FOUND!")
	}

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
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

		// Block if user not in session
		if u == nil {
			u = NewUser(bot, update)
			u.State = InitState
		}

		// This bot doesn't answare for inline queries
		if update.Message != nil {
			if update.Message.Photo != nil {
				// Take the smallest photo
				small := tgbotapi.PhotoSize{}
				small.Width = 513
				for _, pic := range *update.Message.Photo {
					if pic.Width < small.Width {
						small = pic
					}
				}

				url, err := bot.GetFileDirectURL(small.FileID)
				if err != nil {
					u.Println("Error transforming img: " + err.Error())
				}

				u.Println("PHOTO URL: " + url)

				str, err := TransformImage(url)
				if err != nil {
					u.Println("Error transforming img: " + err.Error())
				}
				u.Println("IMG:")
				u.PrintCode(str)
			}
		}

		// Handle the actual command
		go Handle(u, update)

	}

}
