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

func (h *BookHandler) CreateBook(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ISBN   string `json:"isbn"`
		Title  string `json:"title"`
		Author string `json:"author"`
		Year   int    `json:"year"`
	}

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

func (h *BookHandler) GetBookByISBN(w http.ResponseWriter, r *http.Request, isbn string) {
	b, ok := h.Repo.FindBookByISBN(isbn)
	if !ok {
		http.NotFound(w, r)
		return
	}
	json.NewEncoder(w).Encode(b)
}

func (h *BookHandler) UpdateBook(w http.ResponseWriter, r *http.Request, isbn string) {
	var req struct {
		Title  string `json:"title"`
		Author string `json:"author"`
		Year   int    `json:"year"`
	}
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

func (h *BookHandler) DeleteBook(w http.ResponseWriter, r *http.Request, isbn string) {
	if !h.Repo.DeleteBookByISBN(isbn) {
		http.NotFound(w, r)
		return
	}
	w.WriteHeader(http.StatusOK)
}
