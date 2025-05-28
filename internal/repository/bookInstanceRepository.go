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

func (r *BookInstanceRepository) UpdateBookInstanceById(id int, updatedBookInstance *book.BookInstance) bool {
	r.bookInstancesMutex.Lock()
	defer r.bookInstancesMutex.Unlock()

	found := false
	for i, b := range r.bookInstances {
		if b.GetID() == id {
			r.bookInstances[i] = updatedBookInstance
			found = true
			break
		}
	}

	if found {
		r.saveAllToCSV()
	}

	return found
}

func (r *BookInstanceRepository) FindBookInstanceById(id int) (*book.BookInstance, bool) {
	r.bookInstancesMutex.Lock()
	defer r.bookInstancesMutex.Unlock()
	for _, b := range r.bookInstances {
		if b.GetID() == id {
			return b, true
		}
	}
	return nil, false
}

func (r *BookInstanceRepository) DeleteBookInstanceById(id int) bool {
	r.bookInstancesMutex.Lock()
	defer r.bookInstancesMutex.Unlock()

	found := false
	var updated []*book.BookInstance
	for _, b := range r.bookInstances {
		if b.GetID() != id {
			updated = append(updated, b)
		} else {
			found = true
		}
	}

	if found {
		r.bookInstances = updated
		r.saveAllToCSV()
	}

	return found
}

func (r *BookInstanceRepository) saveBookInstanceToCSV(bi *book.BookInstance) {
	file, _ := os.OpenFile(filepath.Join(r.dataDir, r.filename), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer file.Close()

	w := csv.NewWriter(file)
	defer w.Flush()

	_ = w.Write([]string{strconv.Itoa(bi.GetID()), bi.ISBN})
}

func (r *BookInstanceRepository) saveAllToCSV() {
	file, _ := os.Create(filepath.Join(r.dataDir, r.filename))
	defer file.Close()

	w := csv.NewWriter(file)
	defer w.Flush()

	for _, bi := range r.bookInstances {
		_ = w.Write([]string{strconv.Itoa(bi.GetID()), bi.ISBN})
	}
}
