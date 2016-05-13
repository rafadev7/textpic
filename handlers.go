package main

import "github.com/go-telegram-bot-api/telegram-bot-api"

func Handle(u *User, update tgbotapi.Update) {

	// In case of user sending images
	ImageProcess(u, update)

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

func ImageProcess(u *User, update tgbotapi.Update) {

	if update.Message != nil && update.Message.Document != nil {
		u.Println("Please don't send as a file!")
		return
	}

	// This bot doesn't answare for inline queries
	if update.Message == nil || update.Message.Photo == nil {
		return
	}

	// Take the biggest photo sent
	pic := tgbotapi.PhotoSize{}
	pic.Width = 0
	for _, p := range *update.Message.Photo {
		if p.Width > pic.Width {
			pic = p
		}
	}

	// It should never occours
	if pic.FileID == "" {
		u.Println("Error, FileID is empty")
		return
	}

	url, err := u.Bot.GetFileDirectURL(pic.FileID)
	if err != nil {
		u.Println("Error getting the url: " + err.Error())
		return
	}

	// It should never occours
	if url == "" {
		u.Println("Error, URL is empty")
		return
	}

	text, size, err := ImageToText(url)
	if err != nil {
		u.Println("Error transforming img to text: " + err.Error())
		return
	}

	bytesImg, err := TextToImage(text, size)
	if err != nil {
		u.Println("Error transforming text to img: " + err.Error())
		return
	}

	b := tgbotapi.FileBytes{Name: "textpic.png", Bytes: bytesImg}

	msg := tgbotapi.NewPhotoUpload(update.Message.Chat.ID, b)
	_, err = u.Bot.Send(msg)

	if err != nil {
		u.Println("Error sending you the img: " + err.Error())
	}

}
