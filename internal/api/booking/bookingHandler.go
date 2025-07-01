package handler

import (
	"context"
	"encoding/json"
	"gotus/internal/grpc/api/booking_api"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"google.golang.org/grpc"
)

type BookingHandler struct {
	grpcBookingClient booking_api.BookingServiceClient
}

func NewBookingHandler(grpcAddr string) *BookingHandler {
	connBooking, err := grpc.NewClient(grpcAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Не удалось подключиться к gRPC booking серверу: %v", err)
	}

	bookingClient := booking_api.NewBookingServiceClient(connBooking)

	return &BookingHandler{
		grpcBookingClient: bookingClient,
	}
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

// CreateBookingResponse представляет JSON-ответ с ID созданного бронирования.
type CreateBookingResponse struct {
	BookingID int32 `json:"booking_id"`
}

// CreateBooking godoc
// @Summary Создать бронирование книги
// @Description Бронирует доступный экземпляр книги по ISBN для указанного пользователя
// @Tags booking
// @Accept json
// @Produce json
// @Param request body CreateBookingRequest true "Данные для бронирования"
// @Success 201 {object} CreateBookingResponse "Бронирование успешно создано"
// @Failure 400 {string} string "Неверный запрос"
// @Failure 500 {string} string "Ошибка сервера"
// @Router /api/booking/create [post]
func (h *BookingHandler) CreateBooking(w http.ResponseWriter, r *http.Request) {
	var req CreateBookingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	resp, err := h.grpcBookingClient.CreateBooking(ctx, &booking_api.CreateBookingRequest{
		UserId: int32(req.UserID),
		Isbn:   req.ISBN,
	})
	if err != nil {
		http.Error(w, "Ошибка gRPC: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(CreateBookingResponse{BookingID: resp.BookingId})
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

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	_, err = h.grpcBookingClient.ExtendBooking(ctx, &booking_api.ExtendBookingRequest{
		BookingId:     int32(reservationID),
		ExtensionDays: int32(req.ExtensionDays),
	})
	if err != nil {
		http.Error(w, "gRPC error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Бронирование продлено"))
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

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	_, err = h.grpcBookingClient.CancelBooking(ctx, &booking_api.CancelBookingRequest{
		BookingId: int32(reservationID),
	})
	if err != nil {
		http.Error(w, "gRPC error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Бронирование отменено"))
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

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	_, err = h.grpcBookingClient.EndBooking(ctx, &booking_api.EndBookingRequest{
		BookingId: int32(reservationID),
	})
	if err != nil {
		http.Error(w, "gRPC error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Бронирование завершено"))
}
