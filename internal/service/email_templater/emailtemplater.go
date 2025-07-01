package emailtemplater

import (
	"bytes"
	"errors"
	"html/template"
	"path/filepath"
	"time"

	"gotus/internal/model/reservation"
	"gotus/internal/repository"
)

var (
	ErrReservationNotFound  = errors.New("заказ не найлен")
	ErrBookInstanceNotFound = errors.New("экземпляр книги не найден")
	ErrBookNotFound         = errors.New("книга не найдена")
)

type emailTemplater struct {
	reservationRepo  repository.ReservationRepository
	bookInstanceRepo repository.BookInstanceRepository
	bookRepo         repository.BookRepository
	templateDir      string
}

func NewEmailTemplater(
	reservationRepo repository.ReservationRepository,
	bookInstanceRepo repository.BookInstanceRepository,
	bookRepo repository.BookRepository,
	templateDir string,
) EmailTemplater {
	return &emailTemplater{
		reservationRepo:  reservationRepo,
		bookInstanceRepo: bookInstanceRepo,
		bookRepo:         bookRepo,
		templateDir:      templateDir,
	}
}

type emailData struct {
	BookTitle string
	StartDate string
	EndDate   string
	Now       string
}

func (e *emailTemplater) getBookTitleByReservation(reservation *reservation.Reservation) (string, error) {
	bookInstance, found, err := e.bookInstanceRepo.FindBookInstanceById(reservation.BookInstanceID)
	if err != nil {
		return "", err
	}
	if !found {
		return "", ErrBookInstanceNotFound
	}

	book, found, err := e.bookRepo.FindBookByISBN(bookInstance.ISBN)
	if err != nil {
		return "", err
	}
	if !found {
		return "", ErrBookNotFound
	}

	return book.Title, nil
}

func (e *emailTemplater) renderTemplate(templateName string, data emailData) (string, error) {
	path := filepath.Join(e.templateDir, templateName)
	tmpl, err := template.ParseFiles(path)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func (e *emailTemplater) GetBookingCreatedEmail(reservationId int) (string, error) {
	reservation, found, err := e.reservationRepo.FindReservationById(reservationId)
	if err != nil {
		return "", err
	}
	if !found {
		return "", ErrReservationNotFound
	}

	bookTitle, err := e.getBookTitleByReservation(reservation)
	if err != nil {
		return "", err
	}

	data := emailData{
		BookTitle: bookTitle,
		StartDate: reservation.StartDate.Format("02 Jan 2006"),
		EndDate:   reservation.EndDate.Format("02 Jan 2006"),
	}

	return e.renderTemplate("booking_created.html", data)
}

func (e *emailTemplater) GetBookingCanceledEmail(reservationId int) (string, error) {
	reservation, found, err := e.reservationRepo.FindReservationById(reservationId)
	if err != nil {
		return "", err
	}
	if !found {
		return "", ErrReservationNotFound
	}

	bookTitle, err := e.getBookTitleByReservation(reservation)
	if err != nil {
		return "", err
	}

	data := emailData{
		BookTitle: bookTitle,
		Now:       time.Now().Format("02 Jan 2006 15:04"),
	}

	return e.renderTemplate("booking_canceled.html", data)
}

func (e *emailTemplater) GetBookingEndedEmail(reservationId int) (string, error) {
	reservation, found, err := e.reservationRepo.FindReservationById(reservationId)
	if err != nil {
		return "", err
	}
	if !found {
		return "", ErrReservationNotFound
	}

	bookTitle, err := e.getBookTitleByReservation(reservation)
	if err != nil {
		return "", err
	}

	data := emailData{
		BookTitle: bookTitle,
		EndDate:   reservation.EndDate.Format("02 Jan 2006"),
	}

	return e.renderTemplate("booking_ended.html", data)
}

func (e *emailTemplater) GetBookingExtendedEmail(reservationId int) (string, error) {
	reservation, found, err := e.reservationRepo.FindReservationById(reservationId)
	if err != nil {
		return "", err
	}
	if !found {
		return "", ErrReservationNotFound
	}

	bookTitle, err := e.getBookTitleByReservation(reservation)
	if err != nil {
		return "", err
	}

	data := emailData{
		BookTitle: bookTitle,
		EndDate:   reservation.EndDate.Format("02 Jan 2006"),
	}

	return e.renderTemplate("booking_extended.html", data)
}
