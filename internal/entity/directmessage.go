package entity

import (
	"strconv"

	"github.com/dghubble/go-twitter/twitter"
)

type DirectMessage twitter.DirectMessage
type DirectMessages []DirectMessage

func NewDirectMessageFromDirectMessageEvent(
	event twitter.DirectMessageEvent,
	recipient *twitter.User,
	sender *twitter.User,
) DirectMessage {
	id, _ := strconv.ParseInt(event.ID, 10, 64)
	recipientID, _ := strconv.ParseInt(event.Message.Target.RecipientID, 10, 64)
	senderID, _ := strconv.ParseInt(event.Message.SenderID, 10, 64)

	return DirectMessage{
		CreatedAt:           event.CreatedAt,
		Entities:            event.Message.Data.Entities,
		ID:                  id,
		IDStr:               event.ID,
		Recipient:           recipient,
		RecipientID:         recipientID,
		RecipientScreenName: recipient.ScreenName,
		Sender:              sender,
		SenderID:            senderID,
		SenderScreenName:    sender.ScreenName,
		Text:                event.Message.Data.Text,
	}
}

func (e DirectMessage) NewMessage(userID string) Message {
	return Message{
		ForUserID: userID,
		Type:      MessageTypeDirectMessage,
		Data:      e,
	}
}

func (e DirectMessages) Messages(userID string) []Message {
	ms := make([]Message, 0, len(e))

	for _, m := range e {
		ms = append(ms, m.NewMessage(userID))
	}

	return ms
}
