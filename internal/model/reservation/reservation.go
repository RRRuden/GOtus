package reservation

import (
	"fmt"
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

func (r Reservation) String() string {
	return fmt.Sprintf("[Reservation] ID: %d, BookInstanceID: %d, UserID: %d, StartDate: %s, EndDate: %s, StatusID: %d",
		r.GetID(), r.BookInstanceID, r.UserID, r.StartDate.Format("2006-01-02"), r.EndDate.Format("2006-01-02"), r.ReservationStatusID)
}
