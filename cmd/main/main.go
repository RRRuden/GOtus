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
	b, err := book.NewBook("978-3-16-148410-0", "Капитанская дочка", "А. С. Пушкин", 1502)

	if err != nil {
		fmt.Println("Ошибка при создании книги:", err)
		return
	}

	fmt.Println("Создана книга:", b.Name, "ID:", b.GetISBN(), "Автор произведения:", b.Author)

	//Создание экземпляра книги
	instance := book.NewBookInstance(1, b.GetISBN())
	fmt.Println("Создан экземпляр книги:", instance.GetID(), "ISBN:", instance.ISBN)

	// Создание пользователя
	u := user.NewUser(1, "Иван Иванов", "ivan@example.com")
	fmt.Println("Создан пользователь:", u.Name, "Email:", u.Email, "ID:", u.GetID())

	// Создание бронирования
	startDate := time.Now()
	endDate := startDate.AddDate(0, 0, 14)
	r := reservation.NewReservation(1, instance.GetID(), u.GetID(), 1, startDate, endDate)
	fmt.Println("Создано бронирование ID:", r.GetID(), "Номер экземпляра:", r.BookInstanceID, "Пользователь ID:", r.UserID)

	// Отправка сообщения
	msg := message.NewMessage(1, u.Email, "Подтверждение бронирования", time.Now())
	fmt.Println("Сообщение отправлено на:", msg.Email, "Тема:", msg.EmailSubject)
}
