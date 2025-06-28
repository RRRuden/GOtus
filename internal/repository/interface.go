package repository

import (
	"gotus/internal/model/book"
	"gotus/internal/model/reservation"
	"gotus/internal/model/user"
)

//go:generate mockgen -source=interface.go -destination=../mocks/book_repository_mock.go -package=mocks
type BookRepository interface {
	StoreBook(b *book.Book)
	GetBooks() ([]*book.Book, int)
	UpdateBookByISBN(isbn string, updatedBook *book.Book) bool
	DeleteBookByISBN(isbn string) bool
	FindBookByISBN(isbn string) (*book.Book, bool)
}

//go:generate mockgen -source=interface.go -destination=../mocks/book_instance_repository_mock.go -package=mocks
type BookInstanceRepository interface {
	StoreBookInstance(bi *book.BookInstance)
	GetBookInstances() ([]*book.BookInstance, int)
	UpdateBookInstanceById(id int, updatedBookInstance *book.BookInstance) bool
	FindBookInstanceById(id int) (*book.BookInstance, bool)
	GetBookInstancesByISBN(isbn string) ([]*book.BookInstance, int)
	DeleteBookInstanceById(id int) bool
}

//go:generate mockgen -source=interface.go -destination=../mocks/reservation_repository_mock.go -package=mocks
type ReservationRepository interface {
	StoreReservation(r *reservation.Reservation)
	GetReservations() ([]*reservation.Reservation, int)
	UpdateReservationById(id int, updatedReservation *reservation.Reservation) bool
	FindReservationById(id int) (*reservation.Reservation, bool)
	DeleteReservationById(id int) bool
	HasActiveReservation(bookInstanceID int) bool
}

//go:generate mockgen -source=interface.go -destination=../mocks/user_repository_mock.go -package=mocks
type UserRepository interface {
	StoreUser(u *user.User)
	GetUsers() ([]*user.User, int)
	UpdateUserById(id int, updatedReservation *user.User) bool
	FindUserById(id int) (*user.User, bool)
	FindUserByEmail(email string) (*user.User, bool)
	DeleteUserById(id int) bool
}
