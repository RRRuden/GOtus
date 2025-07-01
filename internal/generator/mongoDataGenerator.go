package generator

import (
	"fmt"
	"gotus/internal/model/book"
	"gotus/internal/model/reservation"
	"gotus/internal/model/user"
	"gotus/internal/repository"
	"time"
)

func GenerateInitialDataForMongoDb(userRepo repository.UserRepository,
	bookRepo repository.BookRepository,
	bookInstanceRepo repository.BookInstanceRepository,
	reservationRepo repository.ReservationRepository) {
	users := []*user.User{
		user.NewUser(1, "Иван Иванов", "ivan@example.com"),
		user.NewUser(2, "Петр Петров", "petrov@example.com"),
		user.NewUser(3, "Анна Сергеева", "anna@example.com"),
		user.NewUser(4, "Мария Кузнецова", "maria@example.com"),
		user.NewUser(5, "Дмитрий Орлов", "dmitry@example.com"),
		user.NewUser(6, "Лионель Месси", "leomessi@mail.ru"),
	}

	reservations := []*reservation.Reservation{
		reservation.NewReservation(1, 1, 1, reservation.StatusBooked, time.Now(), time.Now().Add(time.Hour*24*7)),
		reservation.NewReservation(2, 2, 2, reservation.StatusBooked, time.Now(), time.Now().Add(time.Hour*24*7)),
		reservation.NewReservation(3, 3, 3, reservation.StatusBooked, time.Now(), time.Now().Add(time.Hour*24*7)),
		reservation.NewReservation(4, 4, 4, reservation.StatusEnded, time.Now().Add(time.Hour*24*-7), time.Now()),
		reservation.NewReservation(5, 5, 5, reservation.StatusCancelled, time.Now().Add(time.Hour*24*-7), time.Now().Add(time.Hour*24*-7)),
		reservation.NewReservation(6, 6, 6, reservation.StatusExtended, time.Now().Add(time.Hour*24*-7), time.Now().Add(time.Hour*24*14)),
	}

	books, _ := generateBooks()

	bookInstances := []*book.BookInstance{
		book.NewBookInstance(1, "978-3-16-148410-1"),
		book.NewBookInstance(2, "978-3-16-148410-2"),
		book.NewBookInstance(3, "978-3-16-148410-3"),
		book.NewBookInstance(4, "978-3-16-148410-4"),
		book.NewBookInstance(5, "978-3-16-148410-5"),
		book.NewBookInstance(6, "978-3-16-148410-6"),
	}

	for _, user := range users {
		userRepo.StoreUser(user)
	}

	for _, reservation := range reservations {
		reservationRepo.StoreReservation(reservation)
	}

	for _, book := range books {
		bookRepo.StoreBook(book)
	}

	for _, bi := range bookInstances {
		bookInstanceRepo.StoreBookInstance(bi)
	}
}

func generateBooks() ([]*book.Book, error) {
	var books []*book.Book

	entries := []struct {
		isbn   string
		title  string
		author string
		year   int
	}{
		{"978-3-16-148410-1", "Преступление и наказание", "Ф. М. Достоевский", 1866},
		{"978-3-16-148410-2", "Война и мир", "Л. Н. Толстой", 1869},
		{"978-3-16-148410-3", "1984", "Джордж Оруэлл", 1949},
		{"978-3-16-148410-4", "Унесённые ветром", "Маргарет Митчелл", 1936},
		{"978-3-16-148410-5", "Мастер и Маргарита", "Михаил Булгаков", 1966},
		{"978-3-16-148410-6", "Властелин колец (The Lord of the Rings)", "Джон Рональд Руэл Толкин", 1954},
	}

	for _, e := range entries {
		b, err := book.NewBook(e.isbn, e.title, e.author, e.year)
		if err != nil {
			return nil, fmt.Errorf("не удалось создать книгу '%s': %w", e.title, err)
		}
		books = append(books, b)
	}

	return books, nil
}
