package reservation

const (
	StatusBooked    int = iota + 1 // Забронирована
	StatusExtended                 // Продлена
	StatusCancelled                // Отменена
	StatusEnded                    // Завершена
)

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
