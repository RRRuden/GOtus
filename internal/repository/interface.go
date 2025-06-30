package repository

import (
	"gotus/internal/model/book"
	"gotus/internal/model/message"
	"gotus/internal/model/reservation"
	"gotus/internal/model/user"
)

//go:generate mockgen -source=interface.go -destination=../mocks/repository_mock.go -package=mocks
type BookRepository interface {
	StoreBook(b *book.Book) error
	GetBooks() ([]*book.Book, int, error)
	UpdateBookByISBN(isbn string, updatedBook *book.Book) (bool, error)
	DeleteBookByISBN(isbn string) (bool, error)
	FindBookByISBN(isbn string) (*book.Book, bool, error)
}

type BookInstanceRepository interface {
	StoreBookInstance(bi *book.BookInstance) error
	GetBookInstances() ([]*book.BookInstance, int, error)
	UpdateBookInstanceById(id int, updatedBookInstance *book.BookInstance) (bool, error)
	FindBookInstanceById(id int) (*book.BookInstance, bool, error)
	GetBookInstancesByISBN(isbn string) ([]*book.BookInstance, int, error)
	DeleteBookInstanceById(id int) (bool, error)
}

type ReservationRepository interface {
	StoreReservation(r *reservation.Reservation) (int, error)
	GetReservations() ([]*reservation.Reservation, int, error)
	UpdateReservationById(id int, updatedReservation *reservation.Reservation) (bool, error)
	FindReservationById(id int) (*reservation.Reservation, bool, error)
	DeleteReservationById(id int) (bool, error)
	HasActiveReservation(bookInstanceID int) (bool, error)
}

type UserRepository interface {
	StoreUser(u *user.User) error
	GetUsers() ([]*user.User, int, error)
	UpdateUserById(id int, updatedUser *user.User) (bool, error)
	FindUserById(id int) (*user.User, bool, error)
	FindUserByEmail(email string) (*user.User, bool, error)
	DeleteUserById(id int) (bool, error)
}

type EmailMessageRepository interface {
	InsertMessage(m *message.Message) error
}
