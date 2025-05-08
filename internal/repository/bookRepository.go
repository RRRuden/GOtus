package repository

import (
	"encoding/csv"
	"gotus/internal/model/book"
	"os"
	"strconv"
	"sync"
)

const bookCSVPath = "./data/books.csv"

var (
	books      []*book.Book
	booksMutex sync.Mutex
)

func StoreBook(b *book.Book) {
	booksMutex.Lock()
	defer booksMutex.Unlock()
	books = append(books, b)
	saveBookToCSV(b)
}

func GetBooks() ([]*book.Book, int) {
	booksMutex.Lock()
	defer booksMutex.Unlock()
	return books, len(books)
}

func LoadBooksFromCSV() {
	file, err := os.Open(bookCSVPath)
	if err != nil {
		return
	}
	defer file.Close()

	r := csv.NewReader(file)
	records, _ := r.ReadAll()

	booksMutex.Lock()
	defer booksMutex.Unlock()

	for _, rec := range records {
		year, _ := strconv.Atoi(rec[3])
		b, _ := book.NewBook(rec[0], rec[1], rec[2], year)
		books = append(books, b)
	}
}

func saveBookToCSV(b *book.Book) {
	file, _ := os.OpenFile(bookCSVPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer file.Close()

	w := csv.NewWriter(file)
	defer w.Flush()

	_ = w.Write([]string{b.GetISBN(), b.Title, b.Author, strconv.Itoa(b.Year)})
}
