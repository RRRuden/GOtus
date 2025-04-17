// internal/repository/storage.go
package repository

import (
	"fmt"
	"sync"

	"gotus/internal/model/book"
	"gotus/internal/model/reservation"
	"gotus/internal/model/user"
)

var (
	booksMutex sync.Mutex
	books      []*book.Book

	bookInstancesMutex sync.Mutex
	bookInstances      []*book.BookInstance

	usersMutex sync.Mutex
	users      []*user.User

	reservationsMutex sync.Mutex
	reservations      []*reservation.Reservation
)

func Store(data fmt.Stringer) {
	switch v := data.(type) {
	case *book.Book:
		booksMutex.Lock()
		books = append(books, v)
		booksMutex.Unlock()
	case *book.BookInstance:
		bookInstancesMutex.Lock()
		bookInstances = append(bookInstances, v)
		bookInstancesMutex.Unlock()
	case *user.User:
		usersMutex.Lock()
		users = append(users, v)
		usersMutex.Unlock()
	case *reservation.Reservation:
		reservationsMutex.Lock()
		reservations = append(reservations, v)
		reservationsMutex.Unlock()
	default:
		fmt.Println("Неизвестный тип данных")
	}
}

func GetBooks() ([]*book.Book, int) {
	booksMutex.Lock()
	defer booksMutex.Unlock()
	return books, len(books)
}

func GetBookInstances() ([]*book.BookInstance, int) {
	bookInstancesMutex.Lock()
	defer bookInstancesMutex.Unlock()
	return bookInstances, len(bookInstances)
}

func GetUsers() ([]*user.User, int) {
	usersMutex.Lock()
	defer usersMutex.Unlock()
	return users, len(users)
}

func GetReservations() ([]*reservation.Reservation, int) {
	reservationsMutex.Lock()
	defer reservationsMutex.Unlock()
	return reservations, len(reservations)
}
