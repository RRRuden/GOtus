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
	Repo repository.ReservationRepository
}

type CreateReservationRequest struct {
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

type ErrorResponse struct {
	Message string `json:"message"`
}

func writeError(w http.ResponseWriter, status int, msg string) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(ErrorResponse{Message: msg})
}

func writeJSON(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

// CreateReservation godoc
// @Summary      Создать новую бронь
// @Description  Добавляет новую запись бронирования книги
// @Tags         reservation
// @Accept       json
// @Produce      json
// @Param        reservation body CreateReservationRequest true "Данные брони"
// @Success      201
// @Failure      400 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Router       /api/reservation [post]
func (h *ReservationHandler) CreateReservation(w http.ResponseWriter, r *http.Request) {
	var req CreateReservationRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Неверный формат запроса")
		return
	}

	start, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Неверная дата начала")
		return
	}
	end, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Неверная дата окончания")
		return
	}

	res := reservation.NewReservation(0, req.BookInstanceID, req.UserID, req.StatusID, start, end)
	if _, err := h.Repo.StoreReservation(res); err != nil {
		writeError(w, http.StatusInternalServerError, "Ошибка при сохранении брони")
		return
	}
	w.WriteHeader(http.StatusCreated)
}

// GetAllReservations godoc
// @Summary      Получить список всех броней
// @Description  Возвращает все записи бронирования
// @Tags         reservation
// @Produce      json
// @Success      200 {array} reservation.Reservation
// @Failure      500 {object} ErrorResponse
// @Router       /api/reservations [get]
func (h *ReservationHandler) GetAllReservations(w http.ResponseWriter, r *http.Request) {
	res, _, err := h.Repo.GetReservations()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Ошибка при получении данных")
		return
	}
	writeJSON(w, res)
}

func (h *ReservationHandler) ReservationByIDHandler(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/reservation/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Некорректный ID")
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
		writeError(w, http.StatusMethodNotAllowed, "Метод не поддерживается")
	}
}

// GetReservationByID godoc
// @Summary      Получить бронь по ID
// @Description  Возвращает одну запись бронирования по ID
// @Tags         reservation
// @Produce      json
// @Param        id path int true "ID бронирования"
// @Success      200 {object} reservation.Reservation
// @Failure      404 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Router       /api/reservation/{id} [get]
func (h *ReservationHandler) GetReservationByID(w http.ResponseWriter, r *http.Request, id int) {
	res, _, err := h.Repo.FindReservationById(id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Ошибка при получении брони")
		return
	}
	if res == nil {
		writeError(w, http.StatusNotFound, "Бронь не найдена")
		return
	}
	writeJSON(w, res)
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
// @Failure      400 {object} ErrorResponse
// @Failure      404 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Router       /api/reservation/{id} [put]
func (h *ReservationHandler) UpdateReservation(w http.ResponseWriter, r *http.Request, id int) {
	var req UpdateReservationRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Неверный формат запроса")
		return
	}

	start, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Неверная дата начала")
		return
	}
	end, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Неверная дата окончания")
		return
	}

	res := reservation.NewReservation(id, req.BookInstanceID, req.UserID, req.StatusID, start, end)
	ok, err := h.Repo.UpdateReservationById(id, res)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Ошибка при обновлении брони")
		return
	}
	if !ok {
		writeError(w, http.StatusNotFound, "Бронь не найдена")
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
// @Failure      404 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Router       /api/reservation/{id} [delete]
func (h *ReservationHandler) DeleteReservation(w http.ResponseWriter, r *http.Request, id int) {
	ok, err := h.Repo.DeleteReservationById(id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Ошибка при удалении брони")
		return
	}
	if !ok {
		writeError(w, http.StatusNotFound, "Бронь не найдена")
		return
	}
	w.WriteHeader(http.StatusOK)
}
