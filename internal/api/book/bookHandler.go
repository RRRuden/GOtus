package handler

import (
	"encoding/json"
	"gotus/internal/model/book"
	"gotus/internal/repository"
	"net/http"
	"strings"
)

type BookHandler struct {
	Repo repository.BookRepository
}

type CreateBookRequest struct {
	ISBN   string `json:"isbn"`
	Title  string `json:"title"`
	Author string `json:"author"`
	Year   int    `json:"year"`
}

type UpdateBookRequest struct {
	Title  string `json:"title"`
	Author string `json:"author"`
	Year   int    `json:"year"`
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

// CreateBook godoc
// @Summary      Добавить новую книгу
// @Description  Создание новой книги и добавление её в хранилище
// @Tags         book
// @Accept       json
// @Produce      json
// @Param        book body CreateBookRequest true "Данные книги"
// @Success      201
// @Failure      400 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Router       /api/book [post]
func (h *BookHandler) CreateBook(w http.ResponseWriter, r *http.Request) {
	var req CreateBookRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Неверный формат запроса")
		return
	}

	b, err := book.NewBook(req.ISBN, req.Title, req.Author, req.Year)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.Repo.StoreBook(b); err != nil {
		writeError(w, http.StatusInternalServerError, "Ошибка при сохранении книги")
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// GetAllBooks godoc
// @Summary      Получить список всех книг
// @Description  Возвращает массив всех книг в системе
// @Tags         book
// @Produce      json
// @Success      200 {array} book.Book
// @Failure      500 {object} ErrorResponse
// @Router       /api/books [get]
func (h *BookHandler) GetAllBooks(w http.ResponseWriter, r *http.Request) {
	books, _, err := h.Repo.GetBooks()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Ошибка при получении списка книг")
		return
	}
	writeJSON(w, books)
}

func (h *BookHandler) BookByISBNHandler(w http.ResponseWriter, r *http.Request) {
	isbn := strings.TrimPrefix(r.URL.Path, "/api/book/")
	switch r.Method {
	case http.MethodGet:
		h.GetBookByISBN(w, r, isbn)
	case http.MethodPut:
		h.UpdateBook(w, r, isbn)
	case http.MethodDelete:
		h.DeleteBook(w, r, isbn)
	default:
		writeError(w, http.StatusMethodNotAllowed, "Метод не поддерживается")
	}
}

// GetBookByISBN godoc
// @Summary      Получить книгу по ISBN
// @Description  Возвращает книгу по заданному ISBN
// @Tags         book
// @Produce      json
// @Param        isbn path string true "ISBN книги"
// @Success      200 {object} book.Book
// @Failure      404 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Router       /api/book/{isbn} [get]
func (h *BookHandler) GetBookByISBN(w http.ResponseWriter, r *http.Request, isbn string) {
	b, _, err := h.Repo.FindBookByISBN(isbn)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Ошибка при поиске книги")
		return
	}
	if b == nil {
		writeError(w, http.StatusNotFound, "Книга не найдена")
		return
	}
	writeJSON(w, b)
}

// UpdateBook godoc
// @Summary      Обновить книгу по ISBN
// @Description  Обновляет данные книги по заданному ISBN
// @Tags         book
// @Accept       json
// @Param        isbn path string true "ISBN книги"
// @Param        book body UpdateBookRequest true "Обновлённые данные книги"
// @Success      200
// @Failure      400 {object} ErrorResponse
// @Failure      404 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Router       /api/book/{isbn} [put]
func (h *BookHandler) UpdateBook(w http.ResponseWriter, r *http.Request, isbn string) {
	var req UpdateBookRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Неверный формат запроса")
		return
	}

	b, err := book.NewBook(isbn, req.Title, req.Author, req.Year)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	ok, err := h.Repo.UpdateBookByISBN(isbn, b)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Ошибка при обновлении книги")
		return
	}
	if !ok {
		writeError(w, http.StatusNotFound, "Книга не найдена")
		return
	}

	w.WriteHeader(http.StatusOK)
}

// DeleteBook godoc
// @Summary      Удалить книгу по ISBN
// @Description  Удаляет книгу из хранилища по заданному ISBN
// @Tags         book
// @Param        isbn path string true "ISBN книги"
// @Success      200
// @Failure      404 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Router       /api/book/{isbn} [delete]
func (h *BookHandler) DeleteBook(w http.ResponseWriter, r *http.Request, isbn string) {
	ok, err := h.Repo.DeleteBookByISBN(isbn)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Ошибка при удалении книги")
		return
	}
	if !ok {
		writeError(w, http.StatusNotFound, "Книга не найдена")
		return
	}
	w.WriteHeader(http.StatusOK)
}
