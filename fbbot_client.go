package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/orsonwang/fbot"
)

var mainFBEventHandler *FBBotEventHandler
var mainFBBotClient *fbot.Bot

func main() {

	mainFBBotClient = fbot.NewBot(fbot.Config{
		AccessToken: os.Getenv("FB_PAGE_TOKEN"),
		AppSecret:   os.Getenv("FB_APP_SECRET"),
		VerifyToken: os.Getenv("FB_VERIFY_TOKEN"),
	})
	mainFBEventHandler = new(FBBotEventHandler)
	mainFBEventHandler.SetFBBotClient(mainFBBotClient)

	mainFBBotClient.On(fbot.EventMessage, mainFBEventHandler.OnEventMessage)
	mainFBBotClient.On(fbot.EventDelivery, mainFBEventHandler.OnEventDelivery)
	mainFBBotClient.On(fbot.EventPostback, mainFBEventHandler.OnEventPostback)

	http.Handle("/fb_callback", fbot.Handler(mainFBBotClient))

	port := os.Getenv("FB_BOT_PORT")
	addr := fmt.Sprintf(":%s", port)
	http.ListenAndServe(addr, nil)
}
