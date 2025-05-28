package repository

import (
	"fmt"
	"gotus/internal/model/book"
	"gotus/internal/model/reservation"
	"gotus/internal/model/user"
)

type Storage struct {
	BookRepository         *BookRepository
	BookInstanceRepository *BookInstanceRepository
	UserRepository         *UserRepository
	ReservationRepository  *ReservationRepository
}

func NewStorage(bookRepo *BookRepository, bookInstanceRepo *BookInstanceRepository, userRepo *UserRepository, resRepo *ReservationRepository) *Storage {
	return &Storage{
		BookRepository:         bookRepo,
		BookInstanceRepository: bookInstanceRepo,
		UserRepository:         userRepo,
		ReservationRepository:  resRepo,
	}
}

func (s *Storage) Store(data fmt.Stringer) {
	switch v := data.(type) {
	case *book.Book:
		s.BookRepository.StoreBook(v)
	case *book.BookInstance:
		s.BookInstanceRepository.StoreBookInstance(v)
	case *user.User:
		s.UserRepository.StoreUser(v)
	case *reservation.Reservation:
		s.ReservationRepository.StoreReservation(v)
	default:
		fmt.Println("Неизвестный тип данных")
	}
}

func (s *Storage) LoadAllFromCSV() {
	s.BookRepository.LoadBooksFromCSV()
	s.BookInstanceRepository.LoadBookInstancesFromCSV()
	s.UserRepository.LoadUsersFromCSV()
	s.ReservationRepository.LoadReservationsFromCSV()
}
