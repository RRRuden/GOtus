-- +goose Up
-- +goose StatementBegin
-- Добавим пользователей
INSERT INTO users (name, email) VALUES 
('Иван Иванов', 'ivan@example.com'),
('Петр Петров"', 'petrov@example.com'),
('Анна Сергеева', 'anna@example.com'),
('Мария Кузнецова', 'maria@example.com'),
('Дмитрий Орлов', 'dmitry@example.com');

-- Добавим книги
INSERT INTO books (isbn, title, author, year) VALUES
('978-3-16-148410-1', 'Преступление и наказание', 'Ф. М. Достоевский', 1866),
('978-3-16-148410-2', 'Война и мир', 'Л. Н. Толстой', 1869),
('978-3-16-148410-3', '1984', 'Джордж Оруэлл', 1949),
('978-3-16-148410-4', 'Унесённые ветром', 'Маргарет Митчелл', 1936),
('978-3-16-148410-5', 'Мастер и Маргарита', 'Михаил Булгаков', 1966);

-- Добавим экземпляры книг
INSERT INTO book_instances (isbn) VALUES
('978-3-16-148410-1'),
('978-3-16-148410-2'),
('978-3-16-148410-2'),
('978-3-16-148410-3'),
('978-3-16-148410-4'),
('978-3-16-148410-2'),
('978-3-16-148410-5'),
('978-3-16-148410-5');

-- Добавим статусы бронирования
INSERT INTO reservation_statuses (id,name) VALUES
(1, 'Забронирована'),
(2, 'Продлена'),
(3, 'Отменена'),
(4, 'Завершена');

-- Добавим бронирование
INSERT INTO reservations (book_instance_id, user_id, start_date, end_date, reservation_status_id) VALUES
(1, 1, '2025-06-01', '2025-06-15', 1),
(2, 2, '2025-06-05', '2025-06-20', 2);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- +goose StatementEnd
