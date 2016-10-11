package main

import (
	"fmt"
	logger "log"
	"net/http"
	"os"

	"github.com/nats-io/nats"
	"gopkg.in/maciekmm/messenger-platform-go-sdk.v4"
)

var mainFBEventHandler *FBBotEventHandler
var mainFBBotClient *messenger.Messenger
var log *logger.Logger
var nc *nats.Conn

func main() {
	f, err := os.OpenFile("./fbbot.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0660)
	if err != nil {
		log.Fatalf("Can't open log file: %v\n", err)
	}
	log = new(logger.Logger)
	log.SetOutput(f)
	log.SetFlags(logger.LstdFlags)

	urls := "nats://localhost:4222"
	nc, err = nats.Connect(urls)
	if err != nil {
		log.Fatalf("Can't connect to NATS: %v\n", err)
	}
	defer nc.Close()
	mainFBBotClient := &messenger.Messenger{
		VerifyToken: os.Getenv("FB_VERIFY_TOKEN"),
		AppSecret:   os.Getenv("FB_APP_SECRET"),
		AccessToken: os.Getenv("FB_PAGE_TOKEN"),
		//    Debug: messenger.DebugAll, //All,Info,Warning
	}

	mainFBEventHandler = new(FBBotEventHandler)
	mainFBEventHandler.SetFBBotClient(mainFBBotClient)

	mainFBBotClient.MessageReceived = mainFBEventHandler.OnEventMessage
	mainFBBotClient.MessageDelivered = mainFBEventHandler.OnEventDelivery
	mainFBBotClient.Postback = mainFBEventHandler.OnEventPostback
	http.HandleFunc("/fb_callback", mainFBBotClient.Handler)

	port := os.Getenv("FB_BOT_PORT")
	addr := fmt.Sprintf(":%s", port)
	http.ListenAndServe(addr, nil)
}
