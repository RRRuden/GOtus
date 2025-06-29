package message

import (
	"time"
)

type Message struct {
	id           int
	Email        string
	SendDate     time.Time
	EmailSubject string
	Body         string
}

func NewMessage(id int, email, subject string, sendDate time.Time, body string) *Message {
	return &Message{
		id:           id,
		Email:        email,
		EmailSubject: subject,
		SendDate:     sendDate,
		Body:         body,
	}
}

func (m *Message) GetID() int {
	return m.id
}
