package service

import (
	"fmt"
	"time"

	"gotus/internal/model/book"
	"gotus/internal/model/reservation"
	"gotus/internal/model/user"
)

func GenerateAndStoreData(dataCh chan<- fmt.Stringer) {
	for i := 0; i < 5; i++ {
		start := time.Now()
		end := start.AddDate(0, 0, 14)

		isbn := fmt.Sprintf("978-3-16-148410-%d", i)
		b, _ := book.NewBook(isbn, "Капитанская дочка", "А. С. Пушкин", 1502)
		bi := book.NewBookInstance(i, b.GetISBN())
		u := user.NewUser(i, "Иван Иванов", "ivan@example.com")
		r := reservation.NewReservation(i, bi.GetID(), u.GetID(), i, start, end)

		dataCh <- b
		dataCh <- bi
		dataCh <- u
		dataCh <- r

		// Записываем данные раз в 3 секунды
		time.Sleep(3 * time.Second)
	}
}
