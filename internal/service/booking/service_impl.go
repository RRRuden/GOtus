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

// CreateBooking создаёт бронирование, если пользователь, книга и свободный экземпляр найдены
func (s *bookingService) CreateBooking(userID int, isbn string) (*reservation.Reservation, error) {
	_, found, err := s.UserRepo.FindUserById(userID)
	if err != nil {
		s.Logger.Log("CreateBooking", "ошибка при поиске пользователя: "+err.Error(), 3600)
		return nil, err
	}
	if !found {
		s.Logger.Log("CreateBooking", "пользователь не найден: userID="+strconv.Itoa(userID), 3600)
		return nil, ErrUserNotFound
	}

	_, found, err = s.BookRepo.FindBookByISBN(isbn)

	if err != nil {
		s.Logger.Log("CreateBooking", "ошибка при поиске книги: "+err.Error(), 3600)
		return nil, err
	}
	if !found {
		s.Logger.Log("CreateBooking", "книга не найдена: isbn="+isbn, 3600)
		return nil, ErrBookNotFound
	}

	instances, _, err := s.BookInstanceRepo.GetBookInstancesByISBN(isbn)
	if err != nil {
		s.Logger.Log("CreateBooking", "ошибка при получении экземпляров книги: "+err.Error(), 3600)
		return nil, err
	}

	for _, instance := range instances {
		hasActive, err := s.ReservationRepo.HasActiveReservation(instance.GetID())
		if err != nil {
			s.Logger.Log("CreateBooking", "ошибка проверки активных бронирований: "+err.Error(), 3600)
			return nil, err
		}
		if !hasActive {
			now := time.Now()
			newReservation := reservation.NewReservation(
				0,
				instance.GetID(),
				userID,
				int(reservation.StatusBooked),
				now,
				now.AddDate(0, 0, 7),
			)
			err = s.ReservationRepo.StoreReservation(newReservation)
			if err != nil {
				s.Logger.Log("CreateBooking", "ошибка при создании бронирования: "+err.Error(), 3600)
				return nil, err
			}
			s.Logger.Log("CreateBooking",
				fmt.Sprintf("создано бронирование userID=%d, isbn=%s, bookInstanceID=%d", userID, isbn, instance.GetID()), 86400)
			return newReservation, nil
		}
	}

	s.Logger.Log("CreateBooking", "нет доступных экземпляров для isbn="+isbn, 3600)
	return nil, ErrNoAvailableInstances
}

// ExtendBooking продлевает бронирование на ExtensionTime, max 7 дней
func (s *bookingService) ExtendBooking(reservationID int, extensionDays int) (bool, error) {
	if extensionDays <= 0 || extensionDays > 7 {
		s.Logger.Log("ExtendBooking",
			fmt.Sprintf("недопустимое время продления: reservationID=%d, extensionDays=%d", reservationID, extensionDays), 3600)
		return false, ErrExtensionTooLong
	}

	res, found, err := s.ReservationRepo.FindReservationById(reservationID)
	if err != nil {
		s.Logger.Log("ExtendBooking", "ошибка при поиске бронирования: "+err.Error(), 3600)
		return false, err
	}
	if !found || (res.ReservationStatusID == int(reservation.StatusCancelled) || res.ReservationStatusID == int(reservation.StatusEnded)) {
		s.Logger.Log("ExtendBooking",
			fmt.Sprintf("недопустимый статус бронирования для продления: reservationID=%d", reservationID), 3600)
		return false, ErrInvalidStatus
	}

	res.EndDate = res.EndDate.AddDate(0, 0, extensionDays)
	res.ReservationStatusID = int(reservation.StatusExtended)

	success, err := s.ReservationRepo.UpdateReservationById(reservationID, res)
	if err != nil {
		s.Logger.Log("ExtendBooking", "ошибка обновления бронирования: "+err.Error(), 3600)
		return false, err
	}

	if success {
		s.Logger.Log("ExtendBooking",
			fmt.Sprintf("продлено бронирование reservationID=%d на %d дней", reservationID, extensionDays), 86400)
	} else {
		s.Logger.Log("ExtendBooking",
			fmt.Sprintf("ошибка обновления бронирования reservationID=%d", reservationID), 3600)
	}
	return success, nil
}

// CancelBooking отменяет бронирование, если сегодня день начала
func (s *bookingService) CancelBooking(reservationID int) (bool, error) {
	res, found, err := s.ReservationRepo.FindReservationById(reservationID)
	if err != nil {
		s.Logger.Log("CancelBooking", "ошибка при поиске бронирования: "+err.Error(), 3600)
		return false, err
	}
	if !found {
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
	success, err := s.ReservationRepo.UpdateReservationById(reservationID, res)
	if err != nil {
		s.Logger.Log("CancelBooking", "ошибка обновления бронирования: "+err.Error(), 3600)
		return false, err
	}
	if success {
		s.Logger.Log("CancelBooking", fmt.Sprintf("отменено бронирование reservationID=%d", reservationID), 86400)
	} else {
		s.Logger.Log("CancelBooking", fmt.Sprintf("ошибка обновления бронирования reservationID=%d", reservationID), 3600)
	}
	return success, nil
}

// EndBooking завершает бронирование, если уже не день начала
func (s *bookingService) EndBooking(reservationID int) (bool, error) {
	res, found, err := s.ReservationRepo.FindReservationById(reservationID)
	if err != nil {
		s.Logger.Log("EndBooking", "ошибка при поиске бронирования: "+err.Error(), 3600)
		return false, err
	}
	if !found || res.ReservationStatusID == int(reservation.StatusCancelled) || res.ReservationStatusID == int(reservation.StatusEnded) {
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
	success, err := s.ReservationRepo.UpdateReservationById(reservationID, res)
	if err != nil {
		s.Logger.Log("EndBooking", "ошибка обновления бронирования: "+err.Error(), 3600)
		return false, err
	}
	if success {
		s.Logger.Log("EndBooking", fmt.Sprintf("завершено бронирование reservationID=%d", reservationID), 86400)
	} else {
		s.Logger.Log("EndBooking", fmt.Sprintf("ошибка обновления бронирования reservationID=%d", reservationID), 3600)
	}
	return success, nil
}
