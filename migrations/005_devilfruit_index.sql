-- +goose Up
CREATE INDEX devilfruits_search_idx ON devilfruits USING GIN (to_tsvector('english', name || ' ' || description));

-- +goose Down
DROP INDEX IF EXISTS  devilfruits_search_idx;