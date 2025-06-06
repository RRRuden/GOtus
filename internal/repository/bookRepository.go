package repository

import (
	"encoding/csv"
	"gotus/internal/model/book"
	"os"
	"path/filepath"
	"strconv"
	"sync"
)

type BookRepository struct {
	books      []*book.Book
	dataDir    string
	filename   string
	booksMutex sync.Mutex
}

func NewBookRepository(dataDir string) *BookRepository {
	return &BookRepository{
		books:    []*book.Book{},
		dataDir:  dataDir,
		filename: "books.csv",
	}
}

func (r *BookRepository) StoreBook(b *book.Book) {
	r.booksMutex.Lock()
	defer r.booksMutex.Unlock()
	r.books = append(r.books, b)
	r.saveBookToCSV(b)
}

func (r *BookRepository) GetBooks() ([]*book.Book, int) {
	r.booksMutex.Lock()
	defer r.booksMutex.Unlock()
	return r.books, len(r.books)
}

func (r *BookRepository) LoadBooksFromCSV() {
	file, err := os.Open(filepath.Join(r.dataDir, r.filename))
	if err != nil {
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, _ := reader.ReadAll()

	r.booksMutex.Lock()
	defer r.booksMutex.Unlock()

	for _, rec := range records {
		year, _ := strconv.Atoi(rec[3])
		b, _ := book.NewBook(rec[0], rec[1], rec[2], year)
		r.books = append(r.books, b)
	}
}

func (r *BookRepository) UpdateBookByISBN(isbn string, updatedBook *book.Book) bool {
	r.booksMutex.Lock()
	defer r.booksMutex.Unlock()

	found := false
	for i, b := range r.books {
		if b.GetISBN() == isbn {
			r.books[i] = updatedBook
			found = true
			break
		}
	}

	if found {
		r.saveAllToCSV()
	}

	return found
}

func (r *BookRepository) DeleteBookByISBN(isbn string) bool {
	r.booksMutex.Lock()
	defer r.booksMutex.Unlock()

	found := false
	var updated []*book.Book
	for _, b := range r.books {
		if b.GetISBN() != isbn {
			updated = append(updated, b)
		} else {
			found = true
		}
	}

	if found {
		r.books = updated
		r.saveAllToCSV()
	}

	return found
}

func (r *BookRepository) FindBookByISBN(isbn string) (*book.Book, bool) {
	r.booksMutex.Lock()
	defer r.booksMutex.Unlock()
	for _, b := range r.books {
		if b.GetISBN() == isbn {
			return b, true
		}
	}
	return nil, false
}

func (r *BookRepository) saveBookToCSV(b *book.Book) {
	file, _ := os.OpenFile(filepath.Join(r.dataDir, r.filename), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer file.Close()

	w := csv.NewWriter(file)
	defer w.Flush()

	_ = w.Write([]string{b.GetISBN(), b.Title, b.Author, strconv.Itoa(b.Year)})
}

func (r *BookRepository) saveAllToCSV() {
	file, _ := os.Create(filepath.Join(r.dataDir, r.filename))
	defer file.Close()

	w := csv.NewWriter(file)
	defer w.Flush()

	for _, b := range r.books {
		_ = w.Write([]string{b.GetISBN(), b.Title, b.Author, strconv.Itoa(b.Year)})
	}
}
