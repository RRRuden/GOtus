package handler

import (
	"encoding/json"
	"gotus/internal/service"
	"net/http"
	"strconv"
	"strings"
)

type BookingHandler struct {
	Service *service.BookingService
}

// CreateBookingRequest представляет JSON-запрос на бронирование книги.
type CreateBookingRequest struct {
	UserID int    `json:"user_id"`
	ISBN   string `json:"isbn"`
}

// ExtendBookingRequest представляет JSON-запрос на продление бронирования.
type ExtendBookingRequest struct {
	ExtensionDays int `json:"extension_days"`
}

// CreateBooking godoc
// @Summary Создать бронирование книги
// @Description Бронирует доступный экземпляр книги по ISBN для указанного пользователя
// @Tags booking
// @Accept json
// @Produce json
// @Param request body CreateBookingRequest true "Данные для бронирования"
// @Success 201 {string} string "Бронирование успешно создано"
// @Failure 400 {string} string "Неверный запрос"
// @Router /api/booking/create [post]
func (h *BookingHandler) CreateBooking(w http.ResponseWriter, r *http.Request) {
	var req CreateBookingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	_, err := h.Service.CreateBooking(req.UserID, req.ISBN)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// ExtendBooking godoc
// @Summary Продлить бронирование
// @Description Продлевает бронирование по ID на заданное количество дней
// @Tags booking
// @Accept json
// @Produce json
// @Param id path int true "ID бронирования"
// @Param request body ExtendBookingRequest true "Дни продления"
// @Success 200 {string} string "Бронирование продлено"
// @Failure 400 {string} string "Ошибка запроса"
// @Router /api/booking/extend/{id} [put]
func (h *BookingHandler) ExtendBooking(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/booking/extend/")
	reservationID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid reservation id", http.StatusBadRequest)
		return
	}

	var req ExtendBookingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	_, err = h.Service.ExtendBooking(reservationID, req.ExtensionDays)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// CancelBooking godoc
// @Summary Отменить бронирование
// @Description Отменяет бронирование по ID
// @Tags booking
// @Produce json
// @Param id path int true "ID бронирования"
// @Success 200 {string} string "Бронирование отменено"
// @Failure 400 {string} string "Ошибка запроса"
// @Router /api/booking/cancel/{id} [post]
func (h *BookingHandler) CancelBooking(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/booking/cancel/")
	reservationID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid reservation id", http.StatusBadRequest)
		return
	}

	_, err = h.Service.CancelBooking(reservationID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// EndBooking godoc
// @Summary Завершить бронирование
// @Description Завершает бронирование по ID
// @Tags booking
// @Produce json
// @Param id path int true "ID бронирования"
// @Success 200 {string} string "Бронирование завершено"
// @Failure 400 {string} string "Ошибка запроса"
// @Router /api/booking/end/{id} [post]
func (h *BookingHandler) EndBooking(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/booking/end/")
	reservationID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid reservation id", http.StatusBadRequest)
		return
	}

	_, err = h.Service.EndBooking(reservationID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}
