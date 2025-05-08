package repository

import (
	"encoding/csv"
	"gotus/internal/model/reservation"
	"os"
	"strconv"
	"sync"
	"time"
)

const reservationCSVPath = "./data/reservations.csv"

var (
	reservations      []*reservation.Reservation
	reservationsMutex sync.Mutex
)

func StoreReservation(r *reservation.Reservation) {
	reservationsMutex.Lock()
	defer reservationsMutex.Unlock()
	reservations = append(reservations, r)
	saveReservationToCSV(r)
}

func GetReservations() ([]*reservation.Reservation, int) {
	reservationsMutex.Lock()
	defer reservationsMutex.Unlock()
	return reservations, len(reservations)
}

func LoadReservationsFromCSV() {
	file, err := os.Open(reservationCSVPath)
	if err != nil {
		return
	}
	defer file.Close()

	r := csv.NewReader(file)
	records, _ := r.ReadAll()

	reservationsMutex.Lock()
	defer reservationsMutex.Unlock()

	for _, rec := range records {
		id, _ := strconv.Atoi(rec[0])
		bid, _ := strconv.Atoi(rec[1])
		uid, _ := strconv.Atoi(rec[2])
		sid, _ := strconv.Atoi(rec[3])
		start, _ := time.Parse(time.RFC3339, rec[4])
		end, _ := time.Parse(time.RFC3339, rec[5])

		res := reservation.NewReservation(id, bid, uid, sid, start, end)
		reservations = append(reservations, res)
	}
}

func saveReservationToCSV(r *reservation.Reservation) {
	file, _ := os.OpenFile(reservationCSVPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
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
