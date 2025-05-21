package repository

import (
	"encoding/csv"
	"gotus/internal/model/book"
	"os"
	"path/filepath"
	"strconv"
	"sync"
)

type BookInstanceRepository struct {
	bookInstances      []*book.BookInstance
	dataDir            string
	filename           string
	bookInstancesMutex sync.Mutex
}

func NewBookInstanceRepository(dataDir string) *BookInstanceRepository {
	return &BookInstanceRepository{
		bookInstances: []*book.BookInstance{},
		dataDir:       dataDir,
		filename:      "bookInstances.csv",
	}
}

func (r *BookInstanceRepository) StoreBookInstance(bi *book.BookInstance) {
	r.bookInstancesMutex.Lock()
	defer r.bookInstancesMutex.Unlock()

	r.bookInstances = append(r.bookInstances, bi)
	r.saveBookInstanceToCSV(bi)
}

func (r *BookInstanceRepository) GetBookInstances() ([]*book.BookInstance, int) {
	r.bookInstancesMutex.Lock()
	defer r.bookInstancesMutex.Unlock()

	return r.bookInstances, len(r.bookInstances)
}

func (r *BookInstanceRepository) LoadBookInstancesFromCSV() {
	file, err := os.Open(filepath.Join(r.dataDir, r.filename))
	if err != nil {
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, _ := reader.ReadAll()

	r.bookInstancesMutex.Lock()
	defer r.bookInstancesMutex.Unlock()

	for _, rec := range records {
		id, _ := strconv.Atoi(rec[0])
		bi := book.NewBookInstance(id, rec[1])
		r.bookInstances = append(r.bookInstances, bi)
	}
}

func (r *BookInstanceRepository) saveBookInstanceToCSV(bi *book.BookInstance) {
	file, _ := os.OpenFile(filepath.Join(r.dataDir, r.filename), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer file.Close()

	w := csv.NewWriter(file)
	defer w.Flush()

	_ = w.Write([]string{strconv.Itoa(bi.GetID()), bi.ISBN})
}
