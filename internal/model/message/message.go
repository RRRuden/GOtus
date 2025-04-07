package message

import (
	"time"
)

type Message struct {
	id           int
	Email        string
	SendDate     time.Time
	EmailSubject string
}

func NewMessage(id int, email, subject string, sendDate time.Time) *Message {
	return &Message{
		id:           id,
		Email:        email,
		EmailSubject: subject,
		SendDate:     sendDate,
	}
}

func (m *Message) GetID() int {
	return m.id
}
