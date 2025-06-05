package service

import (
	"errors"
	"time"

	"gotus/internal/model/reservation"
	"gotus/internal/repository"
)

type BookingService struct {
	UserRepo         *repository.UserRepository
	BookRepo         *repository.BookRepository
	BookInstanceRepo *repository.BookInstanceRepository
	ReservationRepo  *repository.ReservationRepository
}

var (
	ErrUserNotFound         = errors.New("пользователь не найден")
	ErrBookNotFound         = errors.New("книга не найдена")
	ErrNoAvailableInstances = errors.New("нет доступных экземпляров книги")
	ErrReservationNotFound  = errors.New("бронирование не найдено")
	ErrInvalidStatus        = errors.New("недопустимый статус бронирования")
	ErrExtensionTooLong     = errors.New("время продления не должно быть больше 7 дней")
	ErrCancelDateMismatch   = errors.New("бронирование можно отменить только в день начала")
	ErrEndDateMismatch      = errors.New("бронирование можно завершить только после дня начала")
)

func NewBookingService(
	userRepo *repository.UserRepository,
	bookRepo *repository.BookRepository,
	bookInstanceRepo *repository.BookInstanceRepository,
	reservationRepo *repository.ReservationRepository,
) *BookingService {
	return &BookingService{
		UserRepo:         userRepo,
		BookRepo:         bookRepo,
		BookInstanceRepo: bookInstanceRepo,
		ReservationRepo:  reservationRepo,
	}
}

// 1. CreateBooking создаёт бронирование, если пользователь, книга и свободный экземпляр найдены
func (s *BookingService) CreateBooking(userID int, isbn string) (*reservation.Reservation, error) {
	// Проверка пользователя
	if _, ok := s.UserRepo.FindUserById(userID); !ok {
		return nil, ErrUserNotFound
	}

	// Проверка книги
	if _, ok := s.BookRepo.FindBookByISBN(isbn); !ok {
		return nil, ErrBookNotFound
	}

	// Поиск свободных экземпляров книги
	instances, _ := s.BookInstanceRepo.GetBookInstancesByISBN(isbn)
	for _, instance := range instances {
		// Проверяем, что для данного экземпляра нет активного бронирования со статусом Забронирована или Продлена
		if !s.ReservationRepo.HasActiveReservation(instance.GetID()) {
			// Создаем бронирование
			now := time.Now()
			newReservation := reservation.NewReservation(
				s.ReservationRepo.GenerateID(),
				instance.GetID(),
				userID,
				int(reservation.StatusBooked),
				now,
				now.AddDate(0, 0, 7), // стандартный срок 7 дней
			)
			s.ReservationRepo.StoreReservation(newReservation)
			return newReservation, nil
		}
	}

	return nil, ErrNoAvailableInstances
}

// 2. ExtendBooking продлевает бронирование на ExtensionTime, max 7 дней
func (s *BookingService) ExtendBooking(reservationID int, extensionDays int) (bool, error) {
	if extensionDays <= 0 || extensionDays > 7 {
		return false, ErrExtensionTooLong
	}

	res, ok := s.ReservationRepo.FindReservationById(reservationID)
	if !ok || (res.ReservationStatusID != int(reservation.StatusBooked) && res.ReservationStatusID != int(reservation.StatusEnded)) {
		return false, ErrInvalidStatus
	}

	// Увеличиваем EndDate
	res.EndDate = res.EndDate.AddDate(0, 0, extensionDays)
	res.ReservationStatusID = int(reservation.StatusExtended)
	return s.ReservationRepo.UpdateReservationById(reservationID, res), nil
}

// 3. CancelBooking отменяет бронирование, если сегодня день начала
func (s *BookingService) CancelBooking(reservationID int) (bool, error) {
	res, ok := s.ReservationRepo.FindReservationById(reservationID)
	if !ok {
		return false, ErrReservationNotFound
	}

	now := time.Now()
	startDate := res.StartDate

	if startDate.Year() != now.Year() || startDate.YearDay() != now.YearDay() {
		return false, ErrCancelDateMismatch
	}

	res.ReservationStatusID = int(reservation.StatusCancelled)
	return s.ReservationRepo.UpdateReservationById(reservationID, res), nil
}

// 4. EndBooking завершает бронирование, если уже не день начала
func (s *BookingService) EndBooking(reservationID int) (bool, error) {
	res, ok := s.ReservationRepo.FindReservationById(reservationID)
	if !ok || (res.ReservationStatusID != int(reservation.StatusCancelled) && res.ReservationStatusID != int(reservation.StatusEnded)) {
		return false, ErrInvalidStatus
	}

	now := time.Now()
	startDate := res.StartDate

	if startDate.Year() == now.Year() && startDate.YearDay() == now.YearDay() {
		return false, ErrEndDateMismatch
	}

	res.ReservationStatusID = int(reservation.StatusEnded)
	return s.ReservationRepo.UpdateReservationById(reservationID, res), nil
}
