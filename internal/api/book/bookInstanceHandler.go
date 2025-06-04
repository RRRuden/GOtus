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

type CreateBookInstanceRequest struct {
	ID   int    `json:"id"`
	ISBN string `json:"isbn"`
}

type UpdateBookInstanceRequest struct {
	ID   int    `json:"id"`
	ISBN string `json:"isbn"`
}

// CreateBookInstance godoc
// @Summary      Создать экземпляр книги
// @Description  Добавляет новый экземпляр книги по ISBN
// @Tags         bookinstance
// @Accept       json
// @Produce      json
// @Param        instance body CreateBookInstanceRequest true "Экземпляр книги"
// @Success      201
// @Failure      400 {string} string "invalid request"
// @Router       /api/bookinstance [post]
func (h *BookInstanceHandler) CreateBookInstance(w http.ResponseWriter, r *http.Request) {
	var req CreateBookInstanceRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	res := book.NewBookInstance(req.ID, req.ISBN)
	h.Repo.StoreBookInstance(res)
	w.WriteHeader(http.StatusCreated)
}

// GetAllBookInstances godoc
// @Summary      Получить все экземпляры книг
// @Description  Возвращает список всех экземпляров книг
// @Tags         bookinstance
// @Produce      json
// @Success      200 {array} book.BookInstance
// @Router       /api/bookinstances [get]
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

// GetBookInstanceByID godoc
// @Summary      Получить экземпляр книги по ID
// @Description  Возвращает экземпляр книги по ID
// @Tags         bookinstance
// @Produce      json
// @Param        id path int true "ID экземпляра книги"
// @Success      200 {object} book.BookInstance
// @Failure      404 {string} string "not found"
// @Router       /api/bookinstance/{id} [get]
func (h *BookInstanceHandler) GetBookInstanceByID(w http.ResponseWriter, r *http.Request, id int) {
	res, ok := h.Repo.FindBookInstanceById(id)
	if !ok {
		http.NotFound(w, r)
		return
	}
	json.NewEncoder(w).Encode(res)
}

// UpdateBookInstance godoc
// @Summary      Обновить экземпляр книги
// @Description  Обновляет экземпляр книги по ID
// @Tags         bookinstance
// @Accept       json
// @Produce      json
// @Param        id path int true "ID экземпляра книги"
// @Param        instance body UpdateBookInstanceRequest true "Новые данные экземпляра"
// @Success      200
// @Failure      400 {string} string "invalid request"
// @Failure      404 {string} string "not found"
// @Router       /api/bookinstance/{id} [put]
func (h *BookInstanceHandler) UpdateBookInstance(w http.ResponseWriter, r *http.Request, id int) {
	var req UpdateBookInstanceRequest

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

// DeleteBookInstance godoc
// @Summary      Удалить экземпляр книги
// @Description  Удаляет экземпляр книги по ID
// @Tags         bookinstance
// @Param        id path int true "ID экземпляра книги"
// @Success      200
// @Failure      404 {string} string "not found"
// @Router       /api/bookinstance/{id} [delete]
func (h *BookInstanceHandler) DeleteBookInstance(w http.ResponseWriter, r *http.Request, id int) {
	if !h.Repo.DeleteBookInstanceById(id) {
		http.NotFound(w, r)
		return
	}
	w.WriteHeader(http.StatusOK)
}
