package repository

import (
	"encoding/csv"
	"gotus/internal/model/reservation"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"
)

type ReservationRepository struct {
	reservations      []*reservation.Reservation
	dataDir           string
	filename          string
	reservationsMutex sync.Mutex
}

func NewReservationRepository(dataDir string) *ReservationRepository {
	return &ReservationRepository{
		reservations: []*reservation.Reservation{},
		dataDir:      dataDir,
		filename:     "reservations.csv",
	}
}

func (repo *ReservationRepository) StoreReservation(r *reservation.Reservation) {
	repo.reservationsMutex.Lock()
	defer repo.reservationsMutex.Unlock()
	repo.reservations = append(repo.reservations, r)
	repo.saveReservationToCSV(r)
}

func (repo *ReservationRepository) GetReservations() ([]*reservation.Reservation, int) {
	repo.reservationsMutex.Lock()
	defer repo.reservationsMutex.Unlock()
	return repo.reservations, len(repo.reservations)
}

func (repo *ReservationRepository) LoadReservationsFromCSV() {
	file, err := os.Open(filepath.Join(repo.dataDir, repo.filename))
	if err != nil {
		return
	}
	defer file.Close()

	r := csv.NewReader(file)
	records, _ := r.ReadAll()

	repo.reservationsMutex.Lock()
	defer repo.reservationsMutex.Unlock()

	for _, rec := range records {
		id, _ := strconv.Atoi(rec[0])
		bid, _ := strconv.Atoi(rec[1])
		uid, _ := strconv.Atoi(rec[2])
		sid, _ := strconv.Atoi(rec[3])
		start, _ := time.Parse(time.RFC3339, rec[4])
		end, _ := time.Parse(time.RFC3339, rec[5])

		res := reservation.NewReservation(id, bid, uid, sid, start, end)
		repo.reservations = append(repo.reservations, res)
	}
}

func (r *ReservationRepository) UpdateReservationById(id int, updatedReservation *reservation.Reservation) bool {
	r.reservationsMutex.Lock()
	defer r.reservationsMutex.Unlock()

	found := false
	for i, b := range r.reservations {
		if b.GetID() == id {
			r.reservations[i] = updatedReservation
			found = true
			break
		}
	}

	if found {
		r.saveAllToCSV()
	}

	return found
}

func (repo *ReservationRepository) FindReservationById(id int) (*reservation.Reservation, bool) {
	repo.reservationsMutex.Lock()
	defer repo.reservationsMutex.Unlock()
	for _, r := range repo.reservations {
		if r.GetID() == id {
			return r, true
		}
	}
	return nil, false
}

func (repo *ReservationRepository) DeleteReservationById(id int) bool {
	repo.reservationsMutex.Lock()
	defer repo.reservationsMutex.Unlock()

	found := false
	var updated []*reservation.Reservation
	for _, r := range repo.reservations {
		if r.GetID() != id {
			updated = append(updated, r)
		} else {
			found = true
		}
	}

	if found {
		repo.reservations = updated
		repo.saveAllToCSV()
	}

	return found
}

func (repo *ReservationRepository) saveReservationToCSV(r *reservation.Reservation) {
	file, _ := os.OpenFile(filepath.Join(repo.dataDir, repo.filename), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer file.Close()

	w := csv.NewWriter(file)
	defer w.Flush()

	_ = w.Write([]string{
		strconv.Itoa(r.GetID()),
		strconv.Itoa(r.BookInstanceID),
		strconv.Itoa(r.UserID),
		strconv.Itoa(r.ReservationStatusID),
		r.StartDate.Format(time.RFC3339),
		r.EndDate.Format(time.RFC3339),
	})
}

func (repo *ReservationRepository) saveAllToCSV() {
	file, _ := os.Create(filepath.Join(repo.dataDir, repo.filename))
	defer file.Close()

	w := csv.NewWriter(file)
	defer w.Flush()

	for _, r := range repo.reservations {
		_ = w.Write([]string{
			strconv.Itoa(r.GetID()),
			strconv.Itoa(r.BookInstanceID),
			strconv.Itoa(r.UserID),
			strconv.Itoa(r.ReservationStatusID),
			r.StartDate.Format(time.RFC3339),
			r.EndDate.Format(time.RFC3339),
		})
	}
}
