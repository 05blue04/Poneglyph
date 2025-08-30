-- +goose Up
CREATE INDEX characters_search_idx ON characters USING GIN (to_tsvector('english', name || ' ' || description));

-- +goose Down
DROP INDEX IF EXISTS  characters_search_idx;