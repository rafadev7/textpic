package main

import "github.com/go-telegram-bot-api/telegram-bot-api"

func Handle(u *User, update tgbotapi.Update) {
	switch u.State {
	case InitState:
		InitHandler(u, update)
	case GalleryState:
		GalleryHandler(u, update)
	}

	// Insert instructions for next command
	PostProcess(u)
}

func PostProcess(u *User) {
	switch u.State {
	case InitState:
		InitHeader(u)
	case GalleryState:
		GalleryHeader(u)
	}

}

func InitHeader(u *User) {
}

func InitHandler(u *User, update tgbotapi.Update) {

	switch update.Message.Text {
	case "/start":
		u.Println("Transforms any image into ASCII Art.")
		u.Println("Send us a pic  and we will create a ASCII Art based in your image!")

	case "/help":
		u.Println("Just send us a photo media and we will return the ASCII Art from it.")
		u.Println("Type /about to know more about this project")
		u.Println("Type /rate to give us five stars!")

	case "/about":
		u.Println("This is the first bot to transform your images in ASCII Art")
		u.Println("It's an open-source project found in github.com/rafadev7/textpic")
		u.Println("We don't store any information you send")
		u.Println("If you got interested then access our github pages and get involved with the project")
		u.Println("Chose one of the options in the keyboard bellow")

	case "/rate":
		SendRateInline(u)
		u.Println("Then you can choose any options below")

	case "/gallery":
		u.Println("You will see our gallery")
		u.Println("Type /back if you wanna go back")
		u.State = GalleryState

	default:
		SendKeyboard(u, "Welcome to the @TextPicBot!")
	}

}

func GalleryHeader(u *User) {
	SendKeyboard(u, "If you want to see the next Art, just type /next")
}

func GalleryHandler(u *User, update tgbotapi.Update) {

	switch update.Message.Text {

	case "/Previous", "Previous":
		u.Println("Look this art!")

	case "/next", "Next":
		u.Println("Look this art!")

	case string('\U0001F44D'):
		u.Println("You liked it!")

	case string('\U00002764'):
		u.Println("You loved it!")

	case string('\U0001F44E'):
		u.Println("You didn't like it!")

	case "/back":
		u.HideKeyboard("Backing to Init")
		u.State = InitState

	}
}
