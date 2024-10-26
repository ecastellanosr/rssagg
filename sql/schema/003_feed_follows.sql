-- +goose Up
CREATE TABLE feed_follows (
    id UUID NOT NULL PRIMARY KEY,
    created_at timestamp NOT NULL,
    updated_at timestamp NOT NULL,
    feed_id UUID NOT NULL REFERENCES feeds(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE, 
    CONSTRAINT feed_user UNIQUE (feed_id,user_id)
);

-- +goose Down
DROP TABLE feed_follows;