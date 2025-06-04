package handler

import (
	"encoding/json"
	"gotus/internal/model/reservation"
	"gotus/internal/repository"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type ReservationHandler struct {
	Repo *repository.ReservationRepository
}

type CreateReservationRequest struct {
	ID             int    `json:"id"`
	BookInstanceID int    `json:"book_instance_id"`
	UserID         int    `json:"user_id"`
	StatusID       int    `json:"status_id"`
	StartDate      string `json:"start_date"`
	EndDate        string `json:"end_date"`
}

type UpdateReservationRequest struct {
	BookInstanceID int    `json:"book_instance_id"`
	UserID         int    `json:"user_id"`
	StatusID       int    `json:"status_id"`
	StartDate      string `json:"start_date"`
	EndDate        string `json:"end_date"`
}

// CreateReservation godoc
// @Summary      Создать новую бронь
// @Description  Добавляет новую запись бронирования книги
// @Tags         reservation
// @Accept       json
// @Produce      json
// @Param        reservation body CreateReservationRequest true "Данные брони"
// @Success      201
// @Failure      400 {string} string "invalid request"
// @Router       /api/reservation [post]
func (h *ReservationHandler) CreateReservation(w http.ResponseWriter, r *http.Request) {
	var req CreateReservationRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	start, _ := time.Parse("2006-01-02", req.StartDate)
	end, _ := time.Parse("2006-01-02", req.EndDate)

	res := reservation.NewReservation(req.ID, req.BookInstanceID, req.UserID, req.StatusID, start, end)
	h.Repo.StoreReservation(res)
	w.WriteHeader(http.StatusCreated)
}

// GetAllReservations godoc
// @Summary      Получить список всех броней
// @Description  Возвращает все записи бронирования
// @Tags         reservation
// @Produce      json
// @Success      200 {array} reservation.Reservation
// @Router       /api/reservations [get]
func (h *ReservationHandler) GetAllReservations(w http.ResponseWriter, r *http.Request) {
	res, _ := h.Repo.GetReservations()
	json.NewEncoder(w).Encode(res)
}

func (h *ReservationHandler) ReservationByIDHandler(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/reservation/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid reservation id", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.GetReservationByID(w, r, id)
	case http.MethodPut:
		h.UpdateReservation(w, r, id)
	case http.MethodDelete:
		h.DeleteReservation(w, r, id)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

// GetReservationByID godoc
// @Summary      Получить бронь по ID
// @Description  Возвращает одну запись бронирования по ID
// @Tags         reservation
// @Produce      json
// @Param        id path int true "ID бронирования"
// @Success      200 {object} reservation.Reservation
// @Failure      404 {string} string "not found"
// @Router       /api/reservation/{id} [get]
func (h *ReservationHandler) GetReservationByID(w http.ResponseWriter, r *http.Request, id int) {
	res, ok := h.Repo.FindReservationById(id)
	if !ok {
		http.NotFound(w, r)
		return
	}
	json.NewEncoder(w).Encode(res)
}

// UpdateReservation godoc
// @Summary      Обновить бронь по ID
// @Description  Обновляет запись бронирования
// @Tags         reservation
// @Accept       json
// @Produce      json
// @Param        id path int true "ID бронирования"
// @Param        reservation body UpdateReservationRequest true "Обновлённые данные брони"
// @Success      200
// @Failure      400 {string} string "invalid request"
// @Failure      404 {string} string "not found"
// @Router       /api/reservation/{id} [put]
func (h *ReservationHandler) UpdateReservation(w http.ResponseWriter, r *http.Request, id int) {
	var req UpdateReservationRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	start, _ := time.Parse("2006-01-02", req.StartDate)
	end, _ := time.Parse("2006-01-02", req.EndDate)

	res := reservation.NewReservation(id, req.BookInstanceID, req.UserID, req.StatusID, start, end)
	if !h.Repo.UpdateReservationById(id, res) {
		http.NotFound(w, r)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// DeleteReservation godoc
// @Summary      Удалить бронь по ID
// @Description  Удаляет запись бронирования по ID
// @Tags         reservation
// @Param        id path int true "ID бронирования"
// @Success      200
// @Failure      404 {string} string "not found"
// @Router       /api/reservation/{id} [delete]
func (h *ReservationHandler) DeleteReservation(w http.ResponseWriter, r *http.Request, id int) {
	if !h.Repo.DeleteReservationById(id) {
		http.NotFound(w, r)
		return
	}
	w.WriteHeader(http.StatusOK)
}
