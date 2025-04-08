// internal/repository/storage.go
package repository

import (
	"fmt"
	"gotus/internal/model/book"
	"gotus/internal/model/reservation"
	"gotus/internal/model/user"
)

var (
	books         []*book.Book
	bookInstances []*book.BookInstance
	users         []*user.User
	reservations  []*reservation.Reservation
)

func Store(data fmt.Stringer) {
	switch v := data.(type) {
	case *book.Book:
		books = append(books, v)
		fmt.Println(v.String(), "добавлен в хранилище")
	case *book.BookInstance:
		bookInstances = append(bookInstances, v)
		fmt.Println(v.String(), "добавлен в хранилище")
	case *user.User:
		users = append(users, v)
		fmt.Println(v.String(), "добавлен в хранилище")
	case *reservation.Reservation:
		reservations = append(reservations, v)
		fmt.Println(v.String(), "добавлен в хранилище")
	default:
		fmt.Printf("Неизвестный тип данных")
		return
	}
}
