package main

import (
	"regexp"
	"strings"
	"time"

	"github.com/orsonwang/fbot"
)

// FBBotEventHandler ...
type FBBotEventHandler struct {
	botClient *fbot.Bot
}

// SetFBBotClient to assign fbbot client handler
func (s *FBBotEventHandler) SetFBBotClient(bc *fbot.Bot) {
	s.botClient = bc
}

// OnEventMessage ...
func (s *FBBotEventHandler) OnEventMessage(event *fbot.Event) {
	log.Println("=== FB Callback ===")
	log.Println(event.Sender.ID)
	log.Println(event.Recipient.ID)
	log.Println(event.Timestamp)
	log.Println(event.Message.Mid)
	log.Println(event.Message.Seq)
	log.Println(event.Message.Text)

	if len(event.Message.Attachments) != 0 {
		for _, attachment := range event.Message.Attachments {
			log.Println(attachment.Type)
			log.Println(attachment.Payload.URL)
		}
	}
	s.processTextMessage(event.Sender, event.Message.Text)
}

// OnEventDelivery ...
func (s *FBBotEventHandler) OnEventDelivery(event *fbot.Event) {
	log.Println(event.Delivery.Mids[0])
	log.Println(event.Delivery.Watermark)
	log.Println(event.Delivery.Seq)
}

// OnEventPostback ...
func (s *FBBotEventHandler) OnEventPostback(event *fbot.Event) {
	log.Println(event.Postback.Payload)
}

func (s *FBBotEventHandler) matchString(pattern, text string) (result bool) {
	result, _ = regexp.MatchString(pattern, text)
	return
}

// processTextMessage ...
func (s *FBBotEventHandler) processTextMessage(from *fbot.User, text string) {
	strAfterCut := strings.ToUpper(text)
	log.Printf("Received text \"%s\" from %s", text, from)

	subj := "aitc.text.service"
	msg, err := nc.Request(subj, []byte(strAfterCut), 1*time.Second)
	if err != nil {
		log.Fatalf("Error in Request: %v\n", err)
	}
	strResult := string(msg.Data)
	if len(strResult) != 0 {
		s.botClient.Deliver(fbot.DeliverParams{
			Recipient: from,
			Message: &fbot.Message{
				Text: strResult,
			},
		})
	} else {
		s.botClient.Deliver(fbot.DeliverParams{
			Recipient: from,
			Message: &fbot.Message{
				Attachment: &fbot.Attachment{
					Type: "image",
					Payload: &fbot.Payload{
						URL: "https://linebot.gaze.tw/exrate.png",
					},
				},
			},
		})
	}

}
