-- +goose Up
CREATE INDEX IF NOT EXISTS characters_name_idx ON characters USING GIN (to_tsvector('simple', name));

-- +goose Down
DROP INDEX IF EXISTS  characters_title_idx;