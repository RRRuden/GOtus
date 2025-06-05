package handler

import (
	"encoding/json"
	"gotus/internal/model/book"
	"gotus/internal/repository"
	"net/http"
	"strings"
)

type BookHandler struct {
	Repo *repository.BookRepository
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

// CreateBook godoc
// @Summary      Добавить новую книгу
// @Description  Создание новой книги и добавление её в хранилище
// @Tags         book
// @Accept       json
// @Produce      json
// @Param        book body CreateBookRequest true "Данные книги"
// @Success      201
// @Failure      400 {string} string "invalid request"
// @Router       /api/book [post]
func (h *BookHandler) CreateBook(w http.ResponseWriter, r *http.Request) {
	var req CreateBookRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	b, err := book.NewBook(req.ISBN, req.Title, req.Author, req.Year)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	h.Repo.StoreBook(b)
	w.WriteHeader(http.StatusCreated)
}

// GetAllBooks godoc
// @Summary      Получить список всех книг
// @Description  Возвращает массив всех книг в системе
// @Tags         book
// @Produce      json
// @Success      200 {array} book.Book
// @Router       /api/books [get]
func (h *BookHandler) GetAllBooks(w http.ResponseWriter, r *http.Request) {
	books, _ := h.Repo.GetBooks()
	json.NewEncoder(w).Encode(books)
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
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

// GetBookByISBN godoc
// @Summary      Получить книгу по ISBN
// @Description  Возвращает книгу по заданному ISBN
// @Tags         book
// @Produce      json
// @Param        isbn path string true "ISBN книги"
// @Success      200 {object} book.Book
// @Failure      404 {string} string "not found"
// @Router       /api/book/{isbn} [get]
func (h *BookHandler) GetBookByISBN(w http.ResponseWriter, r *http.Request, isbn string) {
	b, ok := h.Repo.FindBookByISBN(isbn)
	if !ok {
		http.NotFound(w, r)
		return
	}
	json.NewEncoder(w).Encode(b)
}

// UpdateBook godoc
// @Summary      Обновить книгу по ISBN
// @Description  Обновляет данные книги по заданному ISBN
// @Tags         book
// @Accept       json
// @Param        isbn path string true "ISBN книги"
// @Param        book body UpdateBookRequest true "Обновлённые данные книги"
// @Success      200
// @Failure      400 {string} string "invalid request"
// @Failure      404 {string} string "not found"
// @Router       /api/book/{isbn} [put]
func (h *BookHandler) UpdateBook(w http.ResponseWriter, r *http.Request, isbn string) {
	var req UpdateBookRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	b, err := book.NewBook(isbn, req.Title, req.Author, req.Year)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if !h.Repo.UpdateBookByISBN(isbn, b) {
		http.NotFound(w, r)
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
// @Failure      404 {string} string "not found"
// @Router       /api/book/{isbn} [delete]
func (h *BookHandler) DeleteBook(w http.ResponseWriter, r *http.Request, isbn string) {
	if !h.Repo.DeleteBookByISBN(isbn) {
		http.NotFound(w, r)
		return
	}
	w.WriteHeader(http.StatusOK)
}
