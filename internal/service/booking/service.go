package booking

import "gotus/internal/model/reservation"

type Service interface {
	CreateBooking(userID int, isbn string) (*reservation.Reservation, error)
	ExtendBooking(reservationID int, extensionDays int) (bool, error)
	CancelBooking(reservationID int) (bool, error)
	EndBooking(reservationID int) (bool, error)
}
