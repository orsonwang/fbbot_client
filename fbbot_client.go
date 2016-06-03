package main

import (
	"fmt"
	logger "log"
	"net/http"
	"os"

	"github.com/nats-io/nats"
	"github.com/orsonwang/fbot"
)

var mainFBEventHandler *FBBotEventHandler
var mainFBBotClient *fbot.Bot
var log *logger.Logger
var nc *nats.Conn

func main() {
	f, err := os.OpenFile("./fbbot.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0660)
	if err != nil {
		log.Fatalf("Can't connect: %v\n", err)
	}
	log = new(logger.Logger)
	log.SetOutput(f)

	urls := "nats://localhost:4222"
	nc, err = nats.Connect(urls)
	if err != nil {
		log.Fatalf("Can't connect: %v\n", err)
	}
	defer nc.Close()

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
