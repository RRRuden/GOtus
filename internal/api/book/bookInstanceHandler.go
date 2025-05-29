package handler

import (
	"encoding/json"
	"gotus/internal/model/book"
	"gotus/internal/repository"
	"net/http"
	"strconv"
	"strings"
)

type BookInstanceHandler struct {
	Repo *repository.BookInstanceRepository
}

func (h *BookInstanceHandler) CreateBookInstance(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ID   int    `json:"id"`
		ISBN string `json:"isbn"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	res := book.NewBookInstance(req.ID, req.ISBN)
	h.Repo.StoreBookInstance(res)
	w.WriteHeader(http.StatusCreated)
}

func (h *BookInstanceHandler) GetAllBookInstances(w http.ResponseWriter, r *http.Request) {
	res, _ := h.Repo.GetBookInstances()
	json.NewEncoder(w).Encode(res)
}

func (h *BookInstanceHandler) BookInstanceByIDHandler(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/bookinstance/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid reservation id", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.GetBookInstanceByID(w, r, id)
	case http.MethodPut:
		h.UpdateBookInstance(w, r, id)
	case http.MethodDelete:
		h.DeleteBookInstance(w, r, id)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *BookInstanceHandler) GetBookInstanceByID(w http.ResponseWriter, r *http.Request, id int) {
	res, ok := h.Repo.FindBookInstanceById(id)
	if !ok {
		http.NotFound(w, r)
		return
	}
	json.NewEncoder(w).Encode(res)
}

func (h *BookInstanceHandler) UpdateBookInstance(w http.ResponseWriter, r *http.Request, id int) {
	var req struct {
		ID   int    `json:"id"`
		ISBN string `json:"isbn"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	res := book.NewBookInstance(id, req.ISBN)
	if !h.Repo.UpdateBookInstanceById(id, res) {
		http.NotFound(w, r)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *BookInstanceHandler) DeleteBookInstance(w http.ResponseWriter, r *http.Request, id int) {
	if !h.Repo.DeleteBookInstanceById(id) {
		http.NotFound(w, r)
		return
	}
	w.WriteHeader(http.StatusOK)
}
