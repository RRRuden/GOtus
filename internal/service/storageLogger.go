package service

import (
	"context"
	"log"
	"time"

	"gotus/internal/repository"
)

func LogUpdatesWorker(ctx context.Context) {
	var lastBookCount, lastBookInstanceCount, lastUserCount, lastReservationCount int

	ticker := time.NewTicker(200 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("[StorageLogger] Graceful shutdown")
			return
		default:
			// Книги
			if books, count := repository.GetBooks(); count != lastBookCount {
				newBooks := books[lastBookCount:]
				for _, b := range newBooks {
					log.Println("[Book]", b.String())
				}
				lastBookCount = count
			}

			// Экземпляры книг
			if bookInstances, count := repository.GetBookInstances(); count != lastBookInstanceCount {
				newBookInstances := bookInstances[lastBookInstanceCount:]
				for _, bi := range newBookInstances {
					log.Println("[BookInstance]", bi.String())
				}
				lastBookInstanceCount = count
			}

			// Пользователи
			if users, count := repository.GetUsers(); count != lastUserCount {
				newUsers := users[lastUserCount:]
				for _, u := range newUsers {
					log.Println("[User]", u.String())
				}
				lastUserCount = count
			}

			// Бронирования
			if reservations, count := repository.GetReservations(); count != lastReservationCount {
				newReservations := reservations[lastReservationCount:]
				for _, r := range newReservations {
					log.Println("[Reservation]", r.String())
				}
				lastReservationCount = count
			}
		}
	}
}
