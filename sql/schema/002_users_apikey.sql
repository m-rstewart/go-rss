-- +goose Up
ALTER TABLE users
ADD COLUMN api_key VARCHAR(64) DEFAULT encode(sha256(random()::text::bytea), 'hex') NOT NULL;

-- +goose Down
ALTER TABLE users
DROP COLUMN api_key;