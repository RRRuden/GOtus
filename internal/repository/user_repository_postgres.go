package repository

import (
	"database/sql"
	"gotus/internal/model/user"
)

type UserPostgresRepo struct {
	db *sql.DB
}

func NewUserPostgresRepo(db *sql.DB) UserRepository {
	return &UserPostgresRepo{db: db}
}

func (r *UserPostgresRepo) StoreUser(u *user.User) error {
	query := `INSERT INTO users (email, name) VALUES ($1, $2)`
	_, err := r.db.Exec(query, u.Email, u.Name)
	return err
}

func (r *UserPostgresRepo) GetUsers() ([]*user.User, int, error) {
	query := `SELECT id, email, name FROM users`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	users := []*user.User{}
	for rows.Next() {
		u := &user.User{}
		if err := rows.Scan(&u.Id, &u.Email, &u.Name); err != nil {
			return nil, 0, err
		}
		users = append(users, u)
	}

	return users, len(users), nil
}

func (r *UserPostgresRepo) UpdateUserById(id int, updated *user.User) (bool, error) {
	query := `UPDATE users SET email=$1, name=$2 WHERE id=$3`
	res, err := r.db.Exec(query, updated.Email, updated.Name, id)
	if err != nil {
		return false, err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return false, err
	}
	return affected > 0, nil
}

func (r *UserPostgresRepo) FindUserById(id int) (*user.User, bool, error) {
	query := `SELECT id, email, name FROM users WHERE id=$1`
	row := r.db.QueryRow(query, id)

	u := &user.User{}
	err := row.Scan(&u.Id, &u.Email, &u.Name)
	if err == sql.ErrNoRows {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}
	return u, true, nil
}

func (r *UserPostgresRepo) FindUserByEmail(email string) (*user.User, bool, error) {
	query := `SELECT id, email, name FROM users WHERE email=$1`
	row := r.db.QueryRow(query, email)

	u := &user.User{}
	err := row.Scan(&u.Id, &u.Email, &u.Name)
	if err == sql.ErrNoRows {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}
	return u, true, nil
}

func (r *UserPostgresRepo) DeleteUserById(id int) (bool, error) {
	query := `DELETE FROM users WHERE id=$1`
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
