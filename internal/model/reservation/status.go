package reservation

const (
	StatusBooked    int = iota + 1 // Забронирована
	StatusExtended                 // Продлена
	StatusCancelled                // Отменена
	StatusEnded                    // Завершена
)
