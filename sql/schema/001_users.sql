-- +goose Up
CREATE TABLE users (
    id uuid NOT NULL PRIMARY KEY,
    created_at timestamp NOT NULL,
    updated_at timestamp NOT NULL,
    name text not null 
);

-- +goose Down
DROP TABLE users;