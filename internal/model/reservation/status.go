package reservation

type ReservationStatus struct {
	id   int
	Name string
}

func NewReservationStatus(id int, name string) *ReservationStatus {
	return &ReservationStatus{
		id:   id,
		Name: name,
	}
}

func (rs *ReservationStatus) GetID() int {
	return rs.id
}
