package repository

import (
	"database/sql"
	"gotus/internal/model/reservation"
	"time"
)

type ReservationPostgresRepo struct {
	db *sql.DB
}

func NewReservationPostgresRepo(db *sql.DB) ReservationRepository {
	return &ReservationPostgresRepo{db: db}
}

func (r *ReservationPostgresRepo) StoreReservation(res *reservation.Reservation) error {
	query := `
		INSERT INTO reservations (book_instance_id, user_id, reservation_status_id, start_date, end_date)
		VALUES ($1, $2, $3, $4, $5)`
	_, err := r.db.Exec(query,
		res.BookInstanceID, res.UserID, res.ReservationStatusID,
		res.StartDate.Format("2006-01-02"),
		res.EndDate.Format("2006-01-02"),
	)
	return err
}

func (r *ReservationPostgresRepo) GetReservations() ([]*reservation.Reservation, int, error) {
	query := `SELECT id, book_instance_id, user_id, reservation_status_id, start_date, end_date FROM reservations`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	reservations := []*reservation.Reservation{}
	for rows.Next() {
		res := &reservation.Reservation{}
		var start, end string
		if err := rows.Scan(&res.Id, &res.BookInstanceID, &res.UserID, &res.ReservationStatusID, &start, &end); err != nil {
			return nil, 0, err
		}
		res.StartDate, _ = time.Parse("2006-01-02", start)
		res.EndDate, _ = time.Parse("2006-01-02", end)
		reservations = append(reservations, res)
	}
	return reservations, len(reservations), nil
}

func (r *ReservationPostgresRepo) UpdateReservationById(id int, updated *reservation.Reservation) (bool, error) {
	query := `
		UPDATE reservations
		SET book_instance_id=$1, user_id=$2, reservation_status_id=$3, start_date=$4, end_date=$5
		WHERE id=$6`
	res, err := r.db.Exec(query,
		updated.BookInstanceID, updated.UserID, updated.ReservationStatusID,
		updated.StartDate.Format("2006-01-02"), updated.EndDate.Format("2006-01-02"),
		id,
	)
	if err != nil {
		return false, err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return false, err
	}
	return affected > 0, nil
}

func (r *ReservationPostgresRepo) FindReservationById(id int) (*reservation.Reservation, bool, error) {
	query := `SELECT id, book_instance_id, user_id, reservation_status_id, start_date, end_date FROM reservations WHERE id=$1`
	row := r.db.QueryRow(query, id)

	res := &reservation.Reservation{}
	var start, end string
	err := row.Scan(&res.Id, &res.BookInstanceID, &res.UserID, &res.ReservationStatusID, &start, &end)
	if err == sql.ErrNoRows {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}

	res.StartDate, _ = time.Parse("2006-01-02", start)
	res.EndDate, _ = time.Parse("2006-01-02", end)
	return res, true, nil
}

func (r *ReservationPostgresRepo) DeleteReservationById(id int) (bool, error) {
	query := `DELETE FROM reservations WHERE id=$1`
	res, err := r.db.Exec(query, id)
	if err != nil {
		return false, err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return false, err
	}
	return affected > 0, nil
}

func (r *ReservationPostgresRepo) HasActiveReservation(bookInstanceID int) (bool, error) {
	query := `
		SELECT EXISTS(
			SELECT 1 FROM reservations 
			WHERE book_instance_id=$1 AND reservation_status_id IN (1, 2)
		)`
	row := r.db.QueryRow(query, bookInstanceID)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}
