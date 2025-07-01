package repository

import (
	"database/sql"
	"gotus/internal/model/book"
)

type BookInstancePostgresRepo struct {
	db *sql.DB
}

func NewBookInstancePostgresRepo(db *sql.DB) BookInstanceRepository {
	return &BookInstancePostgresRepo{db: db}
}

func (r *BookInstancePostgresRepo) StoreBookInstance(bi *book.BookInstance) error {
	query := `INSERT INTO book_instances (isbn) VALUES ($1)`
	_, err := r.db.Exec(query, bi.ISBN)
	return err
}

func (r *BookInstancePostgresRepo) GetBookInstances() ([]*book.BookInstance, int, error) {
	query := `SELECT id, isbn FROM book_instances`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	instances := []*book.BookInstance{}
	for rows.Next() {
		bi := &book.BookInstance{}
		if err := rows.Scan(&bi.Id, &bi.ISBN); err != nil {
			return nil, 0, err
		}
		instances = append(instances, bi)
	}

	return instances, len(instances), nil
}

func (r *BookInstancePostgresRepo) UpdateBookInstanceById(id int, updated *book.BookInstance) (bool, error) {
	query := `UPDATE book_instances SET isbn=$1 WHERE id=$2`
	res, err := r.db.Exec(query, updated.ISBN, id)
	if err != nil {
		return false, err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return false, err
	}
	return affected > 0, nil
}

func (r *BookInstancePostgresRepo) FindBookInstanceById(id int) (*book.BookInstance, bool, error) {
	query := `SELECT id, isbn FROM book_instances WHERE id=$1`
	row := r.db.QueryRow(query, id)

	bi := &book.BookInstance{}
	err := row.Scan(&bi.Id, &bi.ISBN)
	if err == sql.ErrNoRows {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}

	return bi, true, nil
}

func (r *BookInstancePostgresRepo) GetBookInstancesByISBN(isbn string) ([]*book.BookInstance, int, error) {
	query := `SELECT id, isbn FROM book_instances WHERE isbn=$1`
	rows, err := r.db.Query(query, isbn)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	instances := []*book.BookInstance{}
	for rows.Next() {
		bi := &book.BookInstance{}
		if err := rows.Scan(&bi.Id, &bi.ISBN); err != nil {
			return nil, 0, err
		}
		instances = append(instances, bi)
	}

	return instances, len(instances), nil
}

func (r *BookInstancePostgresRepo) DeleteBookInstanceById(id int) (bool, error) {
	query := `DELETE FROM book_instances WHERE id=$1`
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
