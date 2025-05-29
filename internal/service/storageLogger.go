package service

import (
	"context"
	"log"
	"time"

	"gotus/internal/repository"
)

func LogUpdatesWorker(ctx context.Context, s *repository.Storage) {
	var lastBookCount, lastBookInstanceCount, lastUserCount, lastReservationCount int

	ticker := time.NewTicker(200 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("[StorageLogger] Graceful shutdown")
			return
		case <-ticker.C:
			// Книги
			if books, count := s.BookRepository.GetBooks(); count != lastBookCount {
				newBooks := books[lastBookCount:]
				for _, b := range newBooks {
					log.Println("[Book]", b.String())
				}
				lastBookCount = count
			}

			// Экземпляры книг
			if bookInstances, count := s.BookInstanceRepository.GetBookInstances(); count != lastBookInstanceCount {
				newBookInstances := bookInstances[lastBookInstanceCount:]
				for _, bi := range newBookInstances {
					log.Println("[BookInstance]", bi.String())
				}
				lastBookInstanceCount = count
			}

			// Пользователи
			if users, count := s.UserRepository.GetUsers(); count != lastUserCount {
				newUsers := users[lastUserCount:]
				for _, u := range newUsers {
					log.Println("[User]", u.String())
				}
				lastUserCount = count
			}

			// Бронирования
			if reservations, count := s.ReservationRepository.GetReservations(); count != lastReservationCount {
				newReservations := reservations[lastReservationCount:]
				for _, r := range newReservations {
					log.Println("[Reservation]", r.String())
				}
				lastReservationCount = count
			}
		}
	}
}
