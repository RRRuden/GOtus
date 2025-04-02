// main.go
package main

import (
	"fmt"
	"time"

	"gotus/internal/model/book"
	"gotus/internal/model/message"
	"gotus/internal/model/reservation"
	"gotus/internal/model/user"
)

func main() {
	// Создание книги
	b := book.NewBook(1, "1984")
	fmt.Println("Создана книга:", b.Name, "ID:", b.GetID())

	// Создание пользователя
	u := user.NewUser(1, "Иван Иванов", "ivan@example.com")
	fmt.Println("Создан пользователь:", u.Name, "Email:", u.Email, "ID:", u.GetID())

	// Создание бронирования
	startDate := time.Now()
	endDate := startDate.AddDate(0, 0, 14)
	r := reservation.NewReservation(1, b.GetID(), u.GetID(), 1, startDate, endDate)
	fmt.Println("Создано бронирование ID:", r.GetID(), "Книга ID:", r.BookID, "Пользователь ID:", r.UserID)

	// Отправка сообщения
	msg := message.NewMessage(1, u.Email, "Подтверждение бронирования", time.Now())
	fmt.Println("Сообщение отправлено на:", msg.Email, "Тема:", msg.EmailSubject)
}
