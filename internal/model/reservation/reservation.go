package reservation

import (
	"time"
)

type Reservation struct {
	id                  int
	BookInstanceID      int
	UserID              int
	StartDate           time.Time
	EndDate             time.Time
	ReservationStatusID int
}

func NewReservation(id, bookInstanceID, userID, statusID int, startDate, endDate time.Time) *Reservation {
	return &Reservation{
		id:                  id,
		BookInstanceID:      bookInstanceID,
		UserID:              userID,
		StartDate:           startDate,
		EndDate:             endDate,
		ReservationStatusID: statusID,
	}
}

func (r *Reservation) GetID() int {
	return r.id
}
