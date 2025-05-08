package repository

import (
	"encoding/csv"
	"gotus/internal/model/book"
	"os"
	"strconv"
	"sync"
)

const bookInstanceCSVPath = "./data/bookInstances.csv"

var (
	bookInstances      []*book.BookInstance
	bookInstancesMutex sync.Mutex
)

func StoreBookInstance(bi *book.BookInstance) {
	bookInstancesMutex.Lock()
	defer bookInstancesMutex.Unlock()

	bookInstances = append(bookInstances, bi)
	saveBookInstanceToCSV(bi)
}

func GetBookInstances() ([]*book.BookInstance, int) {
	bookInstancesMutex.Lock()
	defer bookInstancesMutex.Unlock()

	return bookInstances, len(bookInstances)
}

func LoadBookInstancesFromCSV() {
	file, err := os.Open(bookInstanceCSVPath)
	if err != nil {
		return
	}
	defer file.Close()

	r := csv.NewReader(file)
	records, _ := r.ReadAll()

	bookInstancesMutex.Lock()
	defer bookInstancesMutex.Unlock()

	for _, rec := range records {
		id, _ := strconv.Atoi(rec[0])
		bi := book.NewBookInstance(id, rec[1])
		bookInstances = append(bookInstances, bi)
	}
}

func saveBookInstanceToCSV(bi *book.BookInstance) {
	file, _ := os.OpenFile(bookInstanceCSVPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer file.Close()

	w := csv.NewWriter(file)
	defer w.Flush()

	_ = w.Write([]string{strconv.Itoa(bi.GetID()), bi.ISBN})
}
