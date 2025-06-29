package booking

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"gotus/internal/logger"
	"gotus/internal/model/reservation"
	"gotus/internal/repository"
)

type bookingService struct {
	UserRepo         repository.UserRepository
	BookRepo         repository.BookRepository
	BookInstanceRepo repository.BookInstanceRepository
	ReservationRepo  repository.ReservationRepository
	Logger           logger.Logger
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
	userRepo repository.UserRepository,
	bookRepo repository.BookRepository,
	bookInstanceRepo repository.BookInstanceRepository,
	reservationRepo repository.ReservationRepository,
	logger logger.Logger) Service {
	return &bookingService{
		UserRepo:         userRepo,
		BookRepo:         bookRepo,
		BookInstanceRepo: bookInstanceRepo,
		ReservationRepo:  reservationRepo,
		Logger:           logger,
	}
}

// 1. CreateBooking создаёт бронирование, если пользователь, книга и свободный экземпляр найдены
func (s *bookingService) CreateBooking(userID int, isbn string) (*reservation.Reservation, error) {
	// Проверка пользователя
	if _, ok := s.UserRepo.FindUserById(userID); !ok {
		s.Logger.Log("CreateBooking", "пользователь не найден: userID="+strconv.Itoa(userID), 3600)
		return nil, ErrUserNotFound
	}

	// Проверка книги
	if _, ok := s.BookRepo.FindBookByISBN(isbn); !ok {
		s.Logger.Log("CreateBooking", "книга не найдена: isbn="+isbn, 3600)
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
				0,
				instance.GetID(),
				userID,
				int(reservation.StatusBooked),
				now,
				now.AddDate(0, 0, 7), // стандартный срок 7 дней
			)
			s.ReservationRepo.StoreReservation(newReservation)
			s.Logger.Log("CreateBooking",
				fmt.Sprintf("создано бронирование userID=%d, isbn=%s, bookInstanceID=%d", userID, isbn, instance.GetID()), 86400)
			return newReservation, nil
		}
	}

	s.Logger.Log("CreateBooking", "нет доступных экземпляров для isbn="+isbn, 3600)
	return nil, ErrNoAvailableInstances
}

// 2. ExtendBooking продлевает бронирование на ExtensionTime, max 7 дней
func (s *bookingService) ExtendBooking(reservationID int, extensionDays int) (bool, error) {
	if extensionDays <= 0 || extensionDays > 7 {
		s.Logger.Log("ExtendBooking",
			fmt.Sprintf("недопустимое время продления: reservationID=%d, extensionDays=%d", reservationID, extensionDays), 3600)
		return false, ErrExtensionTooLong
	}

	res, ok := s.ReservationRepo.FindReservationById(reservationID)
	if !ok || (res.ReservationStatusID == int(reservation.StatusCancelled) || res.ReservationStatusID == int(reservation.StatusEnded)) {
		s.Logger.Log("ExtendBooking",
			fmt.Sprintf("недопустимый статус бронирования для продления: reservationID=%d", reservationID), 3600)
		return false, ErrInvalidStatus
	}

	res.EndDate = res.EndDate.AddDate(0, 0, extensionDays)
	res.ReservationStatusID = int(reservation.StatusExtended)
	success := s.ReservationRepo.UpdateReservationById(reservationID, res)
	if success {
		s.Logger.Log("ExtendBooking",
			fmt.Sprintf("продлено бронирование reservationID=%d на %d дней", reservationID, extensionDays), 86400)
	} else {
		s.Logger.Log("ExtendBooking",
			fmt.Sprintf("ошибка обновления бронирования reservationID=%d", reservationID), 3600)
	}
	return success, nil
}

// 3. CancelBooking отменяет бронирование, если сегодня день начала
func (s *bookingService) CancelBooking(reservationID int) (bool, error) {
	res, ok := s.ReservationRepo.FindReservationById(reservationID)
	if !ok {
		s.Logger.Log("CancelBooking", "бронирование не найдено: reservationID="+strconv.Itoa(reservationID), 3600)
		return false, ErrReservationNotFound
	}

	now := time.Now()
	startDate := res.StartDate

	if startDate.Year() != now.Year() || startDate.YearDay() != now.YearDay() {
		s.Logger.Log("CancelBooking",
			fmt.Sprintf("недопустимая дата отмены: reservationID=%d, startDate=%s, now=%s",
				reservationID, startDate.Format("2006-01-02"), now.Format("2006-01-02")), 3600)
		return false, ErrCancelDateMismatch
	}

	res.ReservationStatusID = int(reservation.StatusCancelled)
	success := s.ReservationRepo.UpdateReservationById(reservationID, res)
	if success {
		s.Logger.Log("CancelBooking", fmt.Sprintf("отменено бронирование reservationID=%d", reservationID), 86400)
	} else {
		s.Logger.Log("CancelBooking", fmt.Sprintf("ошибка обновления бронирования reservationID=%d", reservationID), 3600)
	}
	return success, nil
}

// 4. EndBooking завершает бронирование, если уже не день начала
func (s *bookingService) EndBooking(reservationID int) (bool, error) {
	res, ok := s.ReservationRepo.FindReservationById(reservationID)
	if !ok || (res.ReservationStatusID == reservation.StatusCancelled) || (res.ReservationStatusID == int(reservation.StatusEnded)) {
		s.Logger.Log("EndBooking", fmt.Sprintf("недопустимый статус для завершения: reservationID=%d", reservationID), 3600)
		return false, ErrInvalidStatus
	}

	now := time.Now()
	startDate := res.StartDate

	if startDate.Year() == now.Year() && startDate.YearDay() == now.YearDay() {
		s.Logger.Log("EndBooking", fmt.Sprintf("завершение в день начала запрещено: reservationID=%d", reservationID), 3600)
		return false, ErrEndDateMismatch
	}

	res.ReservationStatusID = int(reservation.StatusEnded)
	success := s.ReservationRepo.UpdateReservationById(reservationID, res)
	if success {
		s.Logger.Log("EndBooking", fmt.Sprintf("завершено бронирование reservationID=%d", reservationID), 86400)
	} else {
		s.Logger.Log("EndBooking", fmt.Sprintf("ошибка обновления бронирования reservationID=%d", reservationID), 3600)
	}
	return success, nil
}
