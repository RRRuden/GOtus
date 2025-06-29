package repository

import (
	"encoding/csv"
	"gotus/internal/model/book"
	"os"
	"path/filepath"
	"strconv"
	"sync"
)

type bookCSVRepository struct {
	books      []*book.Book
	dataDir    string
	filename   string
	booksMutex sync.Mutex
}

func NewBookCSVRepository(dataDir string) BookRepository {
	repo := &bookCSVRepository{
		books:    []*book.Book{},
		dataDir:  dataDir,
		filename: "books.csv",
	}

	repo.loadBooksFromCSV()
	return repo
}

func (r *bookCSVRepository) StoreBook(b *book.Book) {
	r.booksMutex.Lock()
	defer r.booksMutex.Unlock()
	r.books = append(r.books, b)
	r.saveBookToCSV(b)
}

func (r *bookCSVRepository) GetBooks() ([]*book.Book, int) {
	r.booksMutex.Lock()
	defer r.booksMutex.Unlock()
	return r.books, len(r.books)
}

func (r *bookCSVRepository) UpdateBookByISBN(isbn string, updatedBook *book.Book) bool {
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

func (r *bookCSVRepository) DeleteBookByISBN(isbn string) bool {
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

func (r *bookCSVRepository) FindBookByISBN(isbn string) (*book.Book, bool) {
	r.booksMutex.Lock()
	defer r.booksMutex.Unlock()
	for _, b := range r.books {
		if b.GetISBN() == isbn {
			return b, true
		}
	}
	return nil, false
}

func (r *bookCSVRepository) loadBooksFromCSV() {
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

func (r *bookCSVRepository) saveBookToCSV(b *book.Book) {
	file, _ := os.OpenFile(filepath.Join(r.dataDir, r.filename), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer file.Close()

	w := csv.NewWriter(file)
	defer w.Flush()

	_ = w.Write([]string{b.GetISBN(), b.Title, b.Author, strconv.Itoa(b.Year)})
}

func (r *bookCSVRepository) saveAllToCSV() {
	file, _ := os.Create(filepath.Join(r.dataDir, r.filename))
	defer file.Close()

	w := csv.NewWriter(file)
	defer w.Flush()

	for _, b := range r.books {
		_ = w.Write([]string{b.GetISBN(), b.Title, b.Author, strconv.Itoa(b.Year)})
	}
}
