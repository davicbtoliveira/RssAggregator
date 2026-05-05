-- +goose Up
CREATE TABLE feed_follows (
  id UUID PRIMARY KEY NOT NULL,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  user_id UUID UNIQUE NOT NULL REFERENCES users(id)
                  ON DELETE CASCADE,
  feed_id UUID UNIQUE NOT NULL REFERENCES feeds(id)
                  ON DELETE CASCADE
);

-- +goose Down
DROP TABLE feed_follows;
