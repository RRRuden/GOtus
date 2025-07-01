package emailsender

import (
	"fmt"
	"net/smtp"
	"time"

	"gotus/internal/config"
	"gotus/internal/model/message"
	"gotus/internal/repository"
)

type SMTPEmailSender struct {
	host string
	port int
	from string
	repo repository.EmailMessageRepository
	auth smtp.Auth
}

func NewSMTPEmailSender(cfg config.SMTPConfig, repo repository.EmailMessageRepository) EmailSender {
	return &SMTPEmailSender{
		host: cfg.Host,
		port: cfg.Port,
		from: cfg.From,
		repo: repo,
	}
}

func (s *SMTPEmailSender) SendEmail(to, subject, body string) error {
	addr := fmt.Sprintf("%s:%d", s.host, s.port)

	msg := fmt.Sprintf(
		"From: %s\r\nTo: %s\r\nSubject: %s\r\n"+
			"MIME-version: 1.0;\r\nContent-Type: text/html; charset=\"UTF-8\";\r\n\r\n%s",
		s.from, to, subject, body,
	)

	// Отправка письма
	err := smtp.SendMail(addr, s.auth, s.from, []string{to}, []byte(msg))
	if err != nil {
		return err
	}

	// Логируем сообщение в БД
	msgRecord := message.NewMessage(0, to, subject, time.Now(), body)

	err = s.repo.InsertMessage(msgRecord)

	if err != nil {
		return err
	}

	return nil
}
