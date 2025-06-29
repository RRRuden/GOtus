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
	Repo repository.BookInstanceRepository
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
// @Failure      400 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Router       /api/bookinstance [post]
func (h *BookInstanceHandler) CreateBookInstance(w http.ResponseWriter, r *http.Request) {
	var req CreateBookInstanceRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Неверный формат запроса")
		return
	}

	res := book.NewBookInstance(req.ID, req.ISBN)
	if err := h.Repo.StoreBookInstance(res); err != nil {
		writeError(w, http.StatusInternalServerError, "Ошибка при создании экземпляра книги")
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// GetAllBookInstances godoc
// @Summary      Получить все экземпляры книг
// @Description  Возвращает список всех экземпляров книг
// @Tags         bookinstance
// @Produce      json
// @Success      200 {array} book.BookInstance
// @Failure      500 {object} ErrorResponse
// @Router       /api/bookinstances [get]
func (h *BookInstanceHandler) GetAllBookInstances(w http.ResponseWriter, r *http.Request) {
	res, _, err := h.Repo.GetBookInstances()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Ошибка при получении экземпляров книг")
		return
	}
	writeJSON(w, res)
}

func (h *BookInstanceHandler) BookInstanceByIDHandler(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/bookinstance/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Некорректный ID")
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
		writeError(w, http.StatusMethodNotAllowed, "Метод не поддерживается")
	}
}

// GetBookInstanceByID godoc
// @Summary      Получить экземпляр книги по ID
// @Description  Возвращает экземпляр книги по ID
// @Tags         bookinstance
// @Produce      json
// @Param        id path int true "ID экземпляра книги"
// @Success      200 {object} book.BookInstance
// @Failure      404 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Router       /api/bookinstance/{id} [get]
func (h *BookInstanceHandler) GetBookInstanceByID(w http.ResponseWriter, r *http.Request, id int) {
	res, _, err := h.Repo.FindBookInstanceById(id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Ошибка при поиске экземпляра книги")
		return
	}
	if res == nil {
		writeError(w, http.StatusNotFound, "Экземпляр книги не найден")
		return
	}
	writeJSON(w, res)
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
// @Failure      400 {object} ErrorResponse
// @Failure      404 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Router       /api/bookinstance/{id} [put]
func (h *BookInstanceHandler) UpdateBookInstance(w http.ResponseWriter, r *http.Request, id int) {
	var req UpdateBookInstanceRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Неверный формат запроса")
		return
	}

	res := book.NewBookInstance(id, req.ISBN)
	ok, err := h.Repo.UpdateBookInstanceById(id, res)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Ошибка при обновлении экземпляра книги")
		return
	}
	if !ok {
		writeError(w, http.StatusNotFound, "Экземпляр книги не найден")
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
// @Failure      404 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Router       /api/bookinstance/{id} [delete]
func (h *BookInstanceHandler) DeleteBookInstance(w http.ResponseWriter, r *http.Request, id int) {
	ok, err := h.Repo.DeleteBookInstanceById(id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Ошибка при удалении экземпляра книги")
		return
	}
	if !ok {
		writeError(w, http.StatusNotFound, "Экземпляр книги не найден")
		return
	}
	w.WriteHeader(http.StatusOK)
}
