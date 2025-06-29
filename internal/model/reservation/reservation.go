package reservation

import (
	"fmt"
	"time"
)

type Reservation struct {
	Id                  int       `json:"Id" bson:"id"`
	BookInstanceID      int       `json:"BookInstanceID" bson:"bookInstance_id"`
	UserID              int       `json:"UserID" bson:"user_id"`
	StartDate           time.Time `json:"StartDate" bson:"start_date"`
	EndDate             time.Time `json:"EndDate" bson:"end_date"`
	ReservationStatusID int       `json:"ReservationStatusID" bson:"reservation_status_id"`
}

func NewReservation(id, bookInstanceID, userID, statusID int, startDate, endDate time.Time) *Reservation {
	return &Reservation{
		Id:                  id,
		BookInstanceID:      bookInstanceID,
		UserID:              userID,
		StartDate:           startDate,
		EndDate:             endDate,
		ReservationStatusID: statusID,
	}
}

func (r *Reservation) GetID() int {
	return r.Id
}

func (r *Reservation) SetID(id int) {
	r.Id = id
}

func (r Reservation) String() string {
	return fmt.Sprintf("[Reservation] ID: %d, BookInstanceID: %d, UserID: %d, StartDate: %s, EndDate: %s, StatusID: %d",
		r.GetID(), r.BookInstanceID, r.UserID, r.StartDate.Format("2006-01-02"), r.EndDate.Format("2006-01-02"), r.ReservationStatusID)
}
