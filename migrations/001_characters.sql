-- +goose Up
CREATE TABLE IF NOT EXISTS characters (
    id bigserial PRIMARY KEY,
    name text NOT NULL,
    description TEXT NOT NULL,
    age integer NOT NULL,
    fruit text NOT NULL,
    bounty text NOT NULL,
    origin text NOT NULL,
    debut text NOT NULL,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW()
);

-- +goose Down
DROP TABLE characters;