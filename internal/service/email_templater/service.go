package emailtemplater

//go:generate mockgen -source=service.go -destination=../../mocks/emailtemplater_mock.go -package=mocks
type EmailTemplater interface {
	GetBookingCreatedEmail(reservationId int) (string, error)
	GetBookingCanceledEmail(reservationId int) (string, error)
	GetBookingEndedEmail(reservationId int) (string, error)
	GetBookingExtendedEmail(reservationId int) (string, error)
}
