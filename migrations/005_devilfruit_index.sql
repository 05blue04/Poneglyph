-- +goose Up
CREATE INDEX devilfruits_search_idx ON devilfruits USING GIN (to_tsvector('english', name || ' ' || description));
CREATE INDEX crews_search_idx ON crews USING GIN (to_tsvector('english', name || ' ' || description));

-- +goose Down
DROP INDEX IF EXISTS  devilfruits_search_idx;
DROP INDEX IF EXISTS  crews_search_idx;