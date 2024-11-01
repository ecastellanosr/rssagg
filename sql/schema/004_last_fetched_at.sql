-- +goose Up
ALTER TABLE feeds ADD COLUMN last_fetched_at timestamp;

-- +goose Down
DROP TABLE feeds;