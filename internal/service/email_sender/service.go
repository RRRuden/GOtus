package emailsender

//go:generate mockgen -source=service.go -destination=../../mocks/emailsender_mock.go -package=mocks
type EmailSender interface {
	SendEmail(to, subject, body string) error
}
