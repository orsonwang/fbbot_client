package main

import (
	"regexp"
	//	"strings"
	"time"

	"gopkg.in/maciekmm/messenger-platform-go-sdk.v4"
)

// FBBotEventHandler ...
type FBBotEventHandler struct {
	botClient *messenger.Messenger
}

// SetFBBotClient to assign fbbot client handler
func (s *FBBotEventHandler) SetFBBotClient(bc *messenger.Messenger) {
	s.botClient = bc
}

// OnEventMessage ...
func (s *FBBotEventHandler) OnEventMessage(event messenger.Event, opts messenger.MessageOpts, msg messenger.ReceivedMessage) {
	log.Println("=== FB Callback ===")
	log.Printf("Sender ID: %s\n", opts.Sender.ID)
	log.Printf("Recipent ID: %s\n", opts.Recipient.ID)
	log.Printf("Timestame: %d\n", opts.Timestamp)
	log.Printf("Message ID: %s\n", msg.ID)
	log.Printf("Message Seq: %d\n", msg.Seq)
	log.Printf("Message Text: %s\n", msg.Text)

	if len(msg.Attachments) != 0 {
		for _, attachment := range msg.Attachments {
			log.Println(attachment.Type)
			//log.Println(attachment.payload.URL)
		}
	}
	s.processTextMessage(opts.Sender.ID, msg.Text)
}

// OnEventDelivery ...
func (s *FBBotEventHandler) OnEventDelivery(event messenger.Event, opts messenger.MessageOpts, delivery messenger.Delivery) {
	log.Println(delivery.MessageIDS[0])
	log.Println(delivery.Watermark)
	log.Println(delivery.Seq)
}

// OnEventPostback ...
func (s *FBBotEventHandler) OnEventPostback(event messenger.Event, opts messenger.MessageOpts, postback messenger.Postback) {
	log.Println(postback.Payload)
}

func (s *FBBotEventHandler) matchString(pattern, text string) (result bool) {
	result, _ = regexp.MatchString(pattern, text)
	return
}

// processTextMessage ...
func (s *FBBotEventHandler) processTextMessage(from string, text string) {
	log.Printf("Received text \"%s\" from %s", text, from)

	subj := "aitc.text.service"
	msg, err := nc.Request(subj, []byte(text), 3*time.Second)
	if err != nil {
		log.Fatalf("Error in Request: %v\n", err)
	}
	strResult := string(msg.Data)
	log.Printf("Return text \"%s\" from text_service", strResult)
	if len(strResult) != 0 {
		_, err = s.botClient.SendSimpleMessage(from, strResult)
		if err != nil {
			log.Fatalf("Error reply text message: %v\n", err)
		}
	} else {
		s.botClient.SendSimpleMessage(from, "我無法了解您的指令")
	}
}
