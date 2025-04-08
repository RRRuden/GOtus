// internal/service/generator.go
package service

import (
	"time"

	"gotus/internal/model/book"
	"gotus/internal/model/reservation"
	"gotus/internal/model/user"
	"gotus/internal/repository"
)

func GenerateAndStoreData() {
	start := time.Now()
	end := start.AddDate(0, 0, 14)

	b, _ := book.NewBook("978-3-16-148410-0", "Капитанская дочка", "А. С. Пушкин", 1502)
	bi := book.NewBookInstance(1, b.GetISBN())
	u := user.NewUser(1, "Иван Иванов", "ivan@example.com")
	r := reservation.NewReservation(1, bi.GetID(), u.GetID(), 1, start, end)

	repository.Store(b)
	repository.Store(bi)
	repository.Store(u)
	repository.Store(r)
}
