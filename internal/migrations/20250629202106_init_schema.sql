-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    email TEXT UNIQUE NOT NULL
);

CREATE TABLE books (
    isbn VARCHAR(20) PRIMARY KEY,
    title TEXT NOT NULL,
    author TEXT NOT NULL,
    year INT NOT NULL
);

CREATE TABLE book_instances (
    id SERIAL PRIMARY KEY,
    isbn VARCHAR(20) NOT NULL REFERENCES books(isbn) ON DELETE CASCADE
);

CREATE TABLE reservation_statuses (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL
);

CREATE TABLE reservations (
    id SERIAL PRIMARY KEY,
    book_instance_id INT NOT NULL REFERENCES book_instances(id) ON DELETE CASCADE,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    reservation_status_id INT NOT NULL REFERENCES reservation_statuses(id)
);

CREATE TABLE messages (
    id SERIAL PRIMARY KEY,
    email TEXT NOT NULL,
    send_date TIMESTAMP NOT NULL,
    email_subject TEXT NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS messages;
DROP TABLE IF EXISTS reservations;
DROP TABLE IF EXISTS reservation_statuses;
DROP TABLE IF EXISTS book_instances;
DROP TABLE IF EXISTS books;
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
