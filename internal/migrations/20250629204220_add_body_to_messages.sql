-- +goose Up
ALTER TABLE messages
ADD COLUMN body TEXT;

-- +goose Down
ALTER TABLE messages
DROP COLUMN body;
