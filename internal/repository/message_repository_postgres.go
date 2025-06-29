package repository

import (
	"database/sql"
	"gotus/internal/model/message"
)

type EmailMessagePostgresRepo struct {
	db *sql.DB
}

func NewEmailMessagePostgresRepo(db *sql.DB) EmailMessageRepository {
	return &EmailMessagePostgresRepo{db: db}
}

func (r *EmailMessagePostgresRepo) InsertMessage(m *message.Message) error {
	query := `
		INSERT INTO messages (email, send_date, email_subject, body)
		VALUES ($1, $2, $3, $4)
	`
	_, err := r.db.Exec(
		query,
		m.Email,
		m.SendDate,
		m.EmailSubject,
		m.Body,
	)
	return err
}
