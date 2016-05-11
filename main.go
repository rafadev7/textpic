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

		if update.Message != nil && update.Message.Document != nil {
			u.Println("Please don't send as a file!")
		}

		// This bot doesn't answare for inline queries
		if update.Message != nil {
			if update.Message.Photo != nil {
				// Take the biggest
				pic := tgbotapi.PhotoSize{}
				pic.Width = 0
				for _, p := range *update.Message.Photo {
					if p.Width > pic.Width {
						pic = p
					}
				}

				if pic.FileID == "" {
					u.Println("Error, FileID is empty")
				}

				url, err := bot.GetFileDirectURL(pic.FileID)
				if err != nil {
					u.Println("Error transforming img: " + err.Error())
				}

				if url == "" {
					u.Println("Error, URL is empty")
				}

				//u.Println("PHOTO URL: " + url)

				bytesImg, err := TransformImage(url, 1696, 2560)
				if err != nil {
					u.Println("Error transforming img: " + err.Error())
				}
				//u.Println("Sending img...")

				b := tgbotapi.FileBytes{Name: "textpic.png", Bytes: bytesImg}

				msg := tgbotapi.NewPhotoUpload(update.Message.Chat.ID, b)
				//msg.Caption = "Test"
				_, err = bot.Send(msg)

				if err != nil {
					u.Println("Error transforming img: " + err.Error())
				}
				//u.PrintCode(str)
			}
		}

		// Handle the actual command
		go Handle(u, update)

	}

}
