// internal/repository/storage.go
package repository

import (
	"fmt"
	"sync"
	"time"

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

	prevCounts = struct {
		Books         int
		BookInstances int
		Users         int
		Reservations  int
	}{0, 0, 0, 0}
)

func Store(ch <-chan fmt.Stringer) {
	for data := range ch {
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
}

func LogUpdatesWorker() {
	ticker := time.NewTicker(200 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			booksMutex.Lock()
			if len(books) > prevCounts.Books {
				for _, b := range books[prevCounts.Books:] {
					fmt.Println("[Log]", b.String())
				}
				prevCounts.Books = len(books)
			}
			booksMutex.Unlock()

			bookInstancesMutex.Lock()
			if len(bookInstances) > prevCounts.BookInstances {
				for _, bi := range bookInstances[prevCounts.BookInstances:] {
					fmt.Println("[Log]", bi.String())
				}
				prevCounts.BookInstances = len(bookInstances)
			}
			bookInstancesMutex.Unlock()

			usersMutex.Lock()
			if len(users) > prevCounts.Users {
				for _, u := range users[prevCounts.Users:] {
					fmt.Println("[Log]", u.String())
				}
				prevCounts.Users = len(users)
			}
			usersMutex.Unlock()

			reservationsMutex.Lock()
			if len(reservations) > prevCounts.Reservations {
				for _, r := range reservations[prevCounts.Reservations:] {
					fmt.Println("[Log]", r.String())
				}
				prevCounts.Reservations = len(reservations)
			}
			reservationsMutex.Unlock()
		}
	}
}
