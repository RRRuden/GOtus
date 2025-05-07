package repository

import (
	"fmt"
	"gotus/internal/model/book"
	"gotus/internal/model/reservation"
	"gotus/internal/model/user"
)

func Store(data fmt.Stringer) {
	switch v := data.(type) {
	case *book.Book:
		StoreBook(v)
	case *book.BookInstance:
		StoreBookInstance(v)
	case *user.User:
		StoreUser(v)
	case *reservation.Reservation:
		StoreReservation(v)
	default:
		fmt.Println("Неизвестный тип данных")
	}
}

func LoadAllFromCSV() {
	LoadBooksFromCSV()
	LoadBookInstancesFromCSV()
	LoadUsersFromCSV()
	LoadReservationsFromCSV()
}
