package repository

import (
	"database/sql"
	"gotus/internal/model/book"
)

type BookPostgresRepo struct {
	db *sql.DB
}

func NewBookPostgresRepo(db *sql.DB) BookRepository {
	return &BookPostgresRepo{db: db}
}

func (r *BookPostgresRepo) StoreBook(b *book.Book) error {
	query := `INSERT INTO books (isbn, title, author, year) VALUES ($1, $2, $3, $4)`
	_, err := r.db.Exec(query, b.GetISBN(), b.Title, b.Author, b.Year)
	return err
}

func (r *BookPostgresRepo) GetBooks() ([]*book.Book, int, error) {
	query := `SELECT isbn, title, author, year FROM books`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	books := []*book.Book{}
	for rows.Next() {
		b := &book.Book{}
		if err := rows.Scan(&b.Isbn, &b.Title, &b.Author, &b.Year); err != nil {
			return nil, 0, err
		}
		books = append(books, b)
	}

	return books, len(books), nil
}

func (r *BookPostgresRepo) UpdateBookByISBN(isbn string, updatedBook *book.Book) (bool, error) {
	query := `UPDATE books SET title=$1, author=$2, year=$3 WHERE isbn=$4`
	res, err := r.db.Exec(query, updatedBook.Title, updatedBook.Author, updatedBook.Year, isbn)
	if err != nil {
		return false, err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return false, err
	}

	return affected > 0, nil
}

func (r *BookPostgresRepo) DeleteBookByISBN(isbn string) (bool, error) {
	query := `DELETE FROM books WHERE isbn=$1`
	res, err := r.db.Exec(query, isbn)
	if err != nil {
		return false, err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return false, err
	}

	return affected > 0, nil
}

func (r *BookPostgresRepo) FindBookByISBN(isbn string) (*book.Book, bool, error) {
	query := `SELECT isbn, title, author, year FROM books WHERE isbn=$1`
	row := r.db.QueryRow(query, isbn)

	b := &book.Book{}
	err := row.Scan(&b.Isbn, &b.Title, &b.Author, &b.Year)
	if err == sql.ErrNoRows {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}

	return b, true, nil
}
