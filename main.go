package main

import (
	"log"
	"net/http"
	"os"

	"github.com/braintree/manners"
	"github.com/go-telegram-bot-api/telegram-bot-api"
)

const tokenEnv string = "TOKEN"

type State uint

const (
	InitState    = 0
	GalleryState = 1
)

func main() {

	token := os.Getenv(tokenEnv)
	if token == "" {
		token = "167627134:AAH5K1D2IBeNs7U3hscqT7rHggfMhWPLDUs"
		//log.Panic("TOKEN ENV NOT FOUND!")
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

	// WebServer Workaround
	// These routes are pingged for services just to never idles this instance
	mux := http.NewServeMux()
	// Starting a Web Server never idles our instance
	mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("pong"))
	})
	mux.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("<a href=\"https://www.uptimia.com/\" target=\"_blank\"><img src=\"https://www.uptimia.com/status?hash=54bd14d2474753185fb6a66b9239a3f8\" width=\"130\" height=\"auto\" alt=\"Website monitoring | Uptimia\" title=\"Website monitoring | Uptimia\"></a>"))
	})
	// Shut the server down gracefully
	processStopedBySignal()
	// Manners allows you to shut your Go webserver down gracefully, without dropping any requests
	err = manners.ListenAndServe(ip+":"+port, mux)
	if err != nil {
		log.Panic(err)
		return
	} else {
		log.Println("Server listening at " + ip + ":" + port)
	}
	defer manners.Close()
	// END of WebServer Workaround

	for update := range updates {

		//log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		// This bot doesn't answare for inline queries
		if update.InlineQuery != nil {
			continue
		}

		u := GetUser(update)

		// Block if user not in session
		if u == nil {
			u = NewUser(bot, update)
			u.State = InitState
		}

		// Handle the actual command
		go Handle(u, update)

	}

}
