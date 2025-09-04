-- +goose Up
CREATE TABLE IF NOT EXISTS api_keys (
    id bigserial PRIMARY KEY,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    updated_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    name text NOT NULL, 
    key_hash text UNIQUE NOT NULL, 
    is_active boolean NOT NULL DEFAULT true,
    last_used_at timestamp(0) with time zone,
    expires_at timestamp(0) with time zone 
);

CREATE INDEX idx_api_keys_hash ON api_keys(key_hash);

-- +goose Down
DROP TABLE IF EXISTS api_keys;
DROP INDEX IF EXISTS idx_api_keys_hash;
