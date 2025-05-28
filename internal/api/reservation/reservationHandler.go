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

func (h *ReservationHandler) CreateReservation(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ID             int    `json:"id"`
		BookInstanceID int    `json:"book_instance_id"`
		UserID         int    `json:"user_id"`
		StatusID       int    `json:"status_id"`
		StartDate      string `json:"start_date"`
		EndDate        string `json:"end_date"`
	}

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

func (h *ReservationHandler) GetReservationByID(w http.ResponseWriter, r *http.Request, id int) {
	res, ok := h.Repo.FindReservationById(id)
	if !ok {
		http.NotFound(w, r)
		return
	}
	json.NewEncoder(w).Encode(res)
}

func (h *ReservationHandler) UpdateReservation(w http.ResponseWriter, r *http.Request, id int) {
	var req struct {
		BookInstanceID int    `json:"book_instance_id"`
		UserID         int    `json:"user_id"`
		StatusID       int    `json:"status_id"`
		StartDate      string `json:"start_date"`
		EndDate        string `json:"end_date"`
	}

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

func (h *ReservationHandler) DeleteReservation(w http.ResponseWriter, r *http.Request, id int) {
	if !h.Repo.DeleteReservationById(id) {
		http.NotFound(w, r)
		return
	}
	w.WriteHeader(http.StatusOK)
}
