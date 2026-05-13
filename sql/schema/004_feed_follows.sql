-- +goose Up
CREATE TABLE feed_follows (
    id BIGSERIAL PRIMARY KEY,
    feed_id BIGINT NOT NULL,
    user_id UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    FOREIGN KEY (feed_id) REFERENCES feeds(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE (feed_id, user_id)
);



-- +goose Down
DROP TABLE feed_follows;